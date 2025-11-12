package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Loader carga configuración desde archivos YAML y variables de entorno
type Loader struct {
	configPath string
	configName string
	configType string
	envPrefix  string
}

// LoaderOption función de configuración para Loader
type LoaderOption func(*Loader)

// WithConfigPath establece el path de configuración
func WithConfigPath(path string) LoaderOption {
	return func(l *Loader) {
		l.configPath = path
	}
}

// WithConfigName establece el nombre del archivo de configuración
func WithConfigName(name string) LoaderOption {
	return func(l *Loader) {
		l.configName = name
	}
}

// WithConfigType establece el tipo de archivo de configuración
func WithConfigType(configType string) LoaderOption {
	return func(l *Loader) {
		l.configType = configType
	}
}

// WithEnvPrefix establece el prefijo para variables de entorno
func WithEnvPrefix(prefix string) LoaderOption {
	return func(l *Loader) {
		l.envPrefix = prefix
	}
}

// NewLoader crea un nuevo loader de configuración con valores por defecto
func NewLoader(opts ...LoaderOption) *Loader {
	loader := &Loader{
		configPath: "./config",
		configName: "config",
		configType: "yaml",
		envPrefix:  "",
	}

	for _, opt := range opts {
		opt(loader)
	}

	return loader
}

// Load carga la configuración y la desempaqueta en el struct destino
func (l *Loader) Load(cfg interface{}) error {
	// Configurar viper
	viper.AddConfigPath(l.configPath)
	viper.SetConfigName(l.configName)
	viper.SetConfigType(l.configType)

	// Configurar variables de entorno
	if l.envPrefix != "" {
		viper.SetEnvPrefix(l.envPrefix)
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Leer archivo de configuración
	if err := viper.ReadInConfig(); err != nil {
		// Si el archivo no existe, continuar solo con env vars
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Desempaquetar en el struct
	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// LoadFromFile carga configuración solo desde archivo (sin env vars)
func (l *Loader) LoadFromFile(cfg interface{}) error {
	v := viper.New()
	v.AddConfigPath(l.configPath)
	v.SetConfigName(l.configName)
	v.SetConfigType(l.configType)

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := v.Unmarshal(cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// Get obtiene un valor de configuración por su key
func (l *Loader) Get(key string) interface{} {
	return viper.Get(key)
}

// GetString obtiene un string de configuración
func (l *Loader) GetString(key string) string {
	return viper.GetString(key)
}

// GetInt obtiene un int de configuración
func (l *Loader) GetInt(key string) int {
	return viper.GetInt(key)
}

// GetBool obtiene un bool de configuración
func (l *Loader) GetBool(key string) bool {
	return viper.GetBool(key)
}
