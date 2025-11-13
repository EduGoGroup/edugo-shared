package containers

// Config contiene la configuración para el Manager de containers
type Config struct {
	// Flags para habilitar containers
	UsePostgreSQL bool
	UseMongoDB    bool
	UseRabbitMQ   bool

	// Configuraciones específicas de cada container
	PostgresConfig *PostgresConfig
	MongoConfig    *MongoConfig
	RabbitConfig   *RabbitConfig
}

// PostgresConfig configura el container de PostgreSQL
type PostgresConfig struct {
	Image       string   // Imagen Docker (default: "postgres:15-alpine")
	Database    string   // Nombre de la base de datos (default: "edugo_test")
	Username    string   // Usuario (default: "edugo_user")
	Password    string   // Contraseña (default: "edugo_pass")
	Port        string   // Puerto (default: "5432")
	InitScripts []string // Scripts SQL para ejecutar al iniciar
}

// MongoConfig configura el container de MongoDB
type MongoConfig struct {
	Image    string // Imagen Docker (default: "mongo:7.0")
	Database string // Nombre de la base de datos (default: "edugo_test")
	Username string // Usuario (default: "")
	Password string // Contraseña (default: "")
}

// RabbitConfig configura el container de RabbitMQ
type RabbitConfig struct {
	Image    string // Imagen Docker (default: "rabbitmq:3.12-alpine")
	Username string // Usuario (default: "edugo_user")
	Password string // Contraseña (default: "edugo_pass")
}

// ConfigBuilder permite construir una Config de forma fluida
type ConfigBuilder struct {
	config *Config
}

// NewConfig crea un nuevo ConfigBuilder
func NewConfig() *ConfigBuilder {
	return &ConfigBuilder{
		config: &Config{},
	}
}

// WithPostgreSQL habilita PostgreSQL con la configuración proporcionada.
// Si cfg es nil, se usan valores por defecto.
func (b *ConfigBuilder) WithPostgreSQL(cfg *PostgresConfig) *ConfigBuilder {
	b.config.UsePostgreSQL = true
	if cfg == nil {
		cfg = &PostgresConfig{}
	}
	// Aplicar defaults
	if cfg.Image == "" {
		cfg.Image = "postgres:15-alpine"
	}
	if cfg.Database == "" {
		cfg.Database = "edugo_test"
	}
	if cfg.Username == "" {
		cfg.Username = "edugo_user"
	}
	if cfg.Password == "" {
		cfg.Password = "edugo_pass"
	}
	if cfg.Port == "" {
		cfg.Port = "5432"
	}
	b.config.PostgresConfig = cfg
	return b
}

// WithMongoDB habilita MongoDB con la configuración proporcionada.
// Si cfg es nil, se usan valores por defecto.
func (b *ConfigBuilder) WithMongoDB(cfg *MongoConfig) *ConfigBuilder {
	b.config.UseMongoDB = true
	if cfg == nil {
		cfg = &MongoConfig{}
	}
	// Aplicar defaults
	if cfg.Image == "" {
		cfg.Image = "mongo:7.0"
	}
	if cfg.Database == "" {
		cfg.Database = "edugo_test"
	}
	// MongoDB sin autenticación por defecto en tests
	b.config.MongoConfig = cfg
	return b
}

// WithRabbitMQ habilita RabbitMQ con la configuración proporcionada.
// Si cfg es nil, se usan valores por defecto.
func (b *ConfigBuilder) WithRabbitMQ(cfg *RabbitConfig) *ConfigBuilder {
	b.config.UseRabbitMQ = true
	if cfg == nil {
		cfg = &RabbitConfig{}
	}
	// Aplicar defaults
	if cfg.Image == "" {
		cfg.Image = "rabbitmq:3.12-alpine"
	}
	if cfg.Username == "" {
		cfg.Username = "edugo_user"
	}
	if cfg.Password == "" {
		cfg.Password = "edugo_pass"
	}
	b.config.RabbitConfig = cfg
	return b
}

// Build construye y retorna la Config final
func (b *ConfigBuilder) Build() *Config {
	return b.config
}
