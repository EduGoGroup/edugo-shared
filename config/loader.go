package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Loader carga configuración desde archivos YAML y variables de entorno
type Loader struct {
	configPaths      []string
	configName       string
	configType       string
	envPrefix        string
	environment      string
	explicitBindings map[string]string
	defaults         map[string]interface{}
	envFiles         []string
	viper            *viper.Viper
}

// LoaderOption función de configuración para Loader
type LoaderOption func(*Loader)

// WithConfigPath establece un path de configuración (se puede usar varias veces)
func WithConfigPath(path string) LoaderOption {
	return func(l *Loader) {
		l.configPaths = append(l.configPaths, path)
	}
}

// WithConfigName establece el nombre del archivo de configuración (ej. "config")
func WithConfigName(name string) LoaderOption {
	return func(l *Loader) {
		l.configName = name
	}
}

// WithConfigType establece el tipo de archivo de configuración (yaml, json, toml...)
func WithConfigType(configType string) LoaderOption {
	return func(l *Loader) {
		l.configType = configType
	}
}

// WithEnvPrefix establece el prefijo principal para variables de entorno
func WithEnvPrefix(prefix string) LoaderOption {
	return func(l *Loader) {
		l.envPrefix = prefix
	}
}

// WithEnvironmentOverride permite fusionar un archivo secundario (ej. "config-dev.yaml") sobre el base
func WithEnvironmentOverride(env string) LoaderOption {
	return func(l *Loader) {
		l.environment = env
	}
}

// WithExplicitBindings enlaza llaves de viper con variables de entorno específicas sin prefijo
func WithExplicitBindings(bindings map[string]string) LoaderOption {
	return func(l *Loader) {
		for k, v := range bindings {
			l.explicitBindings[k] = v
		}
	}
}

// WithDefaults establece valores por defecto en memoria antes de leer configuración alguna
func WithDefaults(defaults map[string]interface{}) LoaderOption {
	return func(l *Loader) {
		for k, v := range defaults {
			l.defaults[k] = v
		}
	}
}

// WithEnvFiles permite cargar archivos .env y poblar el entorno de sistema antes que Viper actúe
func WithEnvFiles(files ...string) LoaderOption {
	return func(l *Loader) {
		l.envFiles = append(l.envFiles, files...)
	}
}

// NewLoader crea un nuevo loader de configuración con valores por defecto
func NewLoader(opts ...LoaderOption) *Loader {
	loader := &Loader{
		configPaths:      []string{"./config"},
		configName:       "config",
		configType:       "yaml",
		envPrefix:        "",
		explicitBindings: make(map[string]string),
		defaults:         make(map[string]interface{}),
		envFiles:         make([]string, 0),
	}

	for _, opt := range opts {
		opt(loader)
	}

	return loader
}

// newViper crea y configura una instancia local de viper con todas las opciones del Loader
func (l *Loader) newViper() *viper.Viper {
	v := viper.New()

	v.SetConfigType(l.configType)
	v.SetConfigName(l.configName)

	for _, path := range l.configPaths {
		v.AddConfigPath(path)
	}

	for key, val := range l.defaults {
		v.SetDefault(key, val)
	}

	if l.envPrefix != "" {
		v.SetEnvPrefix(l.envPrefix)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	for viperKey, envKey := range l.explicitBindings {
		_ = v.BindEnv(viperKey, envKey)
	}

	return v
}

// Load carga la configuración y la desempaqueta en el struct destino
func (l *Loader) Load(cfg any) error {
	if len(l.envFiles) > 0 {
		_ = godotenv.Load(l.envFiles...)
	}

	v := l.newViper()

	if err := v.ReadInConfig(); err != nil {
		var configNotFoundErr viper.ConfigFileNotFoundError
		if !errors.As(err, &configNotFoundErr) {
			return fmt.Errorf("failed to read base config file: %w", err)
		}
	}

	if l.environment != "" {
		v.SetConfigName(fmt.Sprintf("%s-%s", l.configName, l.environment))
		if err := v.MergeInConfig(); err != nil {
			var configNotFoundErr viper.ConfigFileNotFoundError
			if !errors.As(err, &configNotFoundErr) {
				return fmt.Errorf("failed to merge environment config file: %w", err)
			}
		}
		v.SetConfigName(l.configName)
	}

	if err := v.Unmarshal(cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	l.viper = v
	return nil
}

// LoadFromFile carga configuración solo desde archivo, sin leer variables de entorno
func (l *Loader) LoadFromFile(cfg any) error {
	v := l.newViper()

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if l.environment != "" {
		v.SetConfigName(fmt.Sprintf("%s-%s", l.configName, l.environment))
		if err := v.MergeInConfig(); err != nil {
			var configNotFoundErr viper.ConfigFileNotFoundError
			if !errors.As(err, &configNotFoundErr) {
				return fmt.Errorf("failed to merge config file: %w", err)
			}
		}
	}

	if err := v.Unmarshal(cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	l.viper = v
	return nil
}

// Get obtiene un valor de configuración por su key
func (l *Loader) Get(key string) any {
	if l.viper == nil {
		return nil
	}
	return l.viper.Get(key)
}

// GetString obtiene un string de configuración
func (l *Loader) GetString(key string) string {
	if l.viper == nil {
		return ""
	}
	return l.viper.GetString(key)
}

// GetInt obtiene un int de configuración
func (l *Loader) GetInt(key string) int {
	if l.viper == nil {
		return 0
	}
	return l.viper.GetInt(key)
}

// GetBool obtiene un bool de configuración
func (l *Loader) GetBool(key string) bool {
	if l.viper == nil {
		return false
	}
	return l.viper.GetBool(key)
}
