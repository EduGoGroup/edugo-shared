package gin

import (
	"errors"
	"testing"

	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/gin-gonic/gin"
)

func TestGetUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)

	// Test: Key no existe
	_, err := GetUserID(c)
	if err == nil {
		t.Error("Debe retornar error cuando user_id no existe")
	}
	if !errors.Is(err, ErrUserIDNotFound) {
		t.Errorf("Error debe ser ErrUserIDNotFound, got: %v", err)
	}

	// Test: Key existe con valor correcto
	c.Set(ContextKeyUserID, "user-123")
	userID, err := GetUserID(c)
	if err != nil {
		t.Errorf("No debe haber error con valor correcto: %v", err)
	}
	if userID != "user-123" {
		t.Errorf("Expected 'user-123', got '%s'", userID)
	}

	// Test: Key existe con tipo incorrecto
	c.Set(ContextKeyUserID, 12345) // int en lugar de string
	_, err = GetUserID(c)
	if err == nil {
		t.Error("Debe retornar error cuando el tipo es incorrecto")
	}
	if !errors.Is(err, ErrInvalidType) {
		t.Errorf("Error debe ser ErrInvalidType, got: %v", err)
	}
}

func TestMustGetUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)

	// Test: Panic cuando no existe
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustGetUserID debe hacer panic cuando user_id no existe")
		}
	}()

	MustGetUserID(c)
}

func TestMustGetUserID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)

	c.Set(ContextKeyUserID, "user-456")

	// No debe hacer panic
	userID := MustGetUserID(c)
	if userID != "user-456" {
		t.Errorf("Expected 'user-456', got '%s'", userID)
	}
}

func TestGetEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)

	// Test: Key no existe
	_, err := GetEmail(c)
	if !errors.Is(err, ErrEmailNotFound) {
		t.Errorf("Error debe ser ErrEmailNotFound, got: %v", err)
	}

	// Test: Key existe con valor correcto
	c.Set(ContextKeyEmail, "jhoan@edugo.com")
	email, err := GetEmail(c)
	if err != nil {
		t.Errorf("No debe haber error: %v", err)
	}
	if email != "jhoan@edugo.com" {
		t.Errorf("Expected 'jhoan@edugo.com', got '%s'", email)
	}

	// Test: Tipo incorrecto
	c.Set(ContextKeyEmail, 999)
	_, err = GetEmail(c)
	if !errors.Is(err, ErrInvalidType) {
		t.Errorf("Error debe ser ErrInvalidType, got: %v", err)
	}
}

func TestGetRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)

	// Test: Key no existe
	_, err := GetRole(c)
	if !errors.Is(err, ErrRoleNotFound) {
		t.Errorf("Error debe ser ErrRoleNotFound, got: %v", err)
	}

	// Test: Valores válidos
	roles := []string{"student", "teacher", "admin", "guardian"}
	for _, expectedRole := range roles {
		c.Set(ContextKeyRole, expectedRole)
		role, err := GetRole(c)
		if err != nil {
			t.Errorf("Error con role '%s': %v", expectedRole, err)
		}
		if role != expectedRole {
			t.Errorf("Expected '%s', got '%s'", expectedRole, role)
		}
	}
}

func TestGetClaims(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)

	// Test: Key no existe
	_, err := GetClaims(c)
	if !errors.Is(err, ErrClaimsNotFound) {
		t.Errorf("Error debe ser ErrClaimsNotFound, got: %v", err)
	}

	// Test: Claims válidos
	expectedClaims := &auth.Claims{
		UserID: "user-789",
		Email:  "claims@test.com",
		Role:   "teacher",
	}

	c.Set(ContextKeyClaims, expectedClaims)
	claims, err := GetClaims(c)
	if err != nil {
		t.Errorf("No debe haber error: %v", err)
	}

	if claims.UserID != expectedClaims.UserID {
		t.Errorf("UserID: expected '%s', got '%s'", expectedClaims.UserID, claims.UserID)
	}
	if claims.Email != expectedClaims.Email {
		t.Errorf("Email: expected '%s', got '%s'", expectedClaims.Email, claims.Email)
	}
	if claims.Role != expectedClaims.Role {
		t.Errorf("Role: expected '%s', got '%s'", expectedClaims.Role, claims.Role)
	}

	// Test: Tipo incorrecto
	c.Set(ContextKeyClaims, "not-claims-object")
	_, err = GetClaims(c)
	if !errors.Is(err, ErrInvalidType) {
		t.Errorf("Error debe ser ErrInvalidType, got: %v", err)
	}
}

func TestMustGetters_Panic(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name   string
		getter func(*gin.Context)
	}{
		{"MustGetUserID", func(c *gin.Context) { MustGetUserID(c) }},
		{"MustGetEmail", func(c *gin.Context) { MustGetEmail(c) }},
		{"MustGetRole", func(c *gin.Context) { MustGetRole(c) }},
		{"MustGetClaims", func(c *gin.Context) { MustGetClaims(c) }},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(nil)

			defer func() {
				if r := recover(); r == nil {
					t.Errorf("%s debe hacer panic cuando la key no existe", tc.name)
				}
			}()

			tc.getter(c)
		})
	}
}

func TestMustGetters_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)

	// Setup contexto con valores
	c.Set(ContextKeyUserID, "user-999")
	c.Set(ContextKeyEmail, "success@test.com")
	c.Set(ContextKeyRole, "admin")
	c.Set(ContextKeyClaims, &auth.Claims{
		UserID: "user-999",
		Email:  "success@test.com",
		Role:   "admin",
	})

	// Test: MustGetUserID
	userID := MustGetUserID(c)
	if userID != "user-999" {
		t.Errorf("Expected 'user-999', got '%s'", userID)
	}

	// Test: MustGetEmail
	email := MustGetEmail(c)
	if email != "success@test.com" {
		t.Errorf("Expected 'success@test.com', got '%s'", email)
	}

	// Test: MustGetRole
	role := MustGetRole(c)
	if role != "admin" {
		t.Errorf("Expected 'admin', got '%s'", role)
	}

	// Test: MustGetClaims
	claims := MustGetClaims(c)
	if claims.UserID != "user-999" {
		t.Errorf("Claims UserID incorrect: %s", claims.UserID)
	}
}

func TestContextKeys_Constants(t *testing.T) {
	// Verificar que las constantes tienen valores esperados
	if ContextKeyUserID != "user_id" {
		t.Errorf("ContextKeyUserID debe ser 'user_id', got '%s'", ContextKeyUserID)
	}
	if ContextKeyEmail != "email" {
		t.Errorf("ContextKeyEmail debe ser 'email', got '%s'", ContextKeyEmail)
	}
	if ContextKeyRole != "role" {
		t.Errorf("ContextKeyRole debe ser 'role', got '%s'", ContextKeyRole)
	}
	if ContextKeyClaims != "jwt_claims" {
		t.Errorf("ContextKeyClaims debe ser 'jwt_claims', got '%s'", ContextKeyClaims)
	}
}
