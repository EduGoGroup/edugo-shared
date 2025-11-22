package auth

import (
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "secreto123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error al hashear password: %v", err)
	}

	// Verificar que el hash no está vacío
	if hash == "" {
		t.Error("Hash no debe estar vacío")
	}

	// Verificar que el hash es diferente al password original
	if hash == password {
		t.Error("Hash no debe ser igual al password en texto plano")
	}

	// Verificar que la longitud es apropiada para bcrypt (~60 caracteres)
	if len(hash) < 50 {
		t.Errorf("Hash muy corto para bcrypt, longitud: %d", len(hash))
	}

	// Verificar que empieza con el prefijo de bcrypt
	if !strings.HasPrefix(hash, "$2a$") && !strings.HasPrefix(hash, "$2b$") {
		t.Errorf("Hash debe empezar con prefijo bcrypt ($2a$ o $2b$), pero empieza con: %s", hash[:4])
	}
}

func TestVerifyPassword_Correct(t *testing.T) {
	password := "miPasswordSeguro123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error al hashear: %v", err)
	}

	// Verificar password correcto
	err = VerifyPassword(hash, password)
	if err != nil {
		t.Errorf("Password correcto debe verificar sin error, pero obtuvo: %v", err)
	}
}

func TestVerifyPassword_Incorrect(t *testing.T) {
	password := "passwordCorrecto"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error al hashear: %v", err)
	}

	// Verificar password incorrecto
	err = VerifyPassword(hash, "passwordIncorrecto")
	if err == nil {
		t.Error("Password incorrecto debe fallar la verificación")
	}
}

func TestHashUniqueness(t *testing.T) {
	password := "mismoPa$$word123"

	// Generar dos hashes del mismo password
	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)

	if err1 != nil || err2 != nil {
		t.Fatalf("Error al hashear: %v, %v", err1, err2)
	}

	// Los hashes deben ser diferentes porque bcrypt usa salt aleatorio
	if hash1 == hash2 {
		t.Error("Hashes del mismo password deben ser únicos debido al salt aleatorio de bcrypt")
	}

	// Pero ambos deben verificar correctamente contra el password original
	if err := VerifyPassword(hash1, password); err != nil {
		t.Error("Hash1 debe verificar correctamente")
	}
	if err := VerifyPassword(hash2, password); err != nil {
		t.Error("Hash2 debe verificar correctamente")
	}
}

func TestHashPassword_EmptyPassword(t *testing.T) {
	// bcrypt acepta passwords vacíos, pero es buena idea testear
	hash, err := HashPassword("")
	if err != nil {
		t.Fatalf("Error al hashear password vacío: %v", err)
	}

	// Debe verificar correctamente
	err = VerifyPassword(hash, "")
	if err != nil {
		t.Error("Password vacío debe verificar correctamente")
	}

	// No debe verificar con password no vacío
	err = VerifyPassword(hash, "algo")
	if err == nil {
		t.Error("Verificación debe fallar con password diferente")
	}
}

func TestHashPassword_LongPassword(t *testing.T) {
	// Password muy largo (>72 bytes excede límite de bcrypt)
	longPassword := strings.Repeat("a", 100)

	// Debe retornar error porque excede 72 bytes
	_, err := HashPassword(longPassword)
	if err == nil {
		t.Error("HashPassword debe retornar error para password >72 bytes")
	}

	// Verificar que el error menciona el límite
	if !strings.Contains(err.Error(), "72") {
		t.Errorf("Error debe mencionar el límite de 72 bytes, pero dice: %v", err)
	}
}

func TestHashPassword_MaxLength(t *testing.T) {
	// Password exactamente en el límite (72 bytes)
	maxPassword := strings.Repeat("a", 72)

	hash, err := HashPassword(maxPassword)
	if err != nil {
		t.Fatalf("Password de 72 bytes debe ser aceptado: %v", err)
	}

	// Debe verificar correctamente
	err = VerifyPassword(hash, maxPassword)
	if err != nil {
		t.Error("Password de 72 bytes debe verificar correctamente")
	}
}

func TestVerifyPassword_InvalidHash(t *testing.T) {
	// Hash inválido (no es bcrypt)
	invalidHash := "not-a-valid-bcrypt-hash"

	err := VerifyPassword(invalidHash, "password")
	if err == nil {
		t.Error("Verificación debe fallar con hash inválido")
	}
}

// Benchmark para medir performance de hash
func BenchmarkHashPassword(b *testing.B) {
	password := "benchmarkPassword123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = HashPassword(password)
	}
}

// Benchmark para medir performance de verificación
func BenchmarkVerifyPassword(b *testing.B) {
	password := "benchmarkPassword123"
	hash, _ := HashPassword(password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = VerifyPassword(hash, password)
	}
}
