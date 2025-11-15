# Testing Guide - Sprint 3

## Objetivo

Aumentar coverage de **config/** y **bootstrap/** a >80%.

## Requisitos

- **Go 1.24.10** (requerido)
- Docker instalado y corriendo (para tests de integración)
- Acceso a internet

## Tests Creados

### config/
- `config/loader_test.go` - Tests comprehensivos para el loader de configuración
  - NewLoader con opciones
  - Load con archivos YAML
  - Load con variables de entorno
  - LoadFromFile
  - Get*, GetString, GetInt, GetBool
  - Múltiples opciones y overrides

**Coverage esperado:** >80%

### bootstrap/
- `bootstrap/options_test.go` - Tests para opciones de bootstrap
  - DefaultBootstrapOptions
  - WithRequiredResources
  - WithOptionalResources
  - WithSkipHealthCheck
  - WithMockFactories
  - WithStopOnFirstError
  - ApplyOptions

- `bootstrap/resources_test.go` - Tests para contenedor de recursos
  - HasLogger
  - HasPostgreSQL
  - HasMongoDB
  - HasMessagePublisher
  - HasStorageClient
  - Configuraciones parciales y completas

**Coverage esperado:** >80%

## Ejecutar Tests

### config/

```bash
cd /home/user/edugo-shared/config
go test -v -cover ./...
```

**Salida esperada:**
```
=== RUN   TestNewLoader_DefaultValues
--- PASS: TestNewLoader_DefaultValues
=== RUN   TestNewLoader_WithOptions
--- PASS: TestNewLoader_WithOptions
...
PASS
coverage: >80% of statements
ok      github.com/EduGoGroup/edugo-shared/config
```

### bootstrap/

```bash
cd /home/user/edugo-shared/bootstrap
go test -v -cover ./...
```

**Salida esperada:**
```
=== RUN   TestDefaultBootstrapOptions
--- PASS: TestDefaultBootstrapOptions
=== RUN   TestResources_HasLogger
--- PASS: TestResources_HasLogger
...
PASS
coverage: >80% of statements
ok      github.com/EduGoGroup/edugo-shared/bootstrap
```

## Verificación de Coverage Total

Para verificar coverage exacto:

```bash
# config/
cd /home/user/edugo-shared/config
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# bootstrap/
cd /home/user/edugo-shared/bootstrap
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

## Troubleshooting

### Error: "go: downloading go1.24.10"

Si estás en un ambiente sin acceso a internet:

1. **Solución**: Usa Claude Code **desktop/local** que tiene acceso a internet
2. Los tests fueron diseñados para ejecutarse en ambiente local con go 1.24.10

### Tests fallan por dependencias

```bash
# Dentro de cada módulo:
go mod tidy
go test -v -cover ./...
```

## Para Claude Code Desktop

Estos tests fueron creados en **Claude Code web** (que no tiene acceso a go 1.24.10 ni internet).

**Para ejecutar en Claude Code desktop:**

1. Verifica que tengas go 1.24.10:
   ```bash
   go version
   # Debe mostrar: go version go1.24.10 ...
   ```

2. Ejecuta los tests:
   ```bash
   cd /home/user/edugo-shared

   # Test config
   cd config && go test -v -cover ./...

   # Test bootstrap
   cd ../bootstrap && go test -v -cover ./...
   ```

3. Verifica coverage >80% en ambos módulos

## Próximos Pasos (Sprint 3 Día 2)

Una vez que los tests pasen:

1. Ejecutar suite completa: `make test-all-modules`
2. Calcular coverage global: `make coverage-all-modules`
3. Validar objetivo: >85% coverage global
4. Commit y push de los tests

## Notas

- Los tests de `loader_test.go` crean archivos temporales YAML
- Los tests de `bootstrap_test.go` usan mocks para evitar dependencias reales
- Todos los tests son unitarios y no requieren Docker (excepto database/postgres)
