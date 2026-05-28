package tracer

import (
	"context"
	"errors"
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
)

// NewSlogHandler returns a slog.Handler that forwards records to the global
// OpenTelemetry LoggerProvider (set up by Init when cfg.Enabled=true and
// cfg.Endpoint!=""). The handler is a no-op when no LoggerProvider is registered.
//
// instrumentationName identifies the emitter inside OTel — pass the service name
// (e.g. "edugo-api-identity") so logs from different services can be filtered
// by scope in the backend.
//
// minLevel implements DA-MPH-5: the OTel/Loki exporter has its own level
// independent of the root slog level (stdout). This prevents flooding the
// observability backend when a developer raises LOGGING_LEVEL=debug locally.
// Use logger.ResolveOtelLevel(os.Getenv("OTEL_LOG_LEVEL"), appEnv) to compute
// minLevel at the call site.
func NewSlogHandler(instrumentationName string, minLevel slog.Level) slog.Handler {
	return &levelHandler{
		inner: otelslog.NewHandler(instrumentationName),
		min:   minLevel,
	}
}

// levelHandler decorates an inner slog.Handler with a minimum-level filter.
// Used to enforce DA-MPH-5 on the OTel branch of the multi-handler tree
// without touching the stdout branch.
type levelHandler struct {
	inner slog.Handler
	min   slog.Level
}

func (h *levelHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if level < h.min {
		return false
	}
	return h.inner.Enabled(ctx, level)
}

func (h *levelHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.Level < h.min {
		return nil
	}
	return h.inner.Handle(ctx, r)
}

func (h *levelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &levelHandler{inner: h.inner.WithAttrs(attrs), min: h.min}
}

func (h *levelHandler) WithGroup(name string) slog.Handler {
	return &levelHandler{inner: h.inner.WithGroup(name), min: h.min}
}

// NewMultiHandler returns a slog.Handler that fans out each record to every
// underlying handler. Used to tee logs to both stdout (existing JSON handler)
// and the OTel exporter without changing call sites.
//
// Enabled returns true if any sub-handler is enabled; Handle dispatches only
// to the sub-handlers that report Enabled for the record's level. WithAttrs and
// WithGroup propagate to all sub-handlers.
func NewMultiHandler(handlers ...slog.Handler) slog.Handler {
	return &multiHandler{handlers: handlers}
}

type multiHandler struct {
	handlers []slog.Handler
}

func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	var errs []error
	for _, h := range m.handlers {
		if !h.Enabled(ctx, r.Level) {
			continue
		}
		if err := h.Handle(ctx, r.Clone()); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	wrapped := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		wrapped[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: wrapped}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	wrapped := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		wrapped[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: wrapped}
}
