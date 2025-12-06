# 📦 EduGo Shared Library - Guía de Actualización

## 🚀 Migrar de v2.0.1 a v2.0.5 (Arquitectura Modular Completa)

### ⚠️ BREAKING CHANGES IMPORTANTES

La versión **v2.0.5** elimina completamente el módulo monolítico `v2` y separa **TODO** en 6 módulos independientes. Este es un cambio **MAJOR** que requiere actualización de imports en todo tu proyecto.

---

### 🎯 Paso 1: Entender la Nueva Arquitectura

#### **Antes (v2.0.1):**
```bash
# Módulo monolítico con TODO incluido
go get github.com/EduGoGroup/edugo-shared/v2@v2.0.1
```
```go
import "github.com/EduGoGroup/edugo-shared/v2/pkg/errors"
import "github.com/EduGoGroup/edugo-shared/v2/pkg/auth"
import "github.com/EduGoGroup/edugo-shared/v2/pkg/logger"
import "github.com/EduGoGroup/edugo-shared/v2/pkg/messaging"
```

**Problema:** Descarga 15+ dependencias (RabbitMQ, JWT, Zap, etc.) aunque solo uses `errors`.

#### **Después (v2.0.5):**
```bash
# Instalación selectiva por módulo
go get github.com/EduGoGroup/edugo-shared/common@v2.0.5
go get github.com/EduGoGroup/edugo-shared/auth@v2.0.5  # Si lo necesitas
```
```go
import "github.com/EduGoGroup/edugo-shared/common/errors"
import "github.com/EduGoGroup/edugo-shared/auth"
```

**Beneficio:** Solo 1-3 dependencias según lo que uses ✅

---

### 📋 Paso 2: Tabla de Migración de Imports

| v2.0.1 (Viejo) | v2.0.5 (Nuevo) | Módulo |
|----------------|----------------|--------|
| `github.com/EduGoGroup/edugo-shared/v2/pkg/errors` | `github.com/EduGoGroup/edugo-shared/common/errors` | `common` |
| `github.com/EduGoGroup/edugo-shared/v2/pkg/types` | `github.com/EduGoGroup/edugo-shared/common/types` | `common` |
| `github.com/EduGoGroup/edugo-shared/v2/pkg/types/enum` | `github.com/EduGoGroup/edugo-shared/common/types/enum` | `common` |
| `github.com/EduGoGroup/edugo-shared/v2/pkg/validator` | `github.com/EduGoGroup/edugo-shared/common/validator` | `common` |
| `github.com/EduGoGroup/edugo-shared/v2/pkg/config` | `github.com/EduGoGroup/edugo-shared/common/config` | `common` |
| `github.com/EduGoGroup/edugo-shared/v2/pkg/auth` | `github.com/EduGoGroup/edugo-shared/auth` | `auth` |
| `github.com/EduGoGroup/edugo-shared/v2/pkg/logger` | `github.com/EduGoGroup/edugo-shared/logger` | `logger` |
| `github.com/EduGoGroup/edugo-shared/v2/pkg/messaging` | `github.com/EduGoGroup/edugo-shared/messaging/rabbit` | `rabbit` |
| `github.com/EduGoGroup/edugo-shared/database/postgres` | Sin cambios ✓ | `postgres` |
| `github.com/EduGoGroup/edugo-shared/database/mongodb` | Sin cambios ✓ | `mongodb` |

---

### 🔧 Paso 3: Actualizar go.mod

**1. Eliminar módulo v2 antiguo:**
```bash
go mod edit -droprequire github.com/EduGoGroup/edugo-shared/v2
```

**2. Agregar solo los módulos que necesites:**
```bash
# Common (errors, types, validator, config) - Casi siempre necesario
go get github.com/EduGoGroup/edugo-shared/common@v2.0.5

# Auth (JWT) - Si usas autenticación
go get github.com/EduGoGroup/edugo-shared/auth@v2.0.5

# Logger (Zap) - Si usas logging
go get github.com/EduGoGroup/edugo-shared/logger@v2.0.5

# RabbitMQ - Si usas messaging
go get github.com/EduGoGroup/edugo-shared/messaging/rabbit@v2.0.5

# PostgreSQL - Si usas Postgres
go get github.com/EduGoGroup/edugo-shared/database/postgres@v2.0.5

# MongoDB - Si usas Mongo
go get github.com/EduGoGroup/edugo-shared/database/mongodb@v2.0.5
```

