package auth

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateRefreshToken(t *testing.T) {
	ttl := 7 * 24 * time.Hour

	token, err := GenerateRefreshToken(ttl)
	if err != nil {
		t.Fatalf("Error al generar token: %v", err)
	}

	// Verificar que el token no está vacío
	if token.Token == "" {
		t.Error("Token no debe estar vacío")
	}

	// Verificar que el hash no está vacío
	if token.TokenHash == "" {
		t.Error("TokenHash no debe estar vacío")
	}

	// Verificar longitud del token (base64 de 32 bytes ≈ 44 chars)
	if len(token.Token) < 40 {
		t.Errorf("Token muy corto, longitud: %d", len(token.Token))
	}

	// Verificar longitud del hash (SHA256 = 64 hex chars)
	if len(token.TokenHash) != 64 {
		t.Errorf("TokenHash debe tener 64 chars, tiene %d", len(token.TokenHash))
	}

	// Verificar que ExpiresAt está en el futuro
	if !token.ExpiresAt.After(time.Now()) {
		t.Error("ExpiresAt debe estar en el futuro")
	}

	// Verificar que ExpiresAt es aproximadamente TTL en el futuro
	expectedExpiry := time.Now().Add(ttl)
	diff := token.ExpiresAt.Sub(expectedExpiry).Abs()
	if diff > time.Second {
		t.Errorf("ExpiresAt debería ser ~%v en el futuro, diferencia: %v", ttl, diff)
	}
}

func TestTokenUniqueness(t *testing.T) {
	token1, err1 := GenerateRefreshToken(time.Hour)
	token2, err2 := GenerateRefreshToken(time.Hour)

	if err1 != nil || err2 != nil {
		t.Fatalf("Error al generar tokens: %v, %v", err1, err2)
	}

	// Tokens deben ser únicos
	if token1.Token == token2.Token {
		t.Error("Tokens generados deben ser únicos (crypto/rand debe generar valores diferentes)")
	}

	// Hashes también deben ser únicos
	if token1.TokenHash == token2.TokenHash {
		t.Error("TokenHashes deben ser únicos")
	}

	// Tokens no deben compartir prefijos largos (verdadera aleatoriedad)
	if len(token1.Token) > 10 && token1.Token[:10] == token2.Token[:10] {
		t.Error("Tokens no deben compartir prefijos largos (puede indicar problema de aleatoriedad)")
	}
}

func TestHashToken(t *testing.T) {
	originalToken := "test-token-123-abc-xyz"

	hash1 := HashToken(originalToken)
	hash2 := HashToken(originalToken)

	// Mismo token debe generar mismo hash (determinístico)
	if hash1 != hash2 {
		t.Error("Mismo token debe generar mismo hash siempre")
	}

	// Verificar longitud (SHA256 = 64 hex chars)
	if len(hash1) != 64 {
		t.Errorf("Hash debe tener 64 caracteres, tiene %d", len(hash1))
	}

	// Verificar que es hexadecimal
	for _, char := range hash1 {
		if !strings.ContainsRune("0123456789abcdef", char) {
			t.Errorf("Hash debe ser hexadecimal, encontrado caracter: %c", char)
			break
		}
	}

	// Diferente token debe generar diferente hash
	hash3 := HashToken("different-token-456")
	if hash1 == hash3 {
		t.Error("Tokens diferentes deben generar hashes diferentes")
	}
}

func TestHashConsistency(t *testing.T) {
	token, err := GenerateRefreshToken(time.Hour)
	if err != nil {
		t.Fatalf("Error al generar token: %v", err)
	}

	// Hash generado debe coincidir con HashToken()
	manualHash := HashToken(token.Token)
	if manualHash != token.TokenHash {
		t.Error("HashToken() debe generar el mismo hash que GenerateRefreshToken()")
	}
}

func TestGenerateRefreshToken_MultipleTTLs(t *testing.T) {
	testCases := []struct {
		name string
		ttl  time.Duration
	}{
		{"1 hora", time.Hour},
		{"24 horas", 24 * time.Hour},
		{"7 días", 7 * 24 * time.Hour},
		{"30 días", 30 * 24 * time.Hour},
		{"1 minuto", time.Minute},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := GenerateRefreshToken(tc.ttl)
			if err != nil {
				t.Fatalf("Error con TTL %v: %v", tc.ttl, err)
			}

			// Verificar que expira aproximadamente en TTL
			expectedExpiry := time.Now().Add(tc.ttl)
			diff := token.ExpiresAt.Sub(expectedExpiry).Abs()
			if diff > time.Second {
				t.Errorf("ExpiresAt incorrecto para TTL %v, diferencia: %v", tc.ttl, diff)
			}
		})
	}
}

func TestHashToken_EmptyToken(t *testing.T) {
	hash := HashToken("")

	// Hash de string vacío debe ser válido (aunque no recomendable)
	if len(hash) != 64 {
		t.Errorf("Hash de token vacío debe tener 64 chars, tiene %d", len(hash))
	}

	// Debe ser determinístico
	hash2 := HashToken("")
	if hash != hash2 {
		t.Error("Hashes de tokens vacíos deben ser iguales")
	}
}

func TestHashToken_SpecialCharacters(t *testing.T) {
	testTokens := []string{
		"token-with-dashes",
		"token_with_underscores",
		"token.with.dots",
		"token+with+plus",
		"token/with/slashes",
		"token=with=equals",
		"token con espacios",
		"token\ncon\nnewlines",
		"token🔐con🎯emojis",
	}

	hashes := make(map[string]string)

	for _, token := range testTokens {
		hash := HashToken(token)

		// Verificar longitud
		if len(hash) != 64 {
			t.Errorf("Hash de '%s' debe tener 64 chars, tiene %d", token, len(hash))
		}

		// Verificar unicidad
		if existingToken, exists := hashes[hash]; exists {
			t.Errorf("Colisión de hash: '%s' y '%s' generan mismo hash", token, existingToken)
		}
		hashes[hash] = token
	}

	// Verificar que generamos N hashes únicos
	if len(hashes) != len(testTokens) {
		t.Errorf("Deberían haber %d hashes únicos, pero hay %d", len(testTokens), len(hashes))
	}
}

// Benchmark para medir performance de generación de tokens
func BenchmarkGenerateRefreshToken(b *testing.B) {
	ttl := 7 * 24 * time.Hour

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GenerateRefreshToken(ttl)
		if err != nil {
			b.Fatalf("Error en benchmark: %v", err)
		}
	}
}

// Benchmark para medir performance de hash de tokens
func BenchmarkHashToken(b *testing.B) {
	token := "sample-refresh-token-abc123xyz"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = HashToken(token)
	}
}
