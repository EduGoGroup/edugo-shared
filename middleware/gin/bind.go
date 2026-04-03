package gin

import (
	"errors"
	"fmt"
	"reflect"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	sharedErrors "github.com/EduGoGroup/edugo-shared/common/errors"
)

// BindJSON hace ShouldBindJSON extrayendo errores de campo detallados.
// Retorna ValidationError de edugo-shared/common/errors con campo-por-campo.
func BindJSON(c *gin.Context, v any) error {
	if err := c.ShouldBindJSON(v); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			fields := make(map[string]string, len(ve))
			for _, fe := range ve {
				fields[ToSnakeCase(fe.Field())] = ValidationMessage(fe)
			}
			return sharedErrors.NewValidationErrorWithFields("validation failed", fields)
		}
		return sharedErrors.NewValidationError("invalid request body")
	}
	return nil
}

// ToSnakeCase convierte CamelCase a snake_case.
func ToSnakeCase(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	result := make([]rune, 0, len(runes))
	for i, r := range runes {
		if i > 0 && unicode.IsUpper(r) {
			prev := runes[i-1]
			hasNext := i+1 < len(runes)
			var next rune
			if hasNext {
				next = runes[i+1]
			}
			if unicode.IsLower(prev) || unicode.IsDigit(prev) ||
				(unicode.IsUpper(prev) && hasNext && unicode.IsLower(next)) {
				result = append(result, '_')
			}
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

// ValidationMessage genera un mensaje legible para un error de validacion.
// Distingue entre length (string, slice, array, map) y value (numeros) para min/max.
func ValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "field is required"
	case "email":
		return "invalid email format"
	case "min":
		switch fe.Kind() {
		case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
			return fmt.Sprintf("minimum length is %s", fe.Param())
		default:
			return fmt.Sprintf("minimum value is %s", fe.Param())
		}
	case "max":
		switch fe.Kind() {
		case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
			return fmt.Sprintf("maximum length is %s", fe.Param())
		default:
			return fmt.Sprintf("maximum value is %s", fe.Param())
		}
	case "uuid":
		return "must be a valid UUID"
	case "oneof":
		return fmt.Sprintf("must be one of: %s", fe.Param())
	default:
		return fmt.Sprintf("failed validation '%s'", fe.Tag())
	}
}
