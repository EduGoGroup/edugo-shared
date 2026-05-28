package tracer

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
	"time"
)

// recordingHandler captura cada Record que recibe para asserts en tests.
// Sirve como "inner" para verificar que levelHandler filtra correctamente.
// No es thread-safe: los tests de este archivo lo usan secuencialmente.
type recordingHandler struct {
	levels  []slog.Level // niveles vistos por Handle (append-only, sin lock)
	enabled slog.Level   // mínimo que el handler reporta como Enabled
	attrs   []slog.Attr
	groups  []string
}

func (r *recordingHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= r.enabled
}

func (r *recordingHandler) Handle(_ context.Context, rec slog.Record) error {
	r.levels = append(r.levels, rec.Level)
	return nil
}

func (r *recordingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	cp := *r
	cp.attrs = append(append([]slog.Attr{}, r.attrs...), attrs...)
	return &cp
}

func (r *recordingHandler) WithGroup(name string) slog.Handler {
	cp := *r
	cp.groups = append(append([]string{}, r.groups...), name)
	return &cp
}

func newLevelHandlerWithInner(min slog.Level, inner *recordingHandler) *levelHandler {
	return &levelHandler{inner: inner, min: min}
}

func TestLevelHandler_EnabledRespectsMin(t *testing.T) {
	inner := &recordingHandler{enabled: slog.LevelDebug}
	h := newLevelHandlerWithInner(slog.LevelInfo, inner)

	if h.Enabled(context.Background(), slog.LevelDebug) {
		t.Fatalf("Enabled debería ser false para DEBUG cuando min=INFO")
	}
	if !h.Enabled(context.Background(), slog.LevelInfo) {
		t.Fatalf("Enabled debería ser true para INFO cuando min=INFO")
	}
	if !h.Enabled(context.Background(), slog.LevelWarn) {
		t.Fatalf("Enabled debería ser true para WARN cuando min=INFO")
	}
}

func TestLevelHandler_EnabledRespectsInner(t *testing.T) {
	// Si el inner reporta no-enabled (e.g. no hay LoggerProvider global),
	// el wrapper debe reflejarlo aunque el level pase el filtro.
	inner := &recordingHandler{enabled: slog.LevelError}
	h := newLevelHandlerWithInner(slog.LevelInfo, inner)

	if h.Enabled(context.Background(), slog.LevelInfo) {
		t.Fatalf("Enabled debería ser false si el inner reporta false")
	}
	if !h.Enabled(context.Background(), slog.LevelError) {
		t.Fatalf("Enabled debería ser true para ERROR cuando inner.enabled=ERROR")
	}
}

func TestLevelHandler_HandleDropsBelowMin(t *testing.T) {
	inner := &recordingHandler{enabled: slog.LevelDebug}
	h := newLevelHandlerWithInner(slog.LevelInfo, inner)
	ctx := context.Background()

	debugRec := slog.NewRecord(time.Now(), slog.LevelDebug, "debug-msg", 0)
	infoRec := slog.NewRecord(time.Now(), slog.LevelInfo, "info-msg", 0)
	warnRec := slog.NewRecord(time.Now(), slog.LevelWarn, "warn-msg", 0)

	if err := h.Handle(ctx, debugRec); err != nil {
		t.Fatalf("Handle debug devolvió error: %v", err)
	}
	if err := h.Handle(ctx, infoRec); err != nil {
		t.Fatalf("Handle info devolvió error: %v", err)
	}
	if err := h.Handle(ctx, warnRec); err != nil {
		t.Fatalf("Handle warn devolvió error: %v", err)
	}

	if len(inner.levels) != 2 {
		t.Fatalf("inner debería haber recibido 2 records (info, warn), recibió %d: %v",
			len(inner.levels), inner.levels)
	}
	if inner.levels[0] != slog.LevelInfo || inner.levels[1] != slog.LevelWarn {
		t.Fatalf("orden esperado [info, warn], obtuvo %v", inner.levels)
	}
}

func TestLevelHandler_WithAttrsPreservesMin(t *testing.T) {
	inner := &recordingHandler{enabled: slog.LevelDebug}
	h := newLevelHandlerWithInner(slog.LevelWarn, inner)

	wrapped := h.WithAttrs([]slog.Attr{slog.String("k", "v")})
	lh, ok := wrapped.(*levelHandler)
	if !ok {
		t.Fatalf("WithAttrs debería retornar *levelHandler, obtuvo %T", wrapped)
	}
	if lh.min != slog.LevelWarn {
		t.Fatalf("min debería preservarse en WithAttrs, esperado WARN obtuvo %v", lh.min)
	}
	if !lh.Enabled(context.Background(), slog.LevelWarn) {
		t.Fatalf("WithAttrs no propagó al inner (Enabled false para WARN)")
	}
	if lh.Enabled(context.Background(), slog.LevelInfo) {
		t.Fatalf("filtro min se rompió tras WithAttrs (INFO pasó cuando min=WARN)")
	}
}