---

### 🔄 Paso 4: Reemplazar Imports en Tu Código

**Opción A: Buscar/Reemplazar Manual**

En tu editor, busca y reemplaza:
```
v2/pkg/errors        → common/errors
v2/pkg/types         → common/types
v2/pkg/validator     → common/validator
v2/pkg/config        → common/config
v2/pkg/auth          → auth
v2/pkg/logger        → logger
v2/pkg/messaging     → messaging/rabbit
```

**Opción B: Script Automatizado (Bash)**

```bash
#!/bin/bash

# Buscar todos los archivos .go
find . -name "*.go" -type f -exec sed -i '' \
  -e 's|github.com/EduGoGroup/edugo-shared/v2/pkg/errors|github.com/EduGoGroup/edugo-shared/common/errors|g' \
  -e 's|github.com/EduGoGroup/edugo-shared/v2/pkg/types|github.com/EduGoGroup/edugo-shared/common/types|g' \
  -e 's|github.com/EduGoGroup/edugo-shared/v2/pkg/validator|github.com/EduGoGroup/edugo-shared/common/validator|g' \
  -e 's|github.com/EduGoGroup/edugo-shared/v2/pkg/config|github.com/EduGoGroup/edugo-shared/common/config|g' \
  -e 's|github.com/EduGoGroup/edugo-shared/v2/pkg/auth|github.com/EduGoGroup/edugo-shared/auth|g' \
  -e 's|github.com/EduGoGroup/edugo-shared/v2/pkg/logger|github.com/EduGoGroup/edugo-shared/logger|g' \
  -e 's|github.com/EduGoGroup/edugo-shared/v2/pkg/messaging|github.com/EduGoGroup/edugo-shared/messaging/rabbit|g' \
  {} \;

echo "✅ Imports actualizados"
```

---

### ✅ Paso 5: Limpiar y Verificar

```bash
# 1. Limpiar dependencias
go mod tidy

# 2. Verificar que compile
go build ./...

# 3. Ejecutar tests
go test ./...

# 4. Verificar que las dependencias correctas estén en go.mod
cat go.mod | grep edugo-shared
```

**Resultado esperado en go.mod:**
```go
require (
    github.com/EduGoGroup/edugo-shared/common v2.0.5
    github.com/EduGoGroup/edugo-shared/auth v2.0.5
    // ... solo los módulos que uses
)
```

**NO deberías ver:**
```go
github.com/EduGoGroup/edugo-shared/v2 v2.0.1  // ❌ ELIMINAR ESTO
```

---

### 🎯 Ejemplos de Migración

#### Ejemplo 1: Proyecto que solo usa Errors

**Antes:**
```go
// go.mod
require github.com/EduGoGroup/edugo-shared/v2 v2.0.1

// main.go
import "github.com/EduGoGroup/edugo-shared/v2/pkg/errors"
```

**Después:**
```go
// go.mod
require github.com/EduGoGroup/edugo-shared/common v2.0.5

// main.go
import "github.com/EduGoGroup/edugo-shared/common/errors"
```

**Beneficio:** De 15+ deps → 1 dep (ahorro ~93%)

---

#### Ejemplo 2: API con Auth + Postgres + Logger

**Antes:**
```go
// go.mod
require (
    github.com/EduGoGroup/edugo-shared/v2 v2.0.1
    github.com/EduGoGroup/edugo-shared/database/postgres v2.0.1
)

// main.go
import (
    "github.com/EduGoGroup/edugo-shared/v2/pkg/auth"
    "github.com/EduGoGroup/edugo-shared/v2/pkg/logger"
    "github.com/EduGoGroup/edugo-shared/v2/pkg/errors"
    "github.com/EduGoGroup/edugo-shared/database/postgres"
)
```

