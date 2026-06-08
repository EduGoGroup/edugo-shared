package logger

import (
	"context"
	"log/slog"
)

// gcpSeverityHandler es un slog.Handler envoltorio (wrapper) que AGREGA un
// atributo top-level `severity` a cada entrada, derivado de record.Level y
// mapeado al vocabulario de Google Cloud Logging (LogSeverity).
//
// Es puramente aditivo: el campo `level` original de slog se conserva intacto.
// GCP solo reconoce `severity` en el top-level del JSON; por eso el wrapper
// mantiene el handler interno SIN grupos y reconstruye los grupos/atributos del
// call-site como atributos anidados del record en Handle. Así `severity` siempre
// aterriza en la raíz, aunque el call-site haya abierto grupos vía WithGroup.
//
// Decisión de diseño: si delegáramos WithGroup directamente al handler interno,
// el atributo `severity` que agregamos en Handle quedaría anidado dentro del
// grupo abierto (GCP no lo vería). Para evitarlo, el wrapper se encarga él mismo
// del prefijado de grupos y los atributos contextuales.
type gcpSeverityHandler struct {
	// inner es el handler base SIN grupos abiertos: emite siempre en la raíz.
	inner slog.Handler
	// stages modela los WithGroup/WithAttrs acumulados del call-site, en orden.
	// El nivel 0 (group == "") agrupa los atributos añadidos antes de abrir el
	// primer grupo; cada stage posterior abre un grupo con sus propios atributos.
	stages []severityStage
}

// severityStage representa un nivel de anidamiento: un grupo (vacío para la raíz)
// con los atributos contextuales agregados mientras ese grupo estaba activo.
type severityStage struct {
	group string
	attrs []slog.Attr
}

// newGCPSeverityHandler envuelve un slog.Handler para inyectar el campo
// `severity` compatible con Google Cloud Logging en cada entrada.
func newGCPSeverityHandler(inner slog.Handler) slog.Handler {
	return &gcpSeverityHandler{
		inner:  inner,
		stages: []severityStage{{}}, // stage raíz, sin grupo
	}
}

// Enabled delega al handler interno.
func (h *gcpSeverityHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

// Handle agrega el atributo `severity` (top-level) y delega al handler interno.
//
// Reconstruye los atributos del call-site (acumulados vía WithAttrs/WithGroup)
// como un árbol anidado de slog.Group y los añade junto a los atributos propios
// del record. `severity` se añade en el stage raíz, garantizando que quede en la
// raíz del JSON. record es una copia (slog pasa Record por valor), por lo que
// mutarlo aquí es seguro.
func (h *gcpSeverityHandler) Handle(ctx context.Context, record slog.Record) error {
	// Recolectar los atributos propios del record (los del call-site del log).
	recordAttrs := make([]slog.Attr, 0, record.NumAttrs())
	record.Attrs(func(a slog.Attr) bool {
		recordAttrs = append(recordAttrs, a)
		return true
	})

	// Construir el árbol de stages de adentro hacia afuera: el stage más interno
	// contiene sus atributos contextuales + los atributos propios del record.
	rebuilt := buildStageAttrs(h.stages, recordAttrs)

	// Record limpio: mismos metadatos, sin los atributos originales (ya están en
	// rebuilt), más severity en la raíz.
	out := slog.NewRecord(record.Time, record.Level, record.Message, record.PC)
	out.AddAttrs(slog.String(FieldSeverity, gcpSeverityFromLevel(record.Level)))
	out.AddAttrs(rebuilt...)

	return h.inner.Handle(ctx, out)
}

// buildStageAttrs colapsa la lista de stages (más los atributos propios del
// record) en una lista de atributos de la raíz, anidando cada grupo con
// slog.GroupValue. recordAttrs pertenece al stage más profundo (el grupo activo
// al momento del log).
//
// Se procesa de adentro hacia afuera: `acc` arranca con los atributos del stage
// más profundo + los del record; en cada paso hacia un stage padre, `acc` se
// envuelve como un grupo y se antepone con los atributos contextuales del padre.
func buildStageAttrs(stages []severityStage, recordAttrs []slog.Attr) []slog.Attr {
	if len(stages) == 0 {
		return recordAttrs
	}

	// Stage más profundo: sus atributos + los del record, sin envolver todavía.
	deepest := stages[len(stages)-1]
	acc := append(append([]slog.Attr{}, deepest.attrs...), recordAttrs...)

	// Plegar hacia la raíz. Para i >= 1, el stage tiene grupo: envolvemos acc en
	// ese grupo y le anteponemos los atributos del stage padre (i-1).
	for i := len(stages) - 1; i >= 1; i-- {
		grouped := slog.Attr{Key: stages[i].group, Value: slog.GroupValue(acc...)}
		parent := stages[i-1]
		acc = append(append([]slog.Attr{}, parent.attrs...), grouped)
	}

	return acc
}

// WithAttrs agrega atributos contextuales al stage actual (el grupo más reciente)
// y devuelve un nuevo wrapper, preservando la inyección de `severity`.
func (h *gcpSeverityHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	stages := cloneStages(h.stages)
	last := &stages[len(stages)-1]
	last.attrs = append(append([]slog.Attr{}, last.attrs...), attrs...)
	return &gcpSeverityHandler{inner: h.inner, stages: stages}
}

// WithGroup abre un nuevo stage de grupo y devuelve un nuevo wrapper.
//
// No delega WithGroup al handler interno: el wrapper gestiona el anidamiento por
// sí mismo (ver Handle) para que `severity` permanezca en la raíz.
func (h *gcpSeverityHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	stages := cloneStages(h.stages)
	stages = append(stages, severityStage{group: name})
	return &gcpSeverityHandler{inner: h.inner, stages: stages}
}

func cloneStages(in []severityStage) []severityStage {
	out := make([]severityStage, len(in))
	for i, s := range in {
		out[i] = severityStage{
			group: s.group,
			attrs: append([]slog.Attr{}, s.attrs...),
		}
	}
	return out
}

// gcpSeverityFromLevel mapea un slog.Level al vocabulario LogSeverity de Google
// Cloud Logging. El mapeo es por umbrales para ser robusto ante niveles custom:
//
//	level <  Info    → "DEBUG"
//	level <  Warn    → "INFO"
//	level <  Error   → "WARNING"   (GCP usa WARNING, no WARN)
//	level <  Error+4 → "ERROR"
//	level >= Error+4 → "CRITICAL"
//
// Referencia: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogSeverity
func gcpSeverityFromLevel(level slog.Level) string {
	switch {
	case level < slog.LevelInfo:
		return "DEBUG"
	case level < slog.LevelWarn:
		return "INFO"
	case level < slog.LevelError:
		return "WARNING"
	case level < slog.LevelError+4:
		return "ERROR"
	default:
		return "CRITICAL"
	}
}
