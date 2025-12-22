# Código Incompleto (TODOs)

> Funciones y código que tienen implementación pendiente marcada con TODO.

---

## Issue #1: GetPresignedURL No Implementado

### Ubicación
```
bootstrap/resource_implementations.go:127-131
```

### Código Actual
```go
// GetPresignedURL genera una URL pre-firmada para acceso temporal
func (c *defaultStorageClient) GetPresignedURL(ctx context.Context, key string, expirationMinutes int) (string, error) {
    // TODO: Implementar con presign client
    return "", fmt.Errorf("presigned URL not implemented yet")
}
```

### Problema
- La función está declarada en la interface `StorageClient`
- La implementación retorna error siempre
- El campo `presignClient` existe pero no se usa
- Cualquier servicio que intente usar URLs pre-firmadas fallará

### Impacto
- **Alto**: Funcionalidad crítica para compartir archivos temporalmente
- Servicios como materials-service pueden necesitar esta funcionalidad

### Solución Sugerida
```go
import (
    "time"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c *defaultStorageClient) GetPresignedURL(ctx context.Context, key string, expirationMinutes int) (string, error) {
    presignClient, ok := c.presignClient.(*s3.PresignClient)
    if !ok {
        return "", fmt.Errorf("presign client not properly initialized")
    }
    
    request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
        Bucket: aws.String(c.bucket),
        Key:    aws.String(key),
    }, func(opts *s3.PresignOptions) {
        opts.Expires = time.Duration(expirationMinutes) * time.Minute
    })
    if err != nil {
        return "", fmt.Errorf("failed to generate presigned URL: %w", err)
    }
    
    return request.URL, nil
}
```

### Cambios Adicionales Requeridos
1. Cambiar tipo de `presignClient` de `interface{}` a `*s3.PresignClient`
2. Agregar tests unitarios
3. Agregar tests de integración con MinIO

### Prioridad: **ALTA**

---

## Issue #2: extractEnvAndVersion Hardcodeado

### Ubicación
```
bootstrap/bootstrap.go:431-436
```

### Código Actual
```go
// extractEnvAndVersion extrae environment y version de la configuración
// Por ahora retorna valores por defecto, será implementado según BaseConfig
func extractEnvAndVersion(config interface{}) (string, string) {
    // TODO: Implementar extracción real cuando BaseConfig esté integrado
    return "local", "0.0.0"
}
```

### Problema
- La función siempre retorna valores hardcodeados
- El logger no refleja el ambiente real de ejecución
- Dificulta debugging en ambientes no-local
- El comentario indica que BaseConfig ya debería estar integrado

### Impacto
- **Medio**: Los logs no muestran el ambiente correcto
- Dificulta troubleshooting en producción

### Solución Sugerida
```go
func extractEnvAndVersion(config interface{}) (string, string) {
    v := reflect.ValueOf(config)
    if v.Kind() == reflect.Ptr {
        v = v.Elem()
    }
    
    if v.Kind() != reflect.Struct {
        return "unknown", "0.0.0"
    }
    
    // Buscar campo Environment
    envField := v.FieldByName("Environment")
    env := "unknown"
    if envField.IsValid() && envField.Kind() == reflect.String {
        env = envField.String()
    }
    
    // Buscar campo Version (opcional)
    versionField := v.FieldByName("Version")
    version := "0.0.0"
    if versionField.IsValid() && versionField.Kind() == reflect.String {
        version = versionField.String()
    }
    
    return env, version
}
```

### Prioridad: **MEDIA**

---

## Issue #3: Tests de Integración MongoDB Requieren Refactoring

### Ubicación
```
bootstrap/factory_mongodb_integration_test.go.skip:1-4
```

### Código Actual
```go
package bootstrap

// TODO: Estos tests de MongoDB integration necesitan refactoring para usar correctamente containers.MongoDB
// Por ahora se skipean para permitir que el coverage validation pase
```

### Problema
- 529 líneas de tests están deshabilitados
- El archivo tiene extensión `.skip` en lugar de `.go`
- Los tests probablemente fallan por cambios en la API de containers
- Coverage del módulo bootstrap está artificialmente inflado

### Impacto
- **Alto**: No hay tests de integración para MongoDB factory
- Bugs pueden llegar a producción sin detectarse

### Ver también
- `TESTS_SKIPPED.md` para análisis detallado

### Prioridad: **ALTA**

---

## Resumen de TODOs

| Archivo | Línea | Descripción | Prioridad |
|---------|-------|-------------|-----------|
| `resource_implementations.go` | 129 | GetPresignedURL | Alta |
| `bootstrap.go` | 434 | extractEnvAndVersion | Media |
| `factory_mongodb_integration_test.go.skip` | 3 | Tests MongoDB | Alta |

---

## Checklist de Resolución

- [ ] Implementar `GetPresignedURL` con `s3.PresignClient`
- [ ] Implementar `extractEnvAndVersion` con reflection
- [ ] Arreglar tests de integración de MongoDB
- [ ] Arreglar tests de integración de PostgreSQL
- [ ] Arreglar tests de integración de RabbitMQ
- [ ] Remover extensión `.skip` de archivos de test
- [ ] Verificar coverage después de arreglos