**Después:**
```go
// go.mod
require (
    github.com/EduGoGroup/edugo-shared/common v2.0.5
    github.com/EduGoGroup/edugo-shared/auth v2.0.5
    github.com/EduGoGroup/edugo-shared/logger v2.0.5
    github.com/EduGoGroup/edugo-shared/database/postgres v2.0.5
)

// main.go
import (
    "github.com/EduGoGroup/edugo-shared/auth"
    "github.com/EduGoGroup/edugo-shared/logger"
    "github.com/EduGoGroup/edugo-shared/common/errors"
    "github.com/EduGoGroup/edugo-shared/database/postgres"
)
```

**Beneficio:** Solo 8 deps en vez de 15+ (ahorro ~47%)

---

### ❓ FAQ

**Q: ¿Puedo mantener v2.0.1 mientras migro?**
A: Sí, pero no es recomendable. v2.0.1 no recibirá actualizaciones futuras.

**Q: ¿Qué pasa si solo uso `common`?**
A: ¡Perfecto! Es el caso de uso ideal. Tendrás mínimas dependencias.

**Q: ¿Los módulos database cambiaron?**
A: No, `database/postgres` y `database/mongodb` mantienen los mismos paths.

**Q: ¿Cómo sé qué módulos necesito?**
A: Revisa tus imports actuales y consulta la tabla de migración arriba.

---

### 🆘 Problemas Comunes

#### Error: "cannot find module"
```bash
# Solución: Asegúrate de instalar el módulo correcto
go get github.com/EduGoGroup/edugo-shared/common@v2.0.5
```

#### Error: "ambiguous import"
```bash
# Solución: Elimina la referencia a v2 en go.mod
go mod edit -droprequire github.com/EduGoGroup/edugo-shared/v2
go mod tidy
```

#### Error: "package ... is not in GOROOT"
```bash
# Solución: Verifica que actualizaste todos los imports
grep -r "v2/pkg/" . --include="*.go"
```

---

## 🚀 Migrar de v1.0.0 a v2.0.0 (Arquitectura Modular)

### ⚠️ BREAKING CHANGES

La versión **v2.0.0** introduce una arquitectura modular con sub-módulos independientes para las bases de datos. Esto **requiere cambios** en tu código.

---

### 🎯 Paso 1: Entender los Cambios

#### **Antes (v1.0.0):**
```bash
# Un solo módulo con todas las dependencias
go get github.com/EduGoGroup/edugo-shared@v1.0.0
```

**Resultado:** Se descargaban drivers de PostgreSQL Y MongoDB (incluso si solo usabas uno).

#### **Después (v2.0.0):**
```bash
# Módulo core (sin bases de datos)
go get github.com/EduGoGroup/edugo-shared@v2.0.0

# Solo el módulo de BD que necesites
go get github.com/EduGoGroup/edugo-shared/database/postgres@v2.0.0
# O
go get github.com/EduGoGroup/edugo-shared/database/mongodb@v2.0.0
```

**Resultado:** Solo descargas las dependencias que realmente necesitas.

---

### 🔄 Paso 2: Actualizar go.mod

#### **Opción A: Usas PostgreSQL**
```bash
cd /path/to/your-project

# Actualizar módulo core
go get github.com/EduGoGroup/edugo-shared@v2.0.0

# Agregar módulo de PostgreSQL
go get github.com/EduGoGroup/edugo-shared/database/postgres@v2.0.0

# Limpiar
go mod tidy
```

#### **Opción B: Usas MongoDB**
```bash
cd /path/to/your-project

# Actualizar módulo core
go get github.com/EduGoGroup/edugo-shared@v2.0.0

# Agregar módulo de MongoDB
go get github.com/EduGoGroup/edugo-shared/database/mongodb@v2.0.0

# Limpiar
go mod tidy
```

