# Deuda Técnica General

> Aspectos del proyecto que no son bugs ni malas prácticas críticas, pero que acumulan deuda técnica y deberían abordarse eventualmente.

---

## 1. Documentación de Código

### Problema
Muchas funciones públicas carecen de documentación GoDoc completa.

### Ejemplos
```go
// Falta documentación de parámetros y retornos
func NewJWTManager(secretKey, issuer string) *JWTManager

// Debería ser:
// NewJWTManager creates a new JWT token manager.
//
// Parameters:
//   - secretKey: The secret key used for signing tokens. Minimum 32 characters recommended.
//   - issuer: The issuer claim to be included in generated tokens.
//
// Returns a configured JWTManager ready to generate and validate tokens.
func NewJWTManager(secretKey, issuer string) *JWTManager
```

### Archivos Afectados
- `auth/jwt.go` - Funciones de JWT
- `database/postgres/connection.go` - Funciones de conexión
- `database/mongodb/connection.go` - Funciones de conexión
- `messaging/rabbit/*.go` - Publisher/Consumer

### Impacto
- **Bajo**: El código funciona, pero dificulta onboarding de nuevos desarrolladores

### Acción Recomendada
Agregar GoDoc comments siguiendo las convenciones de Go.

### Prioridad: **BAJA**

---

## 2. Falta de Ejemplos en Código

### Problema
No hay funciones `Example*` en los archivos de test que sirvan como documentación ejecutable.

### Estado Actual
```bash
$ grep -r "func Example" .
# Sin resultados
```

### Solución
```go
// auth/jwt_test.go
func ExampleJWTManager_GenerateToken() {
    manager := auth.NewJWTManager("my-secret-key-min-32-chars!!", "my-app")
    
    token, err := manager.GenerateToken(
        "user-123",
        "user@example.com",
        enum.SystemRoleStudent,
        time.Hour*24,
    )
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Token generated successfully")
    fmt.Println("Length:", len(token))
    // Output:
    // Token generated successfully
}
```

### Beneficios
- Documentación que se valida en CI
- Aparece en godoc automáticamente
- Ejemplos siempre actualizados

### Prioridad: **BAJA**

---

## 3. Versionado de Módulos Internos

### Problema
Los módulos internos no tienen versionado semántico claro entre ellos.

### Estado Actual
```
auth/go.mod           → depende de common sin versión específica
bootstrap/go.mod      → depende de varios módulos sin versiones
messaging/rabbit/go.mod → igual
```

### Riesgo
- Cambios en `common` pueden romper otros módulos
- No hay forma de hacer rollback de dependencias internas
- Difícil saber qué versión de cada módulo es compatible

### Solución a Largo Plazo
1. Cada módulo con su propio tag de versión
2. Dependencias internas con versiones específicas
3. CI que valide compatibilidad entre versiones

### Prioridad: **BAJA** (el monorepo actual funciona)

---

## 4. Falta de Benchmarks

### Problema
No hay benchmarks para operaciones críticas de performance.

### Operaciones que Deberían Tener Benchmarks
- JWT Generation/Validation
- Password Hashing
- Database Connection Pool
- Message Serialization/Deserialization
- S3 Upload/Download

### Ejemplo de Benchmark Necesario
```go
// auth/jwt_benchmark_test.go
func BenchmarkJWTManager_GenerateToken(b *testing.B) {
    manager := NewJWTManager("test-secret-key-32-characters!!", "test")
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = manager.GenerateToken("user-123", "test@test.com", enum.SystemRoleStudent, time.Hour)
    }
}

func BenchmarkJWTManager_ValidateToken(b *testing.B) {
    manager := NewJWTManager("test-secret-key-32-characters!!", "test")
    token, _ := manager.GenerateToken("user-123", "test@test.com", enum.SystemRoleStudent, time.Hour)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = manager.ValidateToken(token)
    }
}
```

### Prioridad: **BAJA**

---

## 5. Configuración de Linter Incompleta

### Problema
No hay archivo `.golangci.yml` con reglas personalizadas.

### Estado Actual
El Makefile ejecuta `golangci-lint run` sin configuración específica.

