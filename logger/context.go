package logger

import (
	"context"
	"log/slog"
)

type contextKey struct{}

// NewContext retorna un nuevo contexto con el Logger proporcionado embebido.
// Usar en middleware para propagar un logger enriquecido a lo largo de la petición.
func NewContext(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, l)
}

// FromContext extrae el Logger del contexto.
// Retorna un adapter de slog.Default() si no se almacenó ningún logger en el contexto.
func FromContext(ctx context.Context) Logger {
	if l, ok := ctx.Value(contextKey{}).(Logger); ok {
		return l
	}
	return NewSlogAdapter(slog.Default())
}

// L es un alias abreviado de FromContext.
func L(ctx context.Context) Logger {
	return FromContext(ctx)
}
