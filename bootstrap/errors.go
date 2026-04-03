package bootstrap

import "fmt"

// =============================================================================
// ERROR TYPES
// =============================================================================

// ErrMissingFactory se retorna cuando falta una factory requerida.
//
//nolint:errname // Mantener nombre Err* para compatibilidad con API existente
type ErrMissingFactory struct {
	Resource string
}

func (e ErrMissingFactory) Error() string {
	return "missing required factory: " + e.Resource
}

// ErrConnectionFailed se retorna cuando una conexion a un recurso falla.
//
//nolint:errname // Mantener nombre Err* para consistencia con ErrMissingFactory
type ErrConnectionFailed struct {
	Resource string
	Err      error
}

func (e ErrConnectionFailed) Error() string {
	return fmt.Sprintf("bootstrap/%s: connection failed: %v", e.Resource, e.Err)
}

func (e ErrConnectionFailed) Unwrap() error {
	return e.Err
}
