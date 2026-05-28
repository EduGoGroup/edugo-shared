# Testing Guide - database/postgres

## Requisitos

- Docker instalado y corriendo
- Go 1.24+
- Acceso a internet (para descargar imágenes de Docker)

## Ejecutar Tests de Integración

Los tests de integración usan **Testcontainers** para levantar automáticamente un container de PostgreSQL.

### Comando básico

```bash
cd /home/user/edugo-shared/database/postgres
go test -v -cover ./...
```

### Ejecutar solo tests de integración

```bash
go test -v -run Integration -cover ./...
```

### Ejecutar con logs de Docker (debug)

```bash
export TESTCONTAINERS_RYUK_DISABLED=false
go test -v -cover ./...
```

## Estructura de Tests

### connection_test.go

Tests de conexión, health checks y gestión del pool:

- `TestHealthCheck_Integration`: Verifica health checks activos y cerrados
- `TestGetStats_Integration`: Valida estadísticas del connection pool
- `TestClose_Integration`: Prueba cierre de conexiones

**Cobertura esperada:** >70% de connection.go

### transaction_test.go

Tests de transacciones y manejo de errores:

- `TestWithTransaction_Integration`: Commits y rollbacks exitosos
- `TestWithTransactionIsolation_Integration`: Niveles de aislamiento (ReadCommitted, Serializable)
- Tests de panic handling con rollback automático

**Cobertura esperada:** >80% de transaction.go

## Cobertura Total Esperada

**>80%** para el módulo completo

## Verificación de Cobertura

```bash
# Generar reporte de coverage
go test -coverprofile=coverage.out ./...

# Ver coverage por función
go tool cover -func=coverage.out

# Generar HTML interactivo
go tool cover -html=coverage.out -o coverage.html
# Abrir coverage.html en navegador
```

## Troubleshooting

### Error: "Cannot connect to Docker daemon"

```bash
# Verificar que Docker está corriendo
docker ps

# En Linux, verificar permisos
sudo usermod -aG docker $USER
# Logout/login para aplicar cambios
```

### Error: "Testcontainers timeout"

```bash
# Aumentar timeout
export TESTCONTAINERS_TIMEOUT=300

# Limpiar containers viejos
docker system prune -f
```

### Error: "Port already in use"

```bash
# Testcontainers usa puertos aleatorios, pero si hay conflicto:
docker ps -a | grep postgres
docker rm -f <container_id>
```

## Skip Tests de Integración

Si necesitas ejecutar solo tests unitarios (sin Docker):

```bash
go test -v -short ./...
```

Los tests de integración tienen:
```go
if testing.Short() {
    t.Skip("Skipping integration test en modo short")
}
```

## Notas Importantes

1. **Primera ejecución lenta**: Testcontainers descarga la imagen `postgres:15-alpine` (~80MB). Ejecuciones siguientes son rápidas.

2. **Cleanup automático**: Testcontainers limpia los containers automáticamente al terminar los tests.

3. **Tests en CI/CD**: Los tests funcionan en cualquier CI con Docker (GitHub Actions, GitLab CI, etc.)

## Validación de Implementación Completa

Para validar que todos los tests pasan:

```bash
# Desde la raíz del módulo
cd /home/user/edugo-shared/database/postgres

# Ejecutar tests
go test -v -cover ./...

# Verificar salida esperada:
# ✅ TestHealthCheck_Integration PASS
# ✅ TestGetStats_Integration PASS
# ✅ TestClose_Integration PASS
# ✅ TestWithTransaction_Integration PASS
# ✅ TestWithTransactionIsolation_Integration PASS
#
# coverage: >80% of statements
```

## Para Claude Code Local

Si estás ejecutando estos tests desde **Claude Code desktop** (que tiene acceso a Docker):

1. Verifica que Docker está corriendo: `docker ps`
2. Ejecuta los tests normalmente: `go test -v -cover ./...`
3. Los tests deberían pasar sin problemas

Los tests fueron diseñados para ejecutarse en cualquier ambiente con Docker disponible.
