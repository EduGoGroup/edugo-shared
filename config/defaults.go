package config

import (
	"reflect"
	"strconv"
	"time"
)

var durationType = reflect.TypeFor[time.Duration]()

// ExtractDefaults construye un mapa de defaults leyendo los tags `default` y `mapstructure`
// de un struct. Esto permite declarar los valores por defecto junto a la definición del campo
// y que Viper los descubra automáticamente vía AutomaticEnv.
//
// Campos sin tag `default` se registran con su zero value para que Viper conozca la clave
// y pueda resolver la variable de entorno correspondiente.
//
// Ejemplo:
//
//	type ServerConfig struct {
//	    Port int    `mapstructure:"port" default:"8080"`
//	    Host string `mapstructure:"host" default:"0.0.0.0"`
//	}
//
//	ExtractDefaults(ServerConfig{}) // → {"port": int64(8080), "host": "0.0.0.0"}
func ExtractDefaults(v any) map[string]any {
	defaults := make(map[string]any)
	walkStruct(reflect.TypeOf(v), "", defaults)
	return defaults
}

func walkStruct(t reflect.Type, prefix string, out map[string]any) {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}

	for i := range t.NumField() {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		key := field.Tag.Get("mapstructure")
		if key == "" || key == "-" {
			continue
		}

		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		ft := field.Type
		if ft.Kind() == reflect.Pointer {
			ft = ft.Elem()
		}

		// Recurse into nested structs (time.Duration es int64, no struct)
		if ft.Kind() == reflect.Struct {
			walkStruct(ft, fullKey, out)
			continue
		}

		if val, ok := field.Tag.Lookup("default"); ok {
			out[fullKey] = parseTagValue(val, ft)
		} else {
			out[fullKey] = zeroValue(ft)
		}
	}
}

func parseTagValue(val string, t reflect.Type) any {
	if t == durationType {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
		return val
	}

	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if n, err := strconv.ParseInt(val, 10, 64); err == nil {
			return n
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if n, err := strconv.ParseUint(val, 10, 64); err == nil {
			return n
		}
	case reflect.Float32, reflect.Float64:
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	case reflect.Bool:
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	}

	return val
}

func zeroValue(t reflect.Type) any {
	if t == durationType {
		return time.Duration(0)
	}
	return reflect.Zero(t).Interface()
}
