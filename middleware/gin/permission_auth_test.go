package gin

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRouter(middleware gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	return router
}

func createTestContextWithClaims(claims *auth.Claims) *gin.Context {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	if claims != nil {
		c.Set(ContextKeyClaims, claims)
	}
	return c
}

func TestRequirePermission(t *testing.T) {
	t.Run("permite acceso con permiso correcto", func(t *testing.T) {
		middleware := RequirePermission(enum.PermissionUsersRead)
		router := setupTestRouter(middleware)

		claims := &auth.Claims{
			UserID: "user-123",
			Email:  "test@edugo.com",
			ActiveContext: &auth.UserContext{
				RoleID:      "role-teacher",
				RoleName:    "Teacher",
				Permissions: []string{"users:read", "materials:create"},
			},
		}

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set(ContextKeyClaims, claims)

		router.ServeHTTP(w, req)

		// Como no se ejecutó el middleware correctamente, lo hacemos manualmente
		middleware(c)

		assert.False(t, c.IsAborted())
	})

	t.Run("rechaza acceso sin el permiso requerido", func(t *testing.T) {
		c := createTestContextWithClaims(&auth.Claims{
			UserID: "user-123",
			Email:  "test@edugo.com",
			ActiveContext: &auth.UserContext{
				RoleID:      "role-student",
				RoleName:    "Student",
				Permissions: []string{"materials:read"},
			},
		})

		middleware := RequirePermission(enum.PermissionUsersCreate)
		middleware(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusForbidden, c.Writer.Status())
	})

	t.Run("rechaza acceso sin claims en contexto", func(t *testing.T) {
		c := createTestContextWithClaims(nil)

		middleware := RequirePermission(enum.PermissionUsersRead)
		middleware(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusUnauthorized, c.Writer.Status())
	})

	t.Run("rechaza acceso sin ActiveContext", func(t *testing.T) {
		c := createTestContextWithClaims(&auth.Claims{
			UserID: "user-123",
			Email:  "test@edugo.com",
			// ActiveContext es nil
		})

		middleware := RequirePermission(enum.PermissionUsersRead)
		middleware(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusForbidden, c.Writer.Status())
	})

	t.Run("rechaza acceso con claims de tipo incorrecto", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set(ContextKeyClaims, "invalid-claims-type")

		middleware := RequirePermission(enum.PermissionUsersRead)
		middleware(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
	})

	t.Run("permite acceso cuando el usuario tiene el permiso entre varios", func(t *testing.T) {
		c := createTestContextWithClaims(&auth.Claims{
			UserID: "user-123",
			Email:  "test@edugo.com",
			ActiveContext: &auth.UserContext{
				RoleID:   "role-admin",
				RoleName: "Admin",
				Permissions: []string{
					"users:create",
					"users:read",
					"users:update",
					"users:delete",
					"schools:manage",
				},
			},
		})

		middleware := RequirePermission(enum.PermissionUsersUpdate)
		middleware(c)

		assert.False(t, c.IsAborted())
	})
}

