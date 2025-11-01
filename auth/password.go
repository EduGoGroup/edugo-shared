package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// bcryptCost define el costo computacional de bcrypt
// Cost 12 = ~250ms por hash (balance entre seguridad y UX)
const bcryptCost = 12

// maxPasswordLength es el límite de bcrypt (72 bytes)
const maxPasswordLength = 72

// HashPassword genera un hash seguro del password usando bcrypt
// El hash incluye un salt aleatorio automáticamente
// Retorna error si el password excede 72 bytes (límite de bcrypt)
func HashPassword(password string) (string, error) {
	if len(password) > maxPasswordLength {
		return "", fmt.Errorf("password exceeds maximum length of %d bytes", maxPasswordLength)
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// VerifyPassword verifica si un password coincide con su hash bcrypt
// Retorna nil si el password es correcto, error en caso contrario
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
}