func TestLevelHandler_WithGroupPreservesMin(t *testing.T) {
	inner := &recordingHandler{enabled: slog.LevelDebug}
	h := newLevelHandlerWithInner(slog.LevelError, inner)

	wrapped := h.WithGroup("scope")
	lh, ok := wrapped.(*levelHandler)
	if !ok {
		t.Fatalf("WithGroup debería retornar *levelHandler, obtuvo %T", wrapped)
	}
	if lh.min != slog.LevelError {
		t.Fatalf("min debería preservarse en WithGroup, esperado ERROR obtuvo %v", lh.min)
	}
	if !lh.Enabled(context.Background(), slog.LevelError) {
		t.Fatalf("WithGroup no propagó al inner (Enabled false para ERROR)")
	}
}

func TestLevelHandler_WithAttrsChainedPreservesMin(t *testing.T) {
	// Encadenar WithAttrs/WithGroup no debe colapsar el filtro min.
	inner := &recordingHandler{enabled: slog.LevelDebug}
	h := newLevelHandlerWithInner(slog.LevelWarn, inner)

	chained := h.
		WithAttrs([]slog.Attr{slog.String("k1", "v1")}).
		WithGroup("g1").
		WithAttrs([]slog.Attr{slog.String("k2", "v2")})

	lh, ok := chained.(*levelHandler)
	if !ok {
		t.Fatalf("chained debería retornar *levelHandler, obtuvo %T", chained)
	}
	if lh.min != slog.LevelWarn {
		t.Fatalf("min debería preservarse tras encadenar, esperado WARN obtuvo %v", lh.min)
	}
	if lh.Enabled(context.Background(), slog.LevelInfo) {
		t.Fatalf("filtro min se rompió tras encadenar (INFO pasó cuando min=WARN)")
	}
	if !lh.Enabled(context.Background(), slog.LevelWarn) {
		t.Fatalf("WARN debería pasar tras encadenar")
	}
}

func TestNewSlogHandler_FactoryReturnsLevelHandler(t *testing.T) {
	h := NewSlogHandler("test-service", slog.LevelInfo)
	if h == nil {
		t.Fatalf("NewSlogHandler retornó nil")
	}
	lh, ok := h.(*levelHandler)
	if !ok {
		t.Fatalf("NewSlogHandler debería retornar *levelHandler, obtuvo %T", h)
	}
	if lh.min != slog.LevelInfo {
		t.Fatalf("min esperado INFO, obtuvo %v", lh.min)
	}
	if lh.inner == nil {
		t.Fatalf("inner no debería ser nil")
	}
}

// TestMultiHandler_OtelBranchFilteredIndependently es el test integración
// que valida la política DA-MPH-5: el handler stdout recibe DEBUG, el
// handler OTel (con min=INFO) lo descarta — todo via el mismo MultiHandler.
func TestMultiHandler_OtelBranchFilteredIndependently(t *testing.T) {
	// stdout branch: JSON handler con LevelDebug → ve todo.
	var stdoutBuf bytes.Buffer
	stdoutHandler := slog.NewJSONHandler(&stdoutBuf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	// otel branch: levelHandler con min=INFO envolviendo un recordingHandler.
	otelInner := &recordingHandler{enabled: slog.LevelDebug}
	otelHandler := &levelHandler{inner: otelInner, min: slog.LevelInfo}

	logger := slog.New(NewMultiHandler(stdoutHandler, otelHandler))

	logger.Debug("debug-only-stdout")
	logger.Info("info-both-branches")
	logger.Warn("warn-both-branches")

	// stdout vio los 3.
	stdoutLines := strings.Split(strings.TrimRight(stdoutBuf.String(), "\n"), "\n")
	if len(stdoutLines) != 3 {
		t.Fatalf("stdout debería tener 3 líneas, tiene %d: %q", len(stdoutLines), stdoutBuf.String())
	}

	// Cada línea debe ser JSON parseable y contener msg.
	for i, line := range stdoutLines {
		var parsed map[string]any
		if err := json.Unmarshal([]byte(line), &parsed); err != nil {
			t.Fatalf("línea %d no es JSON válido: %v (line=%q)", i, err, line)
		}
		if _, ok := parsed["msg"]; !ok {
			t.Fatalf("línea %d no tiene campo msg: %v", i, parsed)
		}
	}

	// otel inner sólo vio info y warn (debug fue filtrado por levelHandler).
	if len(otelInner.levels) != 2 {
		t.Fatalf("otel inner debería haber recibido 2 records, recibió %d: %v",
			len(otelInner.levels), otelInner.levels)
	}
	if otelInner.levels[0] != slog.LevelInfo || otelInner.levels[1] != slog.LevelWarn {
		t.Fatalf("otel inner orden esperado [info, warn], obtuvo %v", otelInner.levels)
	}
}
