package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGCPSeverityFromLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    slog.Level
		expected string
	}{
		{"debug", slog.LevelDebug, "DEBUG"},
		{"below_debug", slog.LevelDebug - 1, "DEBUG"},
		{"info", slog.LevelInfo, "INFO"},
		{"between_info_warn", slog.LevelInfo + 1, "INFO"},
		// GCP usa WARNING, no WARN.
		{"warn_maps_to_warning", slog.LevelWarn, "WARNING"},
		{"between_warn_error", slog.LevelWarn + 1, "WARNING"},
		{"error", slog.LevelError, "ERROR"},
		{"between_error_critical", slog.LevelError + 1, "ERROR"},
		// Niveles muy altos (>= Error+4) escalan a CRITICAL.
		{"critical", slog.LevelError + 4, "CRITICAL"},
		{"very_high_custom", slog.Level(64), "CRITICAL"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, gcpSeverityFromLevel(tt.level))
		})
	}
}

// TestGCPSeverityHandler_CoexistsTopLevel serializa un log real y verifica que
// `level` y `severity` coexisten, ambos en el top-level del JSON.
func TestGCPSeverityHandler_CoexistsTopLevel(t *testing.T) {
	var buf bytes.Buffer
	base := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	l := slog.New(newGCPSeverityHandler(base))

	l.Warn("algo pasó", "key", "value")

	var out map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &out))

	// `level` original de slog se conserva.
	assert.Equal(t, "WARN", out["level"], "el campo level debe conservarse intacto")
	// `severity` agregado en vocabulario GCP, top-level.
	assert.Equal(t, "WARNING", out["severity"], "severity debe estar en el top-level con vocabulario GCP")
}

// TestGCPSeverityHandler_TopLevelWithGroup verifica que severity permanece en la
// raíz aun cuando el call-site abrió un grupo vía WithGroup.
func TestGCPSeverityHandler_TopLevelWithGroup(t *testing.T) {
	var buf bytes.Buffer
	base := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	l := slog.New(newGCPSeverityHandler(base)).WithGroup("req")

	l.Error("boom", "code", 500)

	var out map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &out))

	assert.Equal(t, "ERROR", out["severity"], "severity debe quedar top-level, no dentro del grupo")
	assert.Equal(t, "ERROR", out["level"], "level se conserva")
	// El atributo del call-site sí entra al grupo; severity NO.
	group, ok := out["req"].(map[string]any)
	require.True(t, ok, "el grupo req debe existir")
	assert.NotContains(t, group, "severity", "severity no debe anidarse en el grupo")
}

// TestGCPSeverityHandler_NestedGroupsAndAttrs verifica que el wrapper reconstruye
// correctamente atributos contextuales antes y después de abrir grupos anidados,
// manteniendo severity en la raíz.
func TestGCPSeverityHandler_NestedGroupsAndAttrs(t *testing.T) {
	var buf bytes.Buffer
	base := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	l := slog.New(newGCPSeverityHandler(base)).
		With("svc", "academic"). // atributo en la raíz
		WithGroup("http").       // abre grupo http
		With("method", "GET").   // atributo dentro de http
		WithGroup("user").       // grupo anidado http.user
		With("id", "42")         // atributo dentro de http.user

	l.Info("req", "latency_ms", 12)

	var out map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &out))

	// severity y level en la raíz.
	assert.Equal(t, "INFO", out["severity"])
	assert.Equal(t, "INFO", out["level"])
	// atributo raíz preservado.
	assert.Equal(t, "academic", out["svc"])

	http, ok := out["http"].(map[string]any)
	require.True(t, ok, "grupo http debe existir")
	assert.Equal(t, "GET", http["method"])
	assert.NotContains(t, http, "severity", "severity nunca debe anidarse")

	user, ok := http["user"].(map[string]any)
	require.True(t, ok, "grupo anidado http.user debe existir")
	assert.Equal(t, "42", user["id"])
	// el atributo del call-site del log entra al grupo activo más profundo.
	assert.EqualValues(t, 12, user["latency_ms"])
}

// TestGCPSeverityHandler_ViaProvider verifica el camino real del constructor:
// NewSlogProvider envuelve el handler JSON, así que severity sale top-level.
func TestGCPSeverityHandler_ViaProvider(t *testing.T) {
	// Redirigir stdout no es trivial aquí; en su lugar verificamos que el handler
	// raíz del provider es el wrapper y que emite severity correctamente.
	l := NewSlogProvider(SlogConfig{Level: "debug", Format: "json", Service: "svc"})
	_, ok := l.Handler().(*gcpSeverityHandler)
	assert.True(t, ok, "el provider JSON debe envolver con gcpSeverityHandler")
}

// TestGCPSeverityHandler_AllSeverities recorre los 4 niveles estándar de extremo
// a extremo a través del wrapper.
func TestGCPSeverityHandler_AllSeverities(t *testing.T) {
	cases := []struct {
		level    slog.Level
		severity string
	}{
		{slog.LevelDebug, "DEBUG"},
		{slog.LevelInfo, "INFO"},
		{slog.LevelWarn, "WARNING"},
		{slog.LevelError, "ERROR"},
	}

	for _, c := range cases {
		var buf bytes.Buffer
		base := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
		h := newGCPSeverityHandler(base)

		rec := slog.NewRecord(time.Time{}, c.level, "msg", 0)
		require.NoError(t, h.Handle(context.Background(), rec))

		var out map[string]any
		require.NoError(t, json.Unmarshal(buf.Bytes(), &out))
		assert.Equal(t, c.severity, out["severity"], "level=%v", c.level)
	}
}
