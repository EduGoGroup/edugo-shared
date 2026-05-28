package bootstrap

// =============================================================================
// GORM FUNCTIONAL OPTIONS
// =============================================================================
//
// Opciones funcionales para configurar conexiones GORM en bootstrap/postgres.
// Definidas en el modulo raiz para que los consumidores puedan configurar
// opciones sin importar el sub-modulo directamente.
// =============================================================================

// GORMOption configura la conexion GORM.
type GORMOption func(*GORMOptions)

// GORMOptions almacena todas las opciones configurables para GORM.
type GORMOptions struct {
	// Logger acepta gorm.logger.Interface sin importar gorm en este modulo.
	// El sub-modulo postgres hace type assertion al tipo concreto.
	Logger any

	// SimpleProtocol habilita pgx.QueryExecModeSimpleProtocol.
	// Requerido para PgBouncer (transaction mode) y Neon.
	// Default: true.
	SimpleProtocol bool

	// PrepareStmt controla gorm.Config.PrepareStmt.
	// Debe ser false cuando SimpleProtocol es true.
	// Default: false.
	PrepareStmt bool
}

// DefaultGORMOptions retorna las opciones por defecto para GORM.
// SimpleProtocol habilitado por defecto (compatible con PgBouncer/Neon).
func DefaultGORMOptions() GORMOptions {
	return GORMOptions{
		SimpleProtocol: true,
		PrepareStmt:    false,
	}
}

// ApplyGORMOptions aplica una lista de opciones sobre los defaults.
func ApplyGORMOptions(opts ...GORMOption) GORMOptions {
	o := DefaultGORMOptions()
	for _, fn := range opts {
		fn(&o)
	}
	return o
}

// WithGORMLogger configura el logger de GORM.
// Acepta gorm.logger.Interface — el sub-modulo postgres hace type assertion.
func WithGORMLogger(logger any) GORMOption {
	return func(o *GORMOptions) {
		o.Logger = logger
	}
}

// WithSimpleProtocol controla pgx SimpleProtocol.
// Usar true para PgBouncer/Neon (default), false para conexiones directas.
func WithSimpleProtocol(enabled bool) GORMOption {
	return func(o *GORMOptions) {
		o.SimpleProtocol = enabled
	}
}

// WithPrepareStmt controla GORM PrepareStmt.
// Debe ser false cuando SimpleProtocol es true.
func WithPrepareStmt(enabled bool) GORMOption {
	return func(o *GORMOptions) {
		o.PrepareStmt = enabled
	}
}