#### **Opción C: Usas ambas**
```bash
cd /path/to/your-project

# Actualizar módulo core
go get github.com/EduGoGroup/edugo-shared@v2.0.0

# Agregar ambos módulos
go get github.com/EduGoGroup/edugo-shared/database/postgres@v2.0.0
go get github.com/EduGoGroup/edugo-shared/database/mongodb@v2.0.0

# Limpiar
go mod tidy
```

---

### 📝 Paso 3: Actualizar Imports en tu Código

#### **Cambios requeridos en imports:**

| Antes (v1.0.0) | Después (v2.0.0) |
|----------------|------------------|
| `github.com/EduGoGroup/edugo-shared/pkg/database/postgres` | `github.com/EduGoGroup/edugo-shared/database/postgres` |
| `github.com/EduGoGroup/edugo-shared/pkg/database/mongodb` | `github.com/EduGoGroup/edugo-shared/database/mongodb` |

#### **Ejemplo de migración:**

**Antes (v1.0.0):**
```go
package main

import (
    "github.com/EduGoGroup/edugo-shared/pkg/database/postgres"
    "github.com/EduGoGroup/edugo-shared/pkg/database/mongodb"
    "github.com/EduGoGroup/edugo-shared/pkg/logger"
)

func main() {
    // Usar PostgreSQL
    db, err := postgres.Connect(&cfg)

    // Usar MongoDB
    client, err := mongodb.Connect(mongoCfg)
}
```

**Después (v2.0.0):**
```go
package main

import (
    "github.com/EduGoGroup/edugo-shared/database/postgres"  // ✅ Cambio aquí
    "github.com/EduGoGroup/edugo-shared/database/mongodb"   // ✅ Cambio aquí
    "github.com/EduGoGroup/edugo-shared/pkg/logger"        // Sin cambios
)

func main() {
    // Usar PostgreSQL (API sin cambios)
    db, err := postgres.Connect(&cfg)

    // Usar MongoDB (API sin cambios)
    client, err := mongodb.Connect(mongoCfg)
}
```

---

### 🔍 Paso 4: Buscar y Reemplazar en tu Proyecto

#### **Comando para encontrar todos los archivos que necesitan actualización:**

```bash
# En Linux/Mac
grep -r "pkg/database/postgres" .
grep -r "pkg/database/mongodb" .

# En Windows (PowerShell)
Get-ChildItem -Recurse -Include *.go | Select-String "pkg/database/postgres"
Get-ChildItem -Recurse -Include *.go | Select-String "pkg/database/mongodb"
```

#### **Reemplazo automático (con precaución):**

```bash
# En Linux/Mac
find . -name "*.go" -type f -exec sed -i '' 's|pkg/database/postgres|database/postgres|g' {} \;
find . -name "*.go" -type f -exec sed -i '' 's|pkg/database/mongodb|database/mongodb|g' {} \;

# En Windows (PowerShell)
Get-ChildItem -Recurse -Filter *.go | ForEach-Object {
    (Get-Content $_.FullName) -replace 'pkg/database/postgres', 'database/postgres' | Set-Content $_.FullName
    (Get-Content $_.FullName) -replace 'pkg/database/mongodb', 'database/mongodb' | Set-Content $_.FullName
}
```

---

### ✅ Paso 5: Verificar que Todo Funciona

#### **1. Compilar el proyecto:**
```bash
go build ./...
```

#### **2. Ejecutar tests:**
```bash
go test ./...
```

#### **3. Verificar dependencias:**
```bash
go mod verify
go mod tidy
```

#### **4. Ver el go.mod final:**
```bash
cat go.mod
```

**Deberías ver algo como:**
```go
module github.com/tu-org/tu-proyecto

go 1.25

require (
    github.com/EduGoGroup/edugo-shared v2.0.0
    github.com/EduGoGroup/edugo-shared/database/postgres v2.0.0
    // ...
)
```

---

### 🎁 Paso 6: Beneficios de la Migración

