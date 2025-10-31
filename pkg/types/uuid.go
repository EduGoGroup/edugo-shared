// Package types provides custom type definitions and utilities,
// including UUID handling for the EduGo shared library.
package types

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
)

const (
	// MinJSONLength is the minimum length for a JSON-encoded UUID string (with quotes)
	MinJSONLength = 2
)

// UUID es un wrapper alrededor de google/uuid con métodos de serialización
type UUID struct {
	uuid.UUID
}

// NewUUID genera un nuevo UUID v4
func NewUUID() UUID {
	return UUID{uuid.New()}
}

// ParseUUID parsea un string a UUID
func ParseUUID(s string) (UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return UUID{}, fmt.Errorf("invalid UUID: %w", err)
	}
	return UUID{id}, nil
}

// MustParseUUID parsea un string a UUID, panic si hay error
func MustParseUUID(s string) UUID {
	id, err := ParseUUID(s)
	if err != nil {
		panic(err)
	}
	return id
}

// String retorna la representación en string del UUID
func (u UUID) String() string {
	return u.UUID.String()
}

// IsZero verifica si el UUID es el valor cero
func (u UUID) IsZero() bool {
	return u.UUID == uuid.Nil
}

// MarshalJSON implementa json.Marshaler
func (u UUID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + u.String() + `"`), nil
}

// UnmarshalJSON implementa json.Unmarshaler
func (u *UUID) UnmarshalJSON(data []byte) error {
	// Remover comillas
	s := string(data)
	if len(s) < MinJSONLength {
		return fmt.Errorf("invalid UUID JSON")
	}
	s = s[1 : len(s)-1]

	parsed, err := ParseUUID(s)
	if err != nil {
		return err
	}

	*u = parsed
	return nil
}

// Value implementa driver.Valuer para PostgreSQL
func (u UUID) Value() (driver.Value, error) {
	return u.String(), nil
}

// Scan implementa sql.Scanner para PostgreSQL
func (u *UUID) Scan(value interface{}) error {
	if value == nil {
		*u = UUID{}
		return nil
	}

	switch v := value.(type) {
	case string:
		parsed, err := ParseUUID(v)
		if err != nil {
			return err
		}
		*u = parsed
		return nil
	case []byte:
		parsed, err := ParseUUID(string(v))
		if err != nil {
			return err
		}
		*u = parsed
		return nil
	default:
		return fmt.Errorf("cannot scan %T into UUID", value)
	}
}
