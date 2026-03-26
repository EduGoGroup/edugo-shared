package logger

import (
	"context"
	"log/slog"
)

type contextKey struct{}

// NewContext returns a new context with the given slog.Logger embedded.
// Use this in middleware to propagate an enriched logger through the request.
func NewContext(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, l)
}

// FromContext extracts the *slog.Logger from the context.
// Returns slog.Default() if no logger was stored in the context.
func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(contextKey{}).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}

// L is a shorthand alias for FromContext.
func L(ctx context.Context) *slog.Logger {
	return FromContext(ctx)
}
