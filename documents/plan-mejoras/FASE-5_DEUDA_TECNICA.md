# FASE 5: Deuda Técnica

> **Prioridad**: BAJA  
> **Duración estimada**: Ongoing (Boy Scout Rule)  
> **Prerrequisitos**: Fases 1-4 completadas (opcional)  
> **Rama**: `fase-5-deuda-tecnica-[item]`  
> **Objetivo**: Reducir deuda técnica acumulada de forma continua

---

## Flujo de Trabajo de Esta Fase

### Filosofía

Esta fase sigue el principio **"Boy Scout Rule"**:
> "Deja el código mejor de lo que lo encontraste"

No es necesario completar todos los pasos en secuencia. En su lugar:
- Elige un paso cuando tengas tiempo libre
- Incluye una mejora pequeña en cada PR de otras fases
- Programa sprints de deuda técnica trimestrales

### Para Cada Item de Deuda

```bash
# 1. Crear rama específica
git checkout dev
git pull origin dev
git checkout -b fase-5-deuda-tecnica-[nombre-item]

# 2. Implementar mejora
# ...

# 3. Commit
git add .
git commit -m "docs/test/chore(módulo): descripción"

# 4. Push y PR
git push origin fase-5-deuda-tecnica-[nombre-item]

# 5. Esperar revisión de GitHub Copilot
# - DESCARTAR: Comentarios de traducción inglés/español
# - CORREGIR: Problemas importantes detectados

# 6. Esperar pipelines (máx 10 min, revisar cada 1 min)
# - Si hay errores: Corregir (regla de 3 intentos)

# 7. Merge a dev
```

---

## Paso 5.1: Agregar Documentación GoDoc

### Problema

Muchas funciones públicas carecen de documentación GoDoc completa.

### Archivos Prioritarios

| Archivo | Funciones a Documentar |
|---------|------------------------|
| `auth/jwt.go` | NewJWTManager, GenerateToken, ValidateToken |
| `database/postgres/connection.go` | Connect, HealthCheck, Close |
| `database/mongodb/connection.go` | Connect, GetDatabase, Close |
| `messaging/rabbit/publisher.go` | NewPublisher, Publish |
| `messaging/rabbit/consumer.go` | NewConsumer, Consume |

### Template de Documentación

```go
// NewJWTManager creates a new JWT token manager.
//
// The manager handles generation, validation, and refresh of JWT tokens
// using HS256 signing algorithm.
//
// Parameters:
//   - secretKey: The secret key used for signing tokens.
//     Minimum 32 characters recommended for security.
//   - issuer: The issuer claim to be included in generated tokens.
//
// Returns a configured JWTManager ready to generate and validate tokens.
//
// Example:
//
//	manager := auth.NewJWTManager("my-32-char-secret-key!!", "my-app")
//	token, err := manager.GenerateToken(userID, email, role, 24*time.Hour)
func NewJWTManager(secretKey, issuer string) *JWTManager {
    // ...
}
```

### Comando para Verificar Cobertura GoDoc

```bash
# Instalar godoc si no está
go install golang.org/x/tools/cmd/godoc@latest

# Iniciar servidor local
godoc -http=:6060

# Navegar a http://localhost:6060/pkg/github.com/EduGoGroup/edugo-shared/
```

### Criterios de Éxito por Función

- [ ] Descripción clara de propósito
- [ ] Parámetros documentados
- [ ] Valores de retorno documentados
- [ ] Ejemplo de uso (opcional pero recomendado)

### Commit

```bash
git add auth/jwt.go
git commit -m "docs(auth): agregar GoDoc a funciones de JWT

- Documenta NewJWTManager, GenerateToken, ValidateToken
- Incluye descripción de parámetros y retornos
- Agrega ejemplos de uso"
```

---

## Paso 5.2: Crear Funciones Example

### Problema

No hay funciones `Example*` que sirvan como documentación ejecutable.

### Template

```go
// auth/jwt_example_test.go
package auth_test

import (
    "fmt"
    "time"
    
    "github.com/EduGoGroup/edugo-shared/auth"
    "github.com/EduGoGroup/edugo-shared/common/types/enum"
)

func ExampleNewJWTManager() {
    manager := auth.NewJWTManager(
        "my-super-secret-key-32-chars!!", 
        "my-app",
    )
    
    token, err := manager.GenerateToken(
        "user-123",
        "user@example.com",
        enum.SystemRoleStudent,
        24*time.Hour,
    )
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    fmt.Println("Token generated successfully")
    fmt.Println("Length:", len(token) > 0)
    // Output:
    // Token generated successfully
    // Length: true
}

func ExampleJWTManager_ValidateToken() {
    manager := auth.NewJWTManager(
        "my-super-secret-key-32-chars!!", 
        "my-app",
    )
    
    token, _ := manager.GenerateToken(
        "user-123",
        "user@example.com",
        enum.SystemRoleStudent,
        time.Hour,
    )
    
    claims, err := manager.ValidateToken(token)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    fmt.Println("UserID:", claims.UserID)
    fmt.Println("Email:", claims.Email)
    // Output:
    // UserID: user-123
    // Email: user@example.com
}
```