func TestRequireAnyPermission(t *testing.T) {
	t.Run("permite acceso con uno de los permisos requeridos", func(t *testing.T) {
		c := createTestContextWithClaims(&auth.Claims{
			UserID: "user-123",
			Email:  "test@edugo.com",
			ActiveContext: &auth.UserContext{
				RoleID:      "role-teacher",
				RoleName:    "Teacher",
				Permissions: []string{"materials:read", "materials:create"},
			},
		})

		middleware := RequireAnyPermission(
			enum.PermissionMaterialsCreate,
			enum.PermissionMaterialsUpdate,
			enum.PermissionMaterialsDelete,
		)
		middleware(c)

		assert.False(t, c.IsAborted())
	})

	t.Run("permite acceso con múltiples permisos requeridos", func(t *testing.T) {
		c := createTestContextWithClaims(&auth.Claims{
			UserID: "user-123",
			Email:  "test@edugo.com",
			ActiveContext: &auth.UserContext{
				RoleID:   "role-admin",
				RoleName: "Admin",
				Permissions: []string{
					"materials:create",
					"materials:update",
					"materials:delete",
				},
			},
		})

		middleware := RequireAnyPermission(
			enum.PermissionMaterialsCreate,
			enum.PermissionMaterialsUpdate,
		)
		middleware(c)

		assert.False(t, c.IsAborted())
	})

	t.Run("rechaza acceso sin ninguno de los permisos requeridos", func(t *testing.T) {
		c := createTestContextWithClaims(&auth.Claims{
			UserID: "user-123",
			Email:  "test@edugo.com",
			ActiveContext: &auth.UserContext{
				RoleID:      "role-student",
				RoleName:    "Student",
				Permissions: []string{"materials:read"},
			},
		})

		middleware := RequireAnyPermission(
			enum.PermissionMaterialsCreate,
			enum.PermissionMaterialsUpdate,
			enum.PermissionMaterialsDelete,
		)
		middleware(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusForbidden, c.Writer.Status())
	})

	t.Run("rechaza acceso sin claims", func(t *testing.T) {
		c := createTestContextWithClaims(nil)

		middleware := RequireAnyPermission(enum.PermissionUsersRead)
		middleware(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusUnauthorized, c.Writer.Status())
	})

	t.Run("rechaza acceso sin ActiveContext", func(t *testing.T) {
		c := createTestContextWithClaims(&auth.Claims{
			UserID: "user-123",
			Email:  "test@edugo.com",
		})

		middleware := RequireAnyPermission(enum.PermissionUsersRead)
		middleware(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusForbidden, c.Writer.Status())
	})

	t.Run("permite acceso con el primer permiso de la lista", func(t *testing.T) {
		c := createTestContextWithClaims(&auth.Claims{
			UserID: "user-123",
			Email:  "test@edugo.com",
			ActiveContext: &auth.UserContext{
				RoleID:      "role-teacher",
				RoleName:    "Teacher",
				Permissions: []string{"users:read"},
			},
		})

		middleware := RequireAnyPermission(
			enum.PermissionUsersRead,
			enum.PermissionUsersCreate,
		)
		middleware(c)

		assert.False(t, c.IsAborted())
	})

	t.Run("permite acceso con el último permiso de la lista", func(t *testing.T) {
		c := createTestContextWithClaims(&auth.Claims{
			UserID: "user-123",
			Email:  "test@edugo.com",
			ActiveContext: &auth.UserContext{
				RoleID:      "role-teacher",
				RoleName:    "Teacher",
				Permissions: []string{"users:delete"},
			},
		})

		middleware := RequireAnyPermission(
			enum.PermissionUsersRead,
			enum.PermissionUsersCreate,
			enum.PermissionUsersDelete,
		)
		middleware(c)

		assert.False(t, c.IsAborted())
	})
}

func TestRequireAllPermissions(t *testing.T) {
	t.Run("permite acceso con todos los permisos requeridos", func(t *testing.T) {
		c := createTestContextWithClaims(&auth.Claims{
			UserID: "user-123",
			Email:  "test@edugo.com",
			ActiveContext: &auth.UserContext{
				RoleID:   "role-admin",
				RoleName: "Admin",
				Permissions: []string{
					"users:create",
					"users:read",
					"users:update",
					"users:delete",
				},
			},
		})

		middleware := RequireAllPermissions(
			enum.PermissionUsersCreate,
			enum.PermissionUsersRead,
			enum.PermissionUsersUpdate,
		)
		middleware(c)

		assert.False(t, c.IsAborted())
	})

	t.Run("rechaza acceso si falta un permiso", func(t *testing.T) {
		c := createTestContextWithClaims(&auth.Claims{
			UserID: "user-123",
			Email:  "test@edugo.com",
			ActiveContext: &auth.UserContext{
				RoleID:      "role-teacher",
				RoleName:    "Teacher",
				Permissions: []string{"users:read", "users:update"},
			},
		})

		middleware := RequireAllPermissions(
			enum.PermissionUsersRead,
			enum.PermissionUsersUpdate,
			enum.PermissionUsersDelete,
		)
		middleware(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusForbidden, c.Writer.Status())
	})

	t.Run("rechaza acceso si faltan todos los permisos", func(t *testing.T) {
		c := createTestContextWithClaims(&auth.Claims{
			UserID: "user-123",
			Email:  "test@edugo.com",
			ActiveContext: &auth.UserContext{
				RoleID:      "role-student",
				RoleName:    "Student",
				Permissions: []string{"materials:read"},
			},
		})

		middleware := RequireAllPermissions(
			enum.PermissionUsersCreate,
			enum.PermissionUsersDelete,
		)
		middleware(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusForbidden, c.Writer.Status())
	})

	t.Run("rechaza acceso sin claims", func(t *testing.T) {
		c := createTestContextWithClaims(nil)

		middleware := RequireAllPermissions(enum.PermissionUsersRead)
		middleware(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusUnauthorized, c.Writer.Status())
	})

	t.Run("rechaza acceso sin ActiveContext", func(t *testing.T) {
		c := createTestContextWithClaims(&auth.Claims{
			UserID: "user-123",
			Email:  "test@edugo.com",
		})

		middleware := RequireAllPermissions(enum.PermissionUsersRead)
		middleware(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusForbidden, c.Writer.Status())
	})

	t.Run("permite acceso con un solo permiso requerido", func(t *testing.T) {
		c := createTestContextWithClaims(&auth.Claims{
			UserID: "user-123",
			Email:  "test@edugo.com",
			ActiveContext: &auth.UserContext{
				RoleID:      "role-teacher",
				RoleName:    "Teacher",
				Permissions: []string{"users:read", "materials:create"},
			},
		})

		middleware := RequireAllPermissions(enum.PermissionUsersRead)
		middleware(c)

		assert.False(t, c.IsAborted())
	})

	t.Run("permite acceso cuando usuario tiene permisos extra", func(t *testing.T) {
		c := createTestContextWithClaims(&auth.Claims{
			UserID: "user-123",
			Email:  "test@edugo.com",
			ActiveContext: &auth.UserContext{
				RoleID:   "role-admin",
				RoleName: "Admin",
				Permissions: []string{
					"users:create",
					"users:read",
					"users:update",
					"users:delete",
					"schools:manage",
					"materials:create",
				},
			},
		})

		middleware := RequireAllPermissions(
			enum.PermissionUsersRead,
			enum.PermissionUsersUpdate,
		)
		middleware(c)

		assert.False(t, c.IsAborted())
	})
}