| Aspecto | v1.0.0 | v2.0.0 |
|---------|--------|--------|
| **go.mod** | ~15 dependencias | ~5-8 dependencias |
| **Dependencias descargadas** | Todas las BDs | Solo las que uses |
| **Builds** | Normal | Más rápidos |
| **Flexibilidad** | Baja | Alta |
| **Mantenibilidad** | Monolítica | Modular |

**Ejemplo real:**
- **Proyecto solo con PostgreSQL:**
  - Antes: Descargaba 15 paquetes (incluyendo MongoDB driver)
  - Después: Descarga 8 paquetes (solo PostgreSQL)
  - **Reducción: ~47%** en dependencias

---

### 🚨 Resolución de Problemas

#### **Error: "package not found"**
```bash
# Asegúrate de haber instalado el módulo correcto
go get github.com/EduGoGroup/edugo-shared/database/postgres@v2.0.0
go mod tidy
```

#### **Error: "ambiguous import"**
```bash
# Verifica que no tengas imports mezclados
grep -r "pkg/database" .  # No debería encontrar nada
```

#### **Error: "version conflict"**
```bash
# Forzar versión 2.0.0
go mod edit -require=github.com/EduGoGroup/edugo-shared@v2.0.0
go mod tidy
```

#### **Si necesitas volver a v1.0.0:**
```bash
go get github.com/EduGoGroup/edugo-shared@v1.0.0
# Revertir cambios en imports
git checkout -- .
```

---

### ⏱️ Tiempo Estimado de Migración

| Tamaño del Proyecto | Tiempo Estimado |
|---------------------|-----------------|
| Pequeño (1-5 archivos) | 5-10 minutos |
| Mediano (5-20 archivos) | 10-20 minutos |
| Grande (20+ archivos) | 30-60 minutos |

---

### 📋 Checklist de Migración

- [ ] Actualizar `go.mod` con módulo core v2.0.0
- [ ] Agregar módulo(s) de base de datos v2.0.0
- [ ] Actualizar imports: `pkg/database/postgres` → `database/postgres`
- [ ] Actualizar imports: `pkg/database/mongodb` → `database/mongodb`
- [ ] Ejecutar `go mod tidy`
- [ ] Compilar proyecto: `go build ./...`
- [ ] Ejecutar tests: `go test ./...`
- [ ] Verificar que `go.mod` solo tiene las dependencias necesarias
- [ ] Commit de cambios

---

## 📦 Guía de Actualización v1.0.0 (Legado)

## 🎯 Para Proyectos Consumidores

### 1️⃣ **Actualizar a la nueva versión v1.0.0**

#### **Opción A: Actualizar a versión específica (Recomendado)**
```bash
# En el directorio de tu proyecto que usa edugo-shared
go get github.com/EduGoGroup/edugo-shared@v1.0.0
```

#### **Opción B: Actualizar a la versión más reciente**
```bash
go get -u github.com/EduGoGroup/edugo-shared
```

#### **Opción C: Ver versiones disponibles**
```bash
go list -m -versions github.com/EduGoGroup/edugo-shared
```

### 2️⃣ **Verificar la actualización**

```bash
# Verificar que la versión se actualizó
go list -m github.com/EduGoGroup/edugo-shared

# Limpiar caché de módulos si es necesario
go mod tidy
```

### 3️⃣ **Cambios en el código (si aplica)**

#### **✅ COMPATIBILIDAD: No hay breaking changes**
- La versión v1.0.0 es **100% compatible** con v0.1.0
- No necesitas cambiar tu código existente
- Todas las APIs mantienen la misma signature

#### **🚀 Nuevas funcionalidades disponibles:**

##### **JWT Authentication (Mejorado)**
```go
import "github.com/EduGoGroup/edugo-shared/pkg/auth"

// Funcionalidad existente sigue igual
manager := auth.NewJWTManager(secretKey, issuer)
token, err := manager.GenerateToken(userID, email, role, expiresIn)

// ✨ Nuevas capacidades agregadas sin cambios en API
```

