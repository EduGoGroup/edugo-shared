package config

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Validator valida configuración usando tags struct
type Validator struct {
	validate *validator.Validate
}

// NewValidator crea un nuevo validador de configuración
func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

// Validate valida un struct de configuración
func (v *Validator) Validate(cfg interface{}) error {
	if err := v.validate.Struct(cfg); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return NewValidationError(validationErrors)
		}
		return fmt.Errorf("config validation failed: %w", err)
	}
	return nil
}

// ValidateField valida un campo específico
func (v *Validator) ValidateField(field interface{}, tag string) error {
	if err := v.validate.Var(field, tag); err != nil {
		return fmt.Errorf("field validation failed: %w", err)
	}
	return nil
}

// ValidationError error de validación con detalles
type ValidationError struct {
	Errors []FieldError
}

// FieldError error de un campo específico
type FieldError struct {
	Field   string
	Tag     string
	Value   interface{}
	Message string
}

// NewValidationError crea un ValidationError desde validator.ValidationErrors
func NewValidationError(errs validator.ValidationErrors) *ValidationError {
	fieldErrors := make([]FieldError, 0, len(errs))

	for _, err := range errs {
		fieldErrors = append(fieldErrors, FieldError{
			Field:   err.Field(),
			Tag:     err.Tag(),
			Value:   err.Value(),
			Message: buildErrorMessage(err),
		})
	}

	return &ValidationError{
		Errors: fieldErrors,
	}
}

// Error implementa la interfaz error
func (e *ValidationError) Error() string {
	if len(e.Errors) == 0 {
		return "config validation failed"
	}

	msg := fmt.Sprintf("config validation failed with %d error(s):", len(e.Errors))
	for _, fieldErr := range e.Errors {
		msg += fmt.Sprintf("\n  - %s: %s", fieldErr.Field, fieldErr.Message)
	}

	return msg
}

// buildErrorMessage construye un mensaje de error amigable
func buildErrorMessage(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()
	param := err.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s", field, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s", field, param)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, param)
	case "email":
		return fmt.Sprintf("%s must be a valid email", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	default:
		return fmt.Sprintf("%s failed validation: %s", field, tag)
	}
}