func TestGetValidatedClaims(t *testing.T) {
	t.Run("retorna claims válidos", func(t *testing.T) {
		expectedClaims := &auth.Claims{
			UserID: "user-123",
			Email:  "test@edugo.com",
			ActiveContext: &auth.UserContext{
				RoleID:      "role-teacher",
				RoleName:    "Teacher",
				Permissions: []string{"users:read"},
			},
		}

		c := createTestContextWithClaims(expectedClaims)
		claims := getValidatedClaims(c)

		assert.NotNil(t, claims)
		assert.Equal(t, expectedClaims, claims)
		assert.False(t, c.IsAborted())
	})

	t.Run("retorna nil sin claims en contexto", func(t *testing.T) {
		c := createTestContextWithClaims(nil)
		claims := getValidatedClaims(c)

		assert.Nil(t, claims)
		assert.Equal(t, http.StatusUnauthorized, c.Writer.Status())
	})

	t.Run("retorna nil con claims de tipo incorrecto", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set(ContextKeyClaims, "wrong-type")

		claims := getValidatedClaims(c)

		assert.Nil(t, claims)
		assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
	})

	t.Run("retorna nil con ActiveContext nil", func(t *testing.T) {
		c := createTestContextWithClaims(&auth.Claims{
			UserID: "user-123",
			Email:  "test@edugo.com",
		})

		claims := getValidatedClaims(c)

		assert.Nil(t, claims)
		assert.Equal(t, http.StatusForbidden, c.Writer.Status())
	})
}

func TestPermissionMiddleware_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("integración completa con token JWT real", func(t *testing.T) {
		// Crear JWTManager y generar token con contexto
		jwtManager := auth.NewJWTManager("test-secret-key", "test-issuer")

		activeContext := &auth.UserContext{
			RoleID:      "role-teacher",
			RoleName:    "Teacher",
			SchoolID:    "school-123",
			SchoolName:  "Test School",
			Permissions: []string{"materials:create", "materials:read", "materials:update"},
		}

		token, _, err := jwtManager.GenerateTokenWithContext(
			"user-123",
			"teacher@edugo.com",
			activeContext,
			time.Hour,
		)
		require.NoError(t, err)

		// Setup router con JWT middleware y permission middleware
		router := gin.New()
		router.Use(JWTAuthMiddleware(jwtManager))
		router.Use(RequirePermission(enum.PermissionMaterialsCreate))
		router.GET("/materials", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "success"})
		})

		// Request con token válido
		req := httptest.NewRequest("GET", "/materials", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("integración rechaza acceso sin permiso", func(t *testing.T) {
		jwtManager := auth.NewJWTManager("test-secret-key", "test-issuer")

		activeContext := &auth.UserContext{
			RoleID:      "role-student",
			RoleName:    "Student",
			Permissions: []string{"materials:read"}, // Solo lectura
		}

		token, _, err := jwtManager.GenerateTokenWithContext(
			"user-456",
			"student@edugo.com",
			activeContext,
			time.Hour,
		)
		require.NoError(t, err)

		router := gin.New()
		router.Use(JWTAuthMiddleware(jwtManager))
		router.Use(RequirePermission(enum.PermissionMaterialsCreate)) // Requiere create
		router.GET("/materials", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "success"})
		})

		req := httptest.NewRequest("GET", "/materials", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