##### **Database Connections (Mejorado)**
```go
import "github.com/EduGoGroup/edugo-shared/pkg/database/postgres"
import "github.com/EduGoGroup/edugo-shared/pkg/database/mongodb"

// PostgreSQL con mejor configuración
cfg := postgres.DefaultConfig()
cfg.Host = "localhost"
cfg.MaxConnections = 25  // ✨ Nuevos defaults optimizados
db, err := postgres.Connect(&cfg)  // ✨ Ahora usa puntero (más eficiente)

// MongoDB con pools optimizados
mongoCfg := mongodb.DefaultConfig()
mongoCfg.MaxPoolSize = 100  // ✨ Defaults profesionales
client, err := mongodb.Connect(mongoCfg)
```

##### **Error Handling (Mejorado)**
```go
import "github.com/EduGoGroup/edugo-shared/pkg/errors"

// ✨ Mejor alineación de memoria en structs
appErr := errors.NewValidationError("invalid input")
// Misma API, mejor performance interno
```

### 4️⃣ **Validar que todo funciona**

#### **Ejecutar tests**
```bash
# En tu proyecto
go test ./...
```

#### **Verificar build**
```bash
go build ./...
```

#### **Verificar imports**
```bash
go mod tidy
go mod verify
```

### 5️⃣ **Aprovechar nuevas funcionalidades**

#### **✨ Nuevos comandos Make (si adoptas el patrón)**
```bash
# Copia el Makefile de edugo-shared a tu proyecto para:
make lint          # Linting profesional
make test-coverage # Tests con coverage
make security      # Análisis de seguridad
make fmt           # Formateo automático
```

#### **✨ Nuevas validaciones disponibles**
```go
import "github.com/EduGoGroup/edugo-shared/pkg/validator"

validator := validator.New()
validator.Email("user@example.com", "email")  // ✨ Mejorado
validator.UUID("123e4567-e89b-12d3-a456-426614174000", "id")  // ✨ Nuevo
validator.URL("https://example.com", "website")  // ✨ Nuevo
```

#### **✨ Logging estructurado**
```go
import "github.com/EduGoGroup/edugo-shared/pkg/logger"

// ✨ Nueva funcionalidad disponible
logger := logger.NewZapLogger("production")
logger.Info("User logged in", "userID", userID)
```

---

## 🔄 **Proceso Completo de Actualización**

### **Para un proyecto típico:**

```bash
# 1. Navegar al proyecto
cd /path/to/your-edugo-project

# 2. Actualizar dependencia
go get github.com/EduGoGroup/edugo-shared@v1.0.0

# 3. Limpiar módulos
go mod tidy

# 4. Verificar que compila
go build ./...

# 5. Ejecutar tests
go test ./...

# 6. Commit de la actualización
git add go.mod go.sum
git commit -m "chore: update edugo-shared to v1.0.0"
git push
```

---

## 🚨 **Resolución de Problemas**

### **Si encuentras errores de compilación:**
```bash
# Limpiar caché de módulos
go clean -modcache
go mod download
go mod tidy
```

### **Si hay conflictos de versiones:**
```bash
# Ver todas las dependencias
go mod graph | grep edugo-shared

# Forzar versión específica
go mod edit -require=github.com/EduGoGroup/edugo-shared@v1.0.0
go mod tidy
```

### **Si necesitas volver a versión anterior:**
```bash
# Downgrade temporal
go get github.com/EduGoGroup/edugo-shared@v0.1.0
```

---

## 📞 **Soporte**

- **Documentación**: Revisa el CHANGELOG.md en el repositorio
- **Issues**: Crea un issue en GitHub si encuentras problemas
- **Tests**: Los tests del edugo-shared cubren el 87.2% del código

---

## 🎉 **Beneficios de v1.0.0**

✅ **Estabilidad**: API garantizada compatible hacia adelante
✅ **Performance**: Optimizaciones de memoria y código
✅ **Calidad**: 100% linter compliance, 0 warnings
✅ **Tooling**: Makefile profesional y CI/CD
✅ **Testing**: Coverage mejorado y tests exhaustivos
✅ **Documentation**: Documentación completa de todos los packages

**¡Tu proyecto ahora usa una librería de calidad profesional!** 🚀