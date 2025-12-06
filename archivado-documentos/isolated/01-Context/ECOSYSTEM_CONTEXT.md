# ECOSYSTEM CONTEXT - SHARED

## Posición en EduGo

**Rol:** Librería de base - Fundación de todo EduGo  
**Interacción:** Importada por TODOS los otros proyectos

---

## Mapa de Ecosistema

```
┌─────────────────────────────────────────┐
│      SHARED Library (Fundación)         │
│  ├─ Logger                              │
│  ├─ Database (PostgreSQL, MongoDB)      │
│  ├─ Auth (JWT)                          │
│  ├─ Messaging (RabbitMQ)                │
│  ├─ Models (User, School, etc)          │
│  ├─ Context (Timeouts, User Info)       │
│  ├─ Errors (Estandarizados)             │
│  └─ Health (Checks)                     │
└──────────┬──────────────────────────────┘
           │
    ┌──────┼──────┬──────┬──────┐
    │      │      │      │      │
    ▼      ▼      ▼      ▼      ▼
  API    API   WORKER  DEV-   (Otros)
 MOBILE ADMIN         ENV
  v1.0  v1.0  v1.0  v1.0
```

---

## Dependencia de Todos Otros Proyectos

### En api-mobile/go.mod
```
require github.com/EduGoGroup/edugo-shared v1.3.0+
```

### En api-admin/go.mod
```
require github.com/EduGoGroup/edugo-shared v1.3.0+
```

### En worker/go.mod
```
require github.com/EduGoGroup/edugo-shared v1.4.0+
```

---

## Versionamiento y Compatibilidad

### Política
- SHARED v1.x = API Mobile v1.x compatible
- SHARED v1.x = API Admin v1.x compatible
- SHARED v1.4.0+ = Worker v1.x compatible
- MAJOR version bump = Breaking changes (nuevo repo/tag)

### Matriz de Compatibilidad

| SHARED | API Mobile | API Admin | Worker |
|--------|-----------|----------|--------|
| v1.0.0 | v0.1.0 | N/A | N/A |
| v1.1.0 | v0.2.0 | N/A | N/A |
| v1.2.0 | v0.3.0 | N/A | N/A |
| v1.3.0 | v1.0.0 | v1.0.0 | v0.1.0 |
| v1.4.0 | v1.1.0 | v1.1.0 | v1.0.0 |

---

## Cambios en SHARED que Afectan Otros Proyectos

### Si Logger cambia
- Todos deben actualizar imports
- Potencial recompilación de todos proyectos

### Si Database module cambia
- API Mobile debe testear
- API Admin debe testear
- Worker debe testear
- Potencial cambios de schema

### Si Auth module cambia
- Todos deben actualizar middleware
- Validación de JWT puede fallar

### Si Messaging module cambia
- API Mobile (publisher/consumer)
- Worker (consumer/publisher)
- Cambio en formato de mensaje = incompatibilidad

---

## Release Workflow

```
SHARED Development:
1. Cambios en rama dev
2. Tests pasen (>80% coverage)
3. Tag v1.3.0
4. Push a GitHub

Otros Proyectos:
5. Notificación de nueva versión
6. Actualizar go.mod
7. Testing con nueva versión
8. Merge a dev
9. Deploy

Coordinación:
- Changelog detallado en SHARED
- Breaking changes anunciadas
- Migration guide si es necesario
```

---

## Dependencias de SHARED

SHARED depende de estas librerías externas (que heredan todos los proyectos):

```
├─ Logger
│  └─ github.com/sirupsen/logrus
├─ Database PostgreSQL
│  ├─ gorm.io/gorm
│  └─ gorm.io/driver/postgres
├─ Database MongoDB
│  └─ go.mongodb.org/mongo-driver
├─ Auth
│  └─ github.com/golang-jwt/jwt/v5
├─ Messaging
│  └─ github.com/streadway/amqp
└─ Configuration
   └─ github.com/spf13/viper
```

---

## Impacto de Cambios

### Si Versión de PostgreSQL sube (15→16)
- SHARED debe adaptar drivers
- Todos proyectos heredan cambio
- Testing requerido

### Si se agrega nuevo módulo a SHARED
- API Mobile puede usar opcionalmente
- API Admin puede usar opcionalmente
- Worker puede usar opcionalmente
- No breaking change (backward compatible)

### Si se remueve módulo de SHARED
- MAJOR version bump (v2.0.0)
- Todos proyectos que lo usan deben migrar
- Requires coordinación

---

## Puntos de Sincronización Críticos

```
Lunes (Planning):
- Revisar qué proyectos necesitan cambios de SHARED
- Asignar developer a SHARED si cambios requeridos

Miércoles (Mid-sprint):
- SHARED si hay cambios, debe estar stable
- Otros proyectos pueden empezar testing

Viernes (Release):
- Tag v1.x.x en SHARED
- Notificar otros proyectos
- Otros proyectos actualizar si es necesario
```

---

## Testing Cross-Project

```bash
# Después de cambios en SHARED:

# Test API Mobile
cd api-mobile
go mod tidy
go get -u github.com/EduGoGroup/edugo-shared@dev
go test ./...

# Test API Admin
cd api-admin
go mod tidy
go get -u github.com/EduGoGroup/edugo-shared@dev
go test ./...

# Test Worker
cd worker
go mod tidy
go get -u github.com/EduGoGroup/edugo-shared@dev
go test ./...
```

---

## Checklist de Integración

- [ ] Todos los módulos importables
- [ ] Compatibilidad con Go 1.21+
- [ ] Tests unitarios >= 80%
- [ ] Documentación completa
- [ ] Ejemplos funcionales
- [ ] Versionamiento semántico
- [ ] Compatible con proyectos (v0.x y v1.x)