### Ejemplos Necesarios

- [ ] `ExampleNewJWTManager`
- [ ] `ExampleJWTManager_GenerateToken`
- [ ] `ExampleJWTManager_ValidateToken`
- [ ] `ExampleHashPassword`
- [ ] `ExampleNewPublisher`
- [ ] `ExampleNewConsumer`

### Verificar que Ejemplos Pasan

```bash
go test -v -run Example ./auth/...
```

### Commit

```bash
git add auth/jwt_example_test.go
git commit -m "docs(auth): agregar funciones Example para JWT

- ExampleNewJWTManager muestra creación de manager
- ExampleJWTManager_ValidateToken muestra validación
- Ejemplos verificables en CI"
```

---

## Paso 5.3: Agregar Benchmarks

### Problema

No hay benchmarks para operaciones críticas de performance.

### Benchmarks Prioritarios

```go
// auth/jwt_benchmark_test.go
package auth

import (
    "testing"
    "time"
    
    "github.com/EduGoGroup/edugo-shared/common/types/enum"
)

func BenchmarkJWTManager_GenerateToken(b *testing.B) {
    manager := NewJWTManager("test-secret-key-32-characters!!", "test")
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = manager.GenerateToken(
            "user-123",
            "test@test.com",
            enum.SystemRoleStudent,
            time.Hour,
        )
    }
}

func BenchmarkJWTManager_ValidateToken(b *testing.B) {
    manager := NewJWTManager("test-secret-key-32-characters!!", "test")
    token, _ := manager.GenerateToken(
        "user-123",
        "test@test.com",
        enum.SystemRoleStudent,
        time.Hour,
    )
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = manager.ValidateToken(token)
    }
}

func BenchmarkHashPassword(b *testing.B) {
    password := "test-password-123"
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = HashPassword(password)
    }
}

func BenchmarkVerifyPassword(b *testing.B) {
    password := "test-password-123"
    hash, _ := HashPassword(password)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = VerifyPassword(hash, password)
    }
}
```

### Ejecutar Benchmarks

```bash
# Benchmark específico
go test -bench=BenchmarkJWTManager -benchmem ./auth/...

# Todos los benchmarks
go test -bench=. -benchmem ./...

# Con comparación
go test -bench=. -benchmem -count=10 ./... | tee bench_new.txt
benchstat bench_old.txt bench_new.txt
```

### Criterios de Éxito

- [ ] Al menos 5 benchmarks creados
- [ ] Operaciones críticas cubiertas (JWT, password, DB)
- [ ] Resultados documentados como baseline

### Commit

```bash
git add auth/jwt_benchmark_test.go
git commit -m "perf(auth): agregar benchmarks para operaciones JWT

- BenchmarkJWTManager_GenerateToken
- BenchmarkJWTManager_ValidateToken
- BenchmarkHashPassword
- BenchmarkVerifyPassword"
```

---

## Paso 5.4: Configurar .golangci.yml

### Problema

No hay configuración de linter personalizada.

### Crear Archivo de Configuración

```yaml
# .golangci.yml
run:
  timeout: 5m
  modules-download-mode: readonly
  tests: true

linters:
  enable:
    # Defectos
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - unused
    
    # Formato
    - gofmt
    - goimports
    
    # Estilo
    - misspell
    - unconvert
    - unparam
    - nakedret
    
    # Performance
    - prealloc
    
    # Seguridad
    - gosec
    
    # Complejidad
    - gocyclo
    - gocritic
    
    # Documentación
    - revive

linters-settings:
  errcheck:
    check-type-assertions: true
    check-blank: true
    
  govet:
    check-shadowing: true
    
  gocyclo:
    min-complexity: 15
    
  revive:
    rules:
      - name: exported
        severity: warning
        arguments:
          - checkPublicInterface
          - disableStutteringCheck
      - name: blank-imports
        severity: warning
      - name: context-as-argument
        severity: warning
      - name: error-return
        severity: warning
      - name: error-strings
        severity: warning
      - name: error-naming
        severity: warning
        
  gosec:
    excludes:
      - G104  # Audit errors not checked
      
  misspell:
    locale: US

issues:
  exclude-rules:
    # Tests pueden tener más libertad
    - path: _test\.go
      linters:
        - errcheck
        - gosec
        - gocritic
        
    # Ejemplos son más relajados
    - path: _example_test\.go
      linters:
        - errcheck
        
    # Código generado
    - path: \.pb\.go
      linters:
        - all
        
  max-issues-per-linter: 50
  max-same-issues: 10
```

### Ejecutar con Nueva Config

