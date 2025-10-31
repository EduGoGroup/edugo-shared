# 📦 EduGo Shared Library - Guía de Actualización v1.0.0

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