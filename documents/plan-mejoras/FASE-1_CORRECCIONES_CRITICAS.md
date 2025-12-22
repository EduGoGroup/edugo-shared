# FASE 1: Correcciones Críticas

> **Prioridad**: ALTA  
> **Duración estimada**: 2-3 días  
> **Prerrequisitos**: Ninguno  
> **Rama**: `fase-1-correcciones-criticas`  
> **Objetivo**: Corregir errores críticos que afectan funcionalidad en producción

---

## Flujo de Trabajo de Esta Fase

### 1. Inicio de la Fase

```bash
# Asegurarse de estar en dev actualizado
git checkout dev
git pull origin dev

# Crear rama de la fase
git checkout -b fase-1-correcciones-criticas

# Verificar estado inicial
make build
make test-all-modules
```

### 2. Durante la Fase

- Ejecutar cada paso en orden
- Commit atómico después de cada paso completado
- Verificar que tests pasen después de cada cambio

### 3. Fin de la Fase

```bash
# Push de la rama
git push origin fase-1-correcciones-criticas

# Crear PR en GitHub hacia dev
# - Título: "feat: Fase 1 - Correcciones Críticas"
# - Descripción: Lista de cambios realizados

# Esperar revisión de GitHub Copilot
# - DESCARTAR: Comentarios de traducción inglés/español
# - CORREGIR: Problemas importantes detectados
# - DOCUMENTAR: Lo que queda como deuda técnica futura

# Esperar pipelines (máx 10 min, revisar cada 1 min)
# - Si hay errores: Corregir (regla de 3 intentos)
# - Todos los errores se corrigen (propios o heredados)

# Merge cuando todo esté verde
```

---

## Resumen de la Fase

Esta fase aborda los problemas más urgentes del código:
- Funciones no implementadas que retornan errores
- Errores silenciados que ocultan problemas reales
- Código hardcodeado que no refleja el ambiente real

---

## Paso 1.1: Implementar GetPresignedURL

### Información General

| Campo | Valor |
|-------|-------|
| **Archivo** | `bootstrap/resource_implementations.go` |
| **Líneas** | 127-131 |
| **Tipo** | Código incompleto |
| **Impacto** | ALTO - Funcionalidad crítica no disponible |

### Descripción del Problema

La función `GetPresignedURL` en `StorageClient` está declarada pero no implementada. Actualmente siempre retorna error.

### Pasos de Implementación

#### Paso 1.1.1: Agregar imports necesarios

Abrir `bootstrap/resource_implementations.go` y agregar al bloque de imports:

```go
import (
    // ... imports existentes ...
    "time"
    
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)
```

#### Paso 1.1.2: Implementar la función

Reemplazar la función actual con:

