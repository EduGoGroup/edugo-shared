// Package timeutil provides time helpers that enforce the EduGo time standard:
// instants are always handled in UTC (ISO-8601 with a Z suffix), while pure
// dates are handled as YYYY-MM-DD without any time zone.
package timeutil

import (
	"fmt"
	"time"
)

// dateLayout es el formato de fecha pura YYYY-MM-DD (sin zona ni hora).
const dateLayout = "2006-01-02"

// NowUTC devuelve el instante actual normalizado a UTC.
func NowUTC() time.Time {
	return time.Now().UTC()
}

// FormatISO serializa un instante como ISO-8601 en UTC con sufijo Z (RFC3339).
func FormatISO(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// ParseISO parsea un instante ISO-8601/RFC3339 y lo devuelve en UTC.
func ParseISO(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid ISO-8601 timestamp: %w", err)
	}
	return t.UTC(), nil
}

// FormatDate serializa una fecha pura como YYYY-MM-DD, sin zona.
func FormatDate(t time.Time) string {
	return t.Format(dateLayout)
}

// ParseDate parsea una fecha pura YYYY-MM-DD (medianoche UTC, sin desplazamiento de zona).
func ParseDate(s string) (time.Time, error) {
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date: %w", err)
	}
	return t, nil
}