```bash
golangci-lint run --config .golangci.yml ./...
```

### Criterios de Éxito

- [ ] Archivo `.golangci.yml` creado
- [ ] `golangci-lint run` pasa sin errores
- [ ] CI actualizado para usar config

### Commit

```bash
git add .golangci.yml
git commit -m "chore: agregar configuración de golangci-lint

- Habilita linters recomendados
- Configura reglas específicas del proyecto
- Excluye tests de reglas estrictas"
```

---

## Paso 5.5: Convertir Tests a Table-Driven

### Problema

Algunos tests no usan el patrón table-driven.

### Ejemplo de Conversión

**Antes:**
```go
func TestValidateToken(t *testing.T) {
    manager := NewJWTManager("secret", "issuer")
    
    _, err := manager.ValidateToken("")
    assert.Error(t, err)
    
    _, err = manager.ValidateToken("invalid")
    assert.Error(t, err)
    
    claims, err := manager.ValidateToken(validToken)
    assert.NoError(t, err)
    assert.NotNil(t, claims)
}
```

**Después:**
```go
func TestValidateToken(t *testing.T) {
    manager := NewJWTManager("test-secret-32-characters-long!!", "test")
    validToken, _ := manager.GenerateToken("user", "email", enum.SystemRoleStudent, time.Hour)
    
    tests := []struct {
        name    string
        token   string
        wantErr bool
        errMsg  string
    }{
        {
            name:    "empty token returns error",
            token:   "",
            wantErr: true,
            errMsg:  "invalid",
        },
        {
            name:    "malformed token returns error",
            token:   "not.a.valid.jwt",
            wantErr: true,
            errMsg:  "invalid",
        },
        {
            name:    "valid token returns claims",
            token:   validToken,
            wantErr: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            claims, err := manager.ValidateToken(tt.token)
            
            if tt.wantErr {
                assert.Error(t, err)
                if tt.errMsg != "" {
                    assert.Contains(t, err.Error(), tt.errMsg)
                }
                assert.Nil(t, claims)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, claims)
            }
        })
    }
}
```

### Tests a Convertir

- [ ] `auth/jwt_test.go`
- [ ] `common/errors/errors_test.go`
- [ ] `config/base_test.go`

### Criterios de Éxito

- [ ] Tests usan `t.Run` para sub-tests
- [ ] Casos de test en slice de structs
- [ ] Fácil agregar nuevos casos

### Commit

```bash
git add auth/jwt_test.go
git commit -m "test(auth): convertir tests JWT a table-driven

- Usa t.Run para sub-tests
- Facilita agregar nuevos casos
- Mejora legibilidad y mantenimiento"
```

---

## Plan de Reducción Continua

### Sprint de Deuda (Trimestral)

Dedicar 1-2 días por trimestre a:
1. Revisar esta documentación
2. Priorizar items pendientes
3. Resolver deuda acumulada
4. Actualizar métricas

### Boy Scout en PRs

Cada PR debería incluir al menos UNA mejora pequeña:
- Agregar GoDoc a una función
- Convertir un test a table-driven
- Agregar un Example
- Corregir un typo

### Métricas a Trackear

| Métrica | Cómo Medir |
|---------|------------|
| Cobertura GoDoc | Revisar godoc localmente |
| Ejemplos | `go test -run Example` |
| Benchmarks | `go test -bench .` |
| Linter warnings | `golangci-lint run` |

---

## Verificación para Cada PR de Deuda

### Antes de Crear el PR

```bash
# GoDoc
go doc ./...

# Examples
go test -v -run Example ./...

# Benchmarks (si aplica)
go test -bench=. -benchmem ./...

# Linter
golangci-lint run ./...

# Tests
make test-all-modules
```

### Revisión de GitHub Copilot

| Tipo de Comentario | Acción |
|-------------------|--------|
| Traducción inglés/español | DESCARTAR |
| Mejora de documentación | EVALUAR |
| Error detectado | CORREGIR |

### Esperar Pipelines

```bash
# Revisar cada minuto durante máximo 10 minutos
# Si hay errores:
#   1. Analizar causa
#   2. Corregir (máx 3 intentos)
#   3. Push y esperar nuevamente
```

---

## Resumen de la Fase 5

| Paso | Descripción | Prioridad |
|------|-------------|-----------|
| 5.1 | Documentación GoDoc | Boy Scout |
| 5.2 | Funciones Example | Boy Scout |
| 5.3 | Benchmarks | Sprint |
| 5.4 | Config golangci-lint | Una vez |
| 5.5 | Tests table-driven | Boy Scout |

---

## Conclusión

Esta fase no tiene un "fin" definido. Es un proceso continuo de mejora que mantiene la calidad del código a largo plazo.

**Recordatorio**: No intentes hacer todo de una vez. Pequeñas mejoras consistentes tienen mayor impacto que grandes refactorizaciones esporádicas.
