# Tests Deshabilitados (Skipped)

> Archivos de test que han sido deshabilitados y necesitan ser arreglados.

---

## Archivos con Extensión `.skip`

Estos archivos fueron renombrados de `.go` a `.skip` para excluirlos de la compilación.

### 1. factory_mongodb_integration_test.go.skip

**Ubicación:** `bootstrap/factory_mongodb_integration_test.go.skip`  
**Líneas:** 529  
**Fecha aproximada de skip:** Sprint 3

#### Contenido del Archivo
```go
package bootstrap

// TODO: Estos tests de MongoDB integration necesitan refactoring 
// para usar correctamente containers.MongoDB
// Por ahora se skipean para permitir que el coverage validation pase

import (
    "context"
    "testing"
    "time"
    
    "github.com/EduGoGroup/edugo-shared/testing/containers"
    // ...
)

func TestMongoDBFactory_CreateConnection_Success(t *testing.T) {
    t.Skip("MongoDB integration tests requieren refactoring - ver TODO en archivo")
    // ...
}
```

#### Tests Afectados
- `TestMongoDBFactory_CreateConnection_Success`
- `TestMongoDBFactory_CreateConnection_InvalidURI`
- `TestMongoDBFactory_GetDatabase`
- `TestMongoDBFactory_Ping_Success`
- `TestMongoDBFactory_Ping_Disconnected`
- `TestMongoDBFactory_Close`
- `TestMongoDBFactory_Integration_FullWorkflow`

#### Problema Raíz
Los tests usan la API antigua de `containers.MongoDB` que fue actualizada. Necesitan:
1. Actualizar llamadas a `manager.MongoDB()`
2. Usar `ConnectionString(ctx)` en lugar de acceso directo
3. Manejar correctamente el contexto

#### Solución Requerida
```go
func TestMongoDBFactory_CreateConnection_Success(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    ctx := context.Background()
    
    // Setup container
    config := containers.NewConfig().
        WithMongoDB(nil).
        Build()

    manager, err := containers.GetManager(t, config)
    require.NoError(t, err)
    
    // Obtener MongoDB container
    mongoContainer := manager.MongoDB()
    require.NotNil(t, mongoContainer)
    
    // Obtener connection string correctamente
    mongoURI, err := mongoContainer.ConnectionString(ctx)
    require.NoError(t, err)
    
    // Crear factory y probar
    factory := NewDefaultMongoDBFactory()
    client, err := factory.CreateConnection(ctx, MongoDBConfig{
        URI:      mongoURI,
        Database: "test_db",
    })
    require.NoError(t, err)
    defer factory.Close(ctx, client)
    
    // Assertions...
}
```

---

### 2. factory_postgresql_integration_test.go.skip

**Ubicación:** `bootstrap/factory_postgresql_integration_test.go.skip`  
**Líneas:** ~500 (estimado)  
**Problema:** Similar a MongoDB

#### Tests Afectados
- `TestPostgreSQLFactory_CreateConnection_Success`
- `TestPostgreSQLFactory_CreateConnection_InvalidConfig`
- `TestPostgreSQLFactory_Ping`
- `TestPostgreSQLFactory_Close`
- `TestPostgreSQLFactory_Integration_FullWorkflow`

#### Solución Requerida
Actualizar para usar `manager.PostgreSQL().ConnectionString(ctx)` y la nueva API.

---

### 3. factory_rabbitmq_integration_test.go.skip

**Ubicación:** `bootstrap/factory_rabbitmq_integration_test.go.skip`  
**Líneas:** ~450 (estimado)  
**Problema:** Similar a MongoDB

#### Tests Afectados
- `TestRabbitMQFactory_CreateConnection_Success`
- `TestRabbitMQFactory_CreateChannel`
- `TestRabbitMQFactory_DeclareQueue`
- `TestRabbitMQFactory_Close`
- `TestRabbitMQFactory_Integration_FullWorkflow`

#### Solución Requerida
Actualizar para usar `manager.RabbitMQ().ConnectionString(ctx)` y la nueva API.

---

## Impacto en Coverage

### Antes de Skip
```
bootstrap/: ~85% coverage (con tests funcionando)
```

### Después de Skip
```
bootstrap/: ~65% coverage (tests de integración no corren)
```

### Coverage Perdido
- ~20% del módulo bootstrap
- ~1500 líneas de código sin tests de integración
- Factories de MongoDB, PostgreSQL y RabbitMQ sin pruebas reales

---

## Plan de Acción

### Fase 1: Análisis (1 día)
1. Revisar cada archivo `.skip`
2. Identificar cambios necesarios en API
3. Documentar breaking changes

### Fase 2: Refactoring (2-3 días)
1. Actualizar `factory_mongodb_integration_test.go`
2. Actualizar `factory_postgresql_integration_test.go`
3. Actualizar `factory_rabbitmq_integration_test.go`
4. Ejecutar tests localmente con Docker

### Fase 3: Restauración (1 día)
1. Renombrar archivos de `.skip` a `.go`
2. Ejecutar suite completa de tests
3. Verificar coverage
4. Actualizar CI/CD si es necesario

---

## Comandos para Verificar

```bash
# Listar archivos skipped
find . -name "*.skip" -type f

# Ver contenido de un archivo skip
cat bootstrap/factory_mongodb_integration_test.go.skip

# Restaurar archivo (después de arreglar)
mv bootstrap/factory_mongodb_integration_test.go.skip \
   bootstrap/factory_mongodb_integration_test.go

# Ejecutar tests de integración
go test -v ./bootstrap/... -run Integration
```

---

## Criterios de Éxito

- [ ] Todos los archivos `.skip` renombrados a `.go`
- [ ] Todos los tests de integración pasan
- [ ] Coverage de bootstrap >= 80%
- [ ] CI/CD verde en todos los ambientes
- [ ] Sin `t.Skip()` innecesarios en los tests
