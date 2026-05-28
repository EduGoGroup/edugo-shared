package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"
)

// RefreshToken representa un token de refrescamiento
type RefreshToken struct {
	Token     string    // Token en texto plano (se retorna al cliente)
	TokenHash string    // SHA256 del token (se guarda en BD)
	ExpiresAt time.Time // Timestamp de expiración
}

// GenerateRefreshToken genera un refresh token criptográficamente seguro
// Usa crypto/rand para generar 32 bytes aleatorios
// Retorna el token en base64 y su hash SHA256
func GenerateRefreshToken(ttl time.Duration) (*RefreshToken, error) {
	// Generar 32 bytes aleatorios usando crypto/rand (seguro)
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Codificar en base64 URL-safe (token que se retorna al cliente)
	token := base64.URLEncoding.EncodeToString(bytes)

	// Generar hash SHA256 del token (lo que se guarda en BD)
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	return &RefreshToken{
		Token:     token,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(ttl),
	}, nil
}

// HashToken genera el hash SHA256 de un token
// Útil para verificar un token recibido contra la base de datos
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