```go
// GetPresignedURL genera una URL pre-firmada para acceso temporal a un objeto en S3.
//
// Parámetros:
//   - ctx: Contexto para cancelación y timeouts
//   - key: Clave del objeto en el bucket
//   - expirationMinutes: Tiempo de expiración en minutos
//
// Retorna:
//   - URL pre-firmada válida por el tiempo especificado
//   - Error si el presign client no está inicializado o falla la generación
func (c *defaultStorageClient) GetPresignedURL(ctx context.Context, key string, expirationMinutes int) (string, error) {
    if c.presignClient == nil {
        return "", fmt.Errorf("presign client not initialized")
    }
    
    presignClient, ok := c.presignClient.(*s3.PresignClient)
    if !ok {
        return "", fmt.Errorf("presign client is not of type *s3.PresignClient")
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

#### Paso 1.1.3: Agregar test unitario

Crear o agregar en `bootstrap/resource_implementations_test.go`:

```go
func TestDefaultStorageClient_GetPresignedURL(t *testing.T) {
    tests := []struct {
        name          string
        presignClient interface{}
        wantErr       bool
        errContains   string
    }{
        {
            name:          "returns error when presign client is nil",
            presignClient: nil,
            wantErr:       true,
            errContains:   "presign client not initialized",
        },
        {
            name:          "returns error when presign client is wrong type",
            presignClient: "wrong-type",
            wantErr:       true,
            errContains:   "not of type",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            client := &defaultStorageClient{
                client:        nil,
                presignClient: tt.presignClient,
                bucket:        "test-bucket",
            }
            
            _, err := client.GetPresignedURL(context.Background(), "test-key", 15)
            
            if tt.wantErr {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errContains)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Criterios de Éxito

- [ ] La función no retorna "not implemented yet"
- [ ] Los tests unitarios pasan
- [ ] `go build ./...` compila sin errores
- [ ] `go test ./bootstrap/...` pasa

### Comando de Verificación

```bash
cd bootstrap && go test -v -run TestDefaultStorageClient_GetPresignedURL
```

### Commit

```bash
git add bootstrap/resource_implementations.go bootstrap/resource_implementations_test.go
git commit -m "feat(bootstrap): implementar GetPresignedURL en StorageClient

- Implementa generación de URLs pre-firmadas para S3
- Agrega validación de presignClient
- Agrega tests unitarios"
```

### Actualización de Documentación

Actualizar `documents/SERVICES.md` si existe documentación de StorageClient para reflejar que `GetPresignedURL` está disponible y cómo usarlo.

---

## Paso 1.2: Corregir Error Handling en Exists()

### Información General

| Campo | Valor |
|-------|-------|
| **Archivo** | `bootstrap/resource_implementations.go` |
| **Líneas** | 133-145 |
| **Tipo** | Mala práctica |
| **Impacto** | ALTO - Errores silenciados pueden causar bugs en producción |

### Descripción del Problema

La función `Exists()` silencia TODOS los errores, tratándolos como "archivo no existe".

### Pasos de Implementación

#### Paso 1.2.1: Agregar import para tipos de error

En el bloque de imports, agregar:

```go
import (
    // ... imports existentes ...
    "errors"
    
    "github.com/aws/aws-sdk-go-v2/service/s3/types"
)
```

#### Paso 1.2.2: Corregir la función Exists

```go
// Exists verifica si un objeto existe en el bucket.
//
// Parámetros:
//   - ctx: Contexto para cancelación y timeouts
//   - key: Clave del objeto a verificar
//
// Retorna:
//   - true si el objeto existe
//   - false si el objeto no existe (sin error)
//   - false con error si hay problemas de red, permisos, etc.
func (c *defaultStorageClient) Exists(ctx context.Context, key string) (bool, error) {
    _, err := c.client.HeadObject(ctx, &s3.HeadObjectInput{
        Bucket: aws.String(c.bucket),
        Key:    aws.String(key),
    })
    if err != nil {
        // Verificar si es específicamente un error "not found"
        var notFound *types.NotFound
        if errors.As(err, &notFound) {
            return false, nil
        }
        
        // Verificar también NoSuchKey que algunos endpoints usan
        var noSuchKey *types.NoSuchKey
        if errors.As(err, &noSuchKey) {
            return false, nil
        }
        
        // Cualquier otro error se propaga
        return false, fmt.Errorf("failed to check object existence for key %s: %w", key, err)
    }
    return true, nil
}
```

### Criterios de Éxito

- [ ] La función diferencia entre "no existe" y otros errores
- [ ] Errores de red/permisos se propagan correctamente
- [ ] `go build ./...` compila sin errores

### Comando de Verificación

```bash
cd bootstrap && go build ./...
```

### Commit

```bash
git add bootstrap/resource_implementations.go
git commit -m "fix(bootstrap): corregir error handling en StorageClient.Exists()

- Detecta específicamente errores NotFound y NoSuchKey
- Propaga otros errores (network, permisos, etc.)
- Evita falsos negativos por errores silenciados"
```

---

## Paso 1.3: Manejar Errores de Ack/Nack en Consumer

### Información General

| Campo | Valor |
|-------|-------|
| **Archivo** | `messaging/rabbit/consumer.go` |
| **Líneas** | 67-70 |
| **Tipo** | Mala práctica |
| **Impacto** | ALTO - Mensajes pueden perderse o duplicarse |

### Descripción del Problema

Los errores de `Ack()` y `Nack()` se ignoran con `_ =`, lo que puede causar pérdida o duplicación de mensajes.

### Pasos de Implementación

#### Paso 1.3.1: Corregir manejo de Ack/Nack con logging

```go
if !c.config.AutoAck {
    if err != nil {
        // Nack con requeue si hubo error en el handler
        if nackErr := msg.Nack(false, true); nackErr != nil {
            log.Printf("[ERROR] failed to nack message (delivery_tag=%d, queue=%s): %v (original error: %v)",
                msg.DeliveryTag, queueName, nackErr, err)
        }
    } else {
        // Ack si el procesamiento fue exitoso
        if ackErr := msg.Ack(false); ackErr != nil {
            log.Printf("[ERROR] failed to ack message (delivery_tag=%d, queue=%s): %v",
                msg.DeliveryTag, queueName, ackErr)
        }
    }
}
```

### Criterios de Éxito

- [ ] No hay `_ =` para ignorar errores de Ack/Nack
- [ ] Errores se registran en logs
- [ ] `go build ./...` compila sin errores
- [ ] Tests existentes siguen pasando

### Comando de Verificación

```bash
cd messaging/rabbit && go test -v ./...
```

### Commit

```bash
git add messaging/rabbit/consumer.go
git commit -m "fix(messaging): manejar errores de Ack/Nack en consumer

- Registra errores de Ack/Nack en lugar de ignorarlos
- Incluye delivery_tag y queue en logs para debugging
- Facilita diagnóstico de problemas de mensajería"
```

---

## Paso 1.4: Implementar extractEnvAndVersion

### Información General

| Campo | Valor |
|-------|-------|
| **Archivo** | `bootstrap/bootstrap.go` |
| **Líneas** | 431-436 |
| **Tipo** | Código incompleto |
| **Impacto** | MEDIO - Logs no muestran ambiente real |

### Descripción del Problema

La función siempre retorna valores hardcodeados `"local"` y `"0.0.0"`.

### Pasos de Implementación

#### Paso 1.4.1: Implementar con reflection

```go
import "reflect"

// extractEnvAndVersion extrae los campos Environment y Version de una configuración.
//
// Busca campos llamados "Environment" y "Version" en el struct proporcionado.
// Si no los encuentra o el config es nil, retorna valores por defecto.
//
// Parámetros:
//   - config: Struct de configuración (puede ser valor o puntero)
//
// Retorna:
//   - environment: Valor del campo Environment o "unknown"
//   - version: Valor del campo Version o "0.0.0"
func extractEnvAndVersion(config interface{}) (string, string) {
    if config == nil {
        return "unknown", "0.0.0"
    }
    
    v := reflect.ValueOf(config)
    if v.Kind() == reflect.Ptr {
        if v.IsNil() {
            return "unknown", "0.0.0"
        }
        v = v.Elem()
    }
    
    if v.Kind() != reflect.Struct {
        return "unknown", "0.0.0"
    }
    
    // Buscar campo Environment
    env := "unknown"
    envField := v.FieldByName("Environment")
    if envField.IsValid() && envField.Kind() == reflect.String {
        env = envField.String()
        if env == "" {
            env = "unknown"
        }
    }
    
    // Buscar campo Version
    version := "0.0.0"
    versionField := v.FieldByName("Version")
    if versionField.IsValid() && versionField.Kind() == reflect.String {
        ver := versionField.String()
        if ver != "" {
            version = ver
        }
    }
    
    return env, version
}
```

#### Paso 1.4.2: Agregar tests unitarios

```go
func TestExtractEnvAndVersion(t *testing.T) {
    tests := []struct {
        name        string
        config      interface{}
        wantEnv     string
        wantVersion string
    }{
        {
            name:        "nil config returns defaults",
            config:      nil,
            wantEnv:     "unknown",
            wantVersion: "0.0.0",
        },
        {
            name: "struct with Environment and Version",
            config: struct {
                Environment string
                Version     string
            }{
                Environment: "prod",
                Version:     "1.2.3",
            },
            wantEnv:     "prod",
            wantVersion: "1.2.3",
        },
        {
            name: "struct with only Environment",
            config: struct {
                Environment string
            }{
                Environment: "dev",
            },
            wantEnv:     "dev",
            wantVersion: "0.0.0",
        },
        {
            name: "pointer to struct",
            config: &struct {
                Environment string
                Version     string
            }{
                Environment: "qa",
                Version:     "2.0.0",
            },
            wantEnv:     "qa",
            wantVersion: "2.0.0",
        },
        {
            name: "empty environment defaults to unknown",
            config: struct {
                Environment string
            }{
                Environment: "",
            },
            wantEnv:     "unknown",
            wantVersion: "0.0.0",
        },
        {
            name:        "non-struct config returns defaults",
            config:      "not a struct",
            wantEnv:     "unknown",
            wantVersion: "0.0.0",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            gotEnv, gotVersion := extractEnvAndVersion(tt.config)
            assert.Equal(t, tt.wantEnv, gotEnv)
            assert.Equal(t, tt.wantVersion, gotVersion)
        })
    }
}
```

### Criterios de Éxito

- [ ] La función extrae correctamente Environment de structs
- [ ] La función extrae correctamente Version de structs
- [ ] Maneja casos edge (nil, non-struct, punteros)
- [ ] Tests pasan
- [ ] `go build ./...` compila sin errores

### Comando de Verificación

```bash
cd bootstrap && go test -v -run TestExtractEnvAndVersion
```

### Commit

```bash
git add bootstrap/bootstrap.go bootstrap/bootstrap_test.go
git commit -m "feat(bootstrap): implementar extractEnvAndVersion con reflection

- Extrae Environment y Version de cualquier struct config
- Maneja nil, punteros, y tipos no-struct
- Defaults seguros: 'unknown' y '0.0.0'
- Agrega tests unitarios completos"
```

---

## Verificación Final de Fase 1

### Antes de Crear el PR

```bash
# Ejecutar todos los tests
make test-all-modules

# Verificar que no hay TODOs en código modificado
grep -n "TODO" bootstrap/resource_implementations.go
grep -n "TODO" bootstrap/bootstrap.go
grep -n "TODO" messaging/rabbit/consumer.go

# Verificar que no hay errores ignorados
grep -n "_ =" messaging/rabbit/consumer.go

# Build limpio
make build

# Linter
make lint
```

### Crear Pull Request

```bash
# Push de la rama
git push origin fase-1-correcciones-criticas

# En GitHub:
# 1. Crear PR hacia dev
# 2. Título: "feat: Fase 1 - Correcciones Críticas"
# 3. Descripción con lista de cambios
```

### Revisión de GitHub Copilot

| Tipo de Comentario | Acción |
|-------------------|--------|
| Traducción inglés/español | DESCARTAR |
| Error de lógica | CORREGIR |
| Sugerencia de mejora menor | DOCUMENTAR como deuda futura |
| Problema de seguridad | CORREGIR |

### Esperar Pipelines

```bash
# Revisar cada minuto durante máximo 10 minutos
# Si hay errores:
#   1. Analizar causa
#   2. Corregir (máx 3 intentos)
#   3. Push y esperar nuevamente
```

### Criterios de Éxito de Fase

- [ ] Todos los pasos completados
- [ ] `make test-all-modules` pasa
- [ ] `make build` compila sin errores
- [ ] No hay nuevos warnings del linter
- [ ] PR aprobado
- [ ] Pipelines verdes
- [ ] Merge a dev completado

---

## Resumen de la Fase 1

| Paso | Descripción | Commit |
|------|-------------|--------|
| 1.1 | Implementar GetPresignedURL | `feat(bootstrap): implementar GetPresignedURL` |
| 1.2 | Corregir Exists() error handling | `fix(bootstrap): corregir error handling en Exists()` |
| 1.3 | Manejar errores Ack/Nack | `fix(messaging): manejar errores de Ack/Nack` |
| 1.4 | Implementar extractEnvAndVersion | `feat(bootstrap): implementar extractEnvAndVersion` |

---

## Siguiente Fase

Después de completar esta fase y hacer merge a dev, continuar con:
→ [FASE-2: Restauración de Tests](./FASE-2_RESTAURACION_TESTS.md)