### Configuración Recomendada
```yaml
# .golangci.yml
run:
  timeout: 5m
  modules-download-mode: readonly

linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - unused
    - gofmt
    - goimports
    - misspell
    - unconvert
    - unparam
    - nakedret
    - prealloc
    - exportloopref
    - gocritic
    - revive

linters-settings:
  errcheck:
    check-type-assertions: true
  govet:
    check-shadowing: true
  revive:
    rules:
      - name: exported
        severity: warning

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - errcheck
```

### Prioridad: **BAJA**

---

## 6. Logs sin Contexto Estructurado Consistente

### Problema
Los logs no tienen un formato consistente de campos.

### Ejemplos de Inconsistencia
```go
// En un lugar
logger.Info("User created", "user_id", userID)

// En otro lugar
logger.Info("User created", "userId", userID)

// En otro lugar
logger.WithField("user", userID).Info("User created")
```

### Solución
Definir constantes para campos comunes:
```go
// logger/fields.go
const (
    FieldUserID      = "user_id"
    FieldRequestID   = "request_id"
    FieldService     = "service"
    FieldOperation   = "operation"
    FieldDuration    = "duration_ms"
    FieldError       = "error"
    FieldStatusCode  = "status_code"
)
```

### Prioridad: **BAJA**

---

## 7. Tests sin Table-Driven Pattern

### Problema
Algunos tests no usan el patrón table-driven, haciendo difícil agregar casos.

### Ejemplo de Test a Mejorar
```go
// Actual - repetitivo
func TestValidateToken(t *testing.T) {
    // Test 1
    _, err := manager.ValidateToken("")
    assert.Error(t, err)
    
    // Test 2
    _, err = manager.ValidateToken("invalid")
    assert.Error(t, err)
    
    // Test 3
    _, err = manager.ValidateToken(expiredToken)
    assert.Error(t, err)
}

// Mejor - table-driven
func TestValidateToken(t *testing.T) {
    tests := []struct {
        name    string
        token   string
        wantErr bool
        errType error
    }{
        {"empty token", "", true, ErrInvalidToken},
        {"malformed token", "invalid", true, ErrInvalidToken},
        {"expired token", expiredToken, true, ErrTokenExpired},
        {"valid token", validToken, false, nil},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := manager.ValidateToken(tt.token)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Prioridad: **BAJA**

---

## 8. Archivos de Test Grandes

### Problema
Algunos archivos de test son muy grandes (>500 líneas).

### Archivos Afectados
- `messaging/rabbit/consumer_test.go` (17182 bytes)
- `messaging/rabbit/publisher_test.go` (16867 bytes)
- `auth/jwt_test.go` (15426 bytes)

### Solución
Dividir por funcionalidad:
```
consumer_test.go → 
    consumer_basic_test.go
    consumer_ack_test.go
    consumer_error_test.go
    consumer_integration_test.go
```

### Prioridad: **BAJA**

---

## Matriz de Deuda Técnica

| Item | Impacto | Esfuerzo | ROI | Prioridad |
|------|---------|----------|-----|-----------|
| Documentación GoDoc | Bajo | Medio | Medio | Baja |
| Ejemplos ejecutables | Bajo | Bajo | Alto | Baja |
| Versionado módulos | Medio | Alto | Bajo | Baja |
| Benchmarks | Bajo | Medio | Medio | Baja |
| Configuración linter | Bajo | Bajo | Alto | Baja |
| Logs estructurados | Bajo | Bajo | Medio | Baja |
| Tests table-driven | Bajo | Medio | Medio | Baja |
| Archivos test grandes | Bajo | Medio | Bajo | Baja |

---

## Plan de Reducción de Deuda

### Enfoque "Boy Scout Rule"
> "Deja el código mejor de lo que lo encontraste"

Cada PR debería incluir al menos UNA mejora pequeña:
- Agregar GoDoc a una función
- Convertir un test a table-driven
- Agregar un Example
- Corregir un typo

### Sprint de Deuda (Trimestral)
Dedicar 1-2 días por trimestre a:
- Revisar esta documentación
- Priorizar items
- Resolver deuda acumulada
- Actualizar documentación de mejoras
