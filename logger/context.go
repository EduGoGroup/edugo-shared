package logger

import (
	"context"
	"log/slog"
)

type contextKey struct{}

// NewContext retorna un nuevo contexto con el slog.Logger proporcionado embebido.
// Usar en middleware para propagar un logger enriquecido a lo largo de la petición.
func NewContext(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, l)
}

// FromContext extrae el *slog.Logger del contexto.
// Retorna slog.Default() si no se almacenó ningún logger en el contexto.
func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(contextKey{}).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}

// L es un alias abreviado de FromContext.
func L(ctx context.Context) *slog.Logger {
	return FromContext(ctx)
}
