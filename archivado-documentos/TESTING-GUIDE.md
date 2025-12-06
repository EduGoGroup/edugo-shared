# Guía de Testing - edugo-shared

**Versión:** 1.0  
**Última actualización:** 2025-11-20

---

## 🎯 Filosofía de Testing

> "Los tests son documentación ejecutable del comportamiento esperado"

### Principios

1. **Tests como Documentación:** El código de test debe ser claro y legible
2. **Independencia:** Cada test debe poder ejecutarse solo
3. **Rapidez:** Tests unitarios <100ms, integración <5s
4. **Confiabilidad:** Tests determinísticos, sin flakiness
5. **Mantenibilidad:** Tests fáciles de actualizar

---

## 📋 Tipos de Tests

### 1. Tests Unitarios

**Propósito:** Validar unidades de código aisladas

**Ejemplo:**
```go
func TestHashPassword(t *testing.T) {
    password := "secret123"
    
    hash, err := HashPassword(password)
    
    assert.NoError(t, err)
    assert.NotEmpty(t, hash)
    assert.NotEqual(t, password, hash)
}
```

### 2. Tests de Integración

**Propósito:** Validar interacción con servicios externos

**Build tag:** `//go:build integration`

**Ejemplo:**
```go
//go:build integration

func TestPostgresConnection_Integration(t *testing.T) {
    container := testing.NewPostgresContainer(t)
    defer container.Close()
    
    db := Connect(container.DSN())
    
    err := db.Ping()
    assert.NoError(t, err)
}
```

### 3. Tests de Tabla (Table-Driven)

**Ejemplo:**
```go
func TestValidateEmail(t *testing.T) {
    tests := []struct{
        name    string
        email   string
        wantErr bool
    }{
        {"valid", "user@example.com", false},
        {"invalid - no @", "user.example.com", true},
        {"empty", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

---

## 🛠️ Herramientas

### Testing Frameworks

- **stdlib testing:** Base
- **testify/assert:** Aserciones claras
- **testify/suite:** Test suites
- **testify/mock:** Mocking

### Coverage

```bash
# Generar coverage
go test ./... -coverprofile=coverage.out

# Ver en terminal
go tool cover -func=coverage.out

# Ver en HTML
go tool cover -html=coverage.out
```

### Test Containers

```go
import "github.com/EduGoGroup/edugo-shared/testing/containers"

// PostgreSQL
pg := containers.NewPostgresContainer(t)
defer pg.Close()

// MongoDB
mongo := containers.NewMongoDBContainer(t)
defer mongo.Close()

// RabbitMQ
rabbit := containers.NewRabbitMQContainer(t)
defer rabbit.Close()
```

---

## 📏 Umbrales de Coverage

Ver: `.coverage-thresholds.yml`

### Por Tipo de Módulo

| Tipo | Umbral | Razón |
|------|--------|-------|
| Seguridad (auth) | >80% | Crítico |
| Core (logger, lifecycle) | >85% | Infraestructura |
| Negocio (evaluation) | >90% | Lógica crítica |
| Base de datos | >55% | Integración |
| Utilities | >70% | Amplio uso |

---

## 🚀 Comandos Útiles

```bash
# Tests de un módulo
cd auth && go test ./...

# Tests con coverage
go test ./... -cover

# Tests con race detector
go test ./... -race

# Tests solo short
go test ./... -short

# Tests solo integration
go test ./... -tags=integration

# Tests específicos
go test -run TestJWT

# Tests verbose
go test -v ./...
```

---

## ✅ Checklist de Test

### Antes de Commit

- [ ] Tests pasan localmente
- [ ] Coverage no disminuye
- [ ] Tests son independientes
- [ ] No hay prints/debugs
- [ ] Nombres descriptivos

### En PR

- [ ] Tests cubren cambios nuevos
- [ ] Tests de edge cases
- [ ] Tests de error handling
- [ ] Documentación actualizada
- [ ] CI/CD pasa

---

## 📖 Recursos

- [Go Testing Guide](https://golang.org/doc/code#Testing)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Coverage Thresholds](./.coverage-thresholds.yml)
- [Coverage Strategy](./docs/cicd/coverage-analysis/STRATEGY.md)

---

**Mantenido por:** EduGo Team  
**Revisión:** Cada sprint
