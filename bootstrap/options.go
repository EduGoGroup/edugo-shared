package bootstrap

// =============================================================================
// BOOTSTRAP OPTIONS
// =============================================================================

// BootstrapOptions define las opciones de configuración para el bootstrap
type BootstrapOptions struct {
	// RequiredResources lista los recursos obligatorios que deben inicializarse
	RequiredResources []string

	// OptionalResources lista los recursos opcionales que se pueden omitir sin error
	OptionalResources []string

	// SkipHealthCheck indica si se debe omitir el health check inicial
	SkipHealthCheck bool

	// MockFactories permite inyectar factories simuladas para testing
	MockFactories *MockFactories

	// StopOnFirstError indica si se debe detener al primer error o continuar
	StopOnFirstError bool
}

// MockFactories contiene factories simuladas para testing
type MockFactories struct {
	Logger     LoggerFactory
	PostgreSQL PostgreSQLFactory
	MongoDB    MongoDBFactory
	RabbitMQ   RabbitMQFactory
	S3         S3Factory
}

// =============================================================================
// OPTION FUNCTIONS
// =============================================================================

// BootstrapOption es una función que modifica BootstrapOptions
type BootstrapOption func(*BootstrapOptions)

// WithRequiredResources especifica los recursos obligatorios
func WithRequiredResources(resources ...string) BootstrapOption {
	return func(opts *BootstrapOptions) {
		opts.RequiredResources = resources
	}
}

// WithOptionalResources especifica los recursos opcionales
func WithOptionalResources(resources ...string) BootstrapOption {
	return func(opts *BootstrapOptions) {
		opts.OptionalResources = resources
	}
}

// WithSkipHealthCheck omite el health check inicial
func WithSkipHealthCheck() BootstrapOption {
	return func(opts *BootstrapOptions) {
		opts.SkipHealthCheck = true
	}
}

// WithMockFactories inyecta factories simuladas para testing
func WithMockFactories(mocks *MockFactories) BootstrapOption {
	return func(opts *BootstrapOptions) {
		opts.MockFactories = mocks
	}
}

// WithStopOnFirstError configura si detener al primer error
func WithStopOnFirstError(stop bool) BootstrapOption {
	return func(opts *BootstrapOptions) {
		opts.StopOnFirstError = stop
	}
}

// =============================================================================
// DEFAULT OPTIONS
// =============================================================================

// DefaultBootstrapOptions retorna las opciones por defecto
func DefaultBootstrapOptions() *BootstrapOptions {
	return &BootstrapOptions{
		RequiredResources: []string{"logger"},
		OptionalResources: []string{},
		SkipHealthCheck:   false,
		MockFactories:     nil,
		StopOnFirstError:  true,
	}
}

// ApplyOptions aplica una lista de opciones a BootstrapOptions
func ApplyOptions(opts *BootstrapOptions, options ...BootstrapOption) {
	for _, opt := range options {
		opt(opts)
	}
}
