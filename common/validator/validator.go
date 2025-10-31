// Package validator provides input validation utilities and error collection
// for the EduGo shared library.
package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"

	"github.com/EduGoGroup/edugo-shared/common/errors"
)

var (
	// Expresión regular para validar emails (RFC 5322 simplificado)
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	// Expresión regular para validar URLs
	urlRegex = regexp.MustCompile(`^https?://[a-zA-Z0-9\-._~:/?#\[\]@!$&'()*+,;=]+$`)

	// Expresión regular para nombres (solo letras, espacios, guiones y apóstrofes)
	nameRegex = regexp.MustCompile(`^[a-zA-ZáéíóúÁÉÍÓÚñÑüÜ\s'\-]+$`)
)

// Validator proporciona métodos de validación comunes
type Validator struct {
	errors []string
}

// New crea un nuevo Validator
func New() *Validator {
	return &Validator{
		errors: []string{},
	}
}

// AddError agrega un error de validación
func (v *Validator) AddError(message string) {
	v.errors = append(v.errors, message)
}

// AddErrorf agrega un error de validación con formato
func (v *Validator) AddErrorf(format string, args ...interface{}) {
	v.errors = append(v.errors, fmt.Sprintf(format, args...))
}

// HasErrors retorna true si hay errores
func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

// GetErrors retorna todos los errores
func (v *Validator) GetErrors() []string {
	return v.errors
}

// GetError retorna un AppError con todos los errores de validación
func (v *Validator) GetError() error {
	if !v.HasErrors() {
		return nil
	}
	return errors.NewValidationError(strings.Join(v.errors, "; "))
}

// Required valida que un campo no esté vacío
func (v *Validator) Required(value, fieldName string) {
	if strings.TrimSpace(value) == "" {
		v.AddErrorf("%s is required", fieldName)
	}
}

// MinLength valida que un string tenga longitud mínima
func (v *Validator) MinLength(value string, minLength int, fieldName string) {
	if len(value) < minLength {
		v.AddErrorf("%s must be at least %d characters", fieldName, minLength)
	}
}

// MaxLength valida que un string tenga longitud máxima
func (v *Validator) MaxLength(value string, maxLength int, fieldName string) {
	if len(value) > maxLength {
		v.AddErrorf("%s must be at most %d characters", fieldName, maxLength)
	}
}

// Email valida que un string sea un email válido
func (v *Validator) Email(value, fieldName string) {
	if value != "" && !IsValidEmail(value) {
		v.AddErrorf("%s must be a valid email address", fieldName)
	}
}

// UUID valida que un string sea un UUID válido
func (v *Validator) UUID(value, fieldName string) {
	if value != "" && !IsValidUUID(value) {
		v.AddErrorf("%s must be a valid UUID", fieldName)
	}
}

// URL valida que un string sea una URL válida
func (v *Validator) URL(value, fieldName string) {
	if value != "" && !IsValidURL(value) {
		v.AddErrorf("%s must be a valid URL", fieldName)
	}
}

// InSlice valida que un valor esté en una lista de valores permitidos
func (v *Validator) InSlice(value string, allowed []string, fieldName string) {
	if value == "" {
		return
	}

	for _, a := range allowed {
		if value == a {
			return
		}
	}

	v.AddErrorf("%s must be one of: %s", fieldName, strings.Join(allowed, ", "))
}

// MinValue valida que un número sea mayor o igual a un mínimo
func (v *Validator) MinValue(value, minValue int, fieldName string) {
	if value < minValue {
		v.AddErrorf("%s must be at least %d", fieldName, minValue)
	}
}

// MaxValue valida que un número sea menor o igual a un máximo
func (v *Validator) MaxValue(value, maxValue int, fieldName string) {
	if value > maxValue {
		v.AddErrorf("%s must be at most %d", fieldName, maxValue)
	}
}

// Range valida que un número esté en un rango
func (v *Validator) Range(value, minValue, maxValue int, fieldName string) {
	if value < minValue || value > maxValue {
		v.AddErrorf("%s must be between %d and %d", fieldName, minValue, maxValue)
	}
}

// Name valida que un string sea un nombre válido
func (v *Validator) Name(value, fieldName string) {
	if value != "" && !IsValidName(value) {
		v.AddErrorf("%s must contain only letters, spaces, hyphens and apostrophes", fieldName)
	}
}

// Funciones helper independientes

// IsValidEmail verifica si un string es un email válido
func IsValidEmail(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}

// IsValidUUID verifica si un string es un UUID válido
func IsValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

// IsValidURL verifica si un string es una URL válida
func IsValidURL(url string) bool {
	return urlRegex.MatchString(url)
}

// IsValidName verifica si un string es un nombre válido
func IsValidName(name string) bool {
	if len(name) < 1 || len(name) > 100 {
		return false
	}
	return nameRegex.MatchString(name)
}

// IsEmpty verifica si un string está vacío (después de trim)
func IsEmpty(value string) bool {
	return strings.TrimSpace(value) == ""
}

// Normalize normaliza un string (trim + lowercase)
func Normalize(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}
