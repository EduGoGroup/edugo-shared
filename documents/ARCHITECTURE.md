# Arquitectura del Sistema

## Diagrama de Arquitectura General

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              MICROSERVICIOS EDUGO                               │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐            │
│  │  Auth API   │  │ Materials   │  │  Progress   │  │   Users     │   ...      │
│  │  Service    │  │  Service    │  │  Service    │  │  Service    │            │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘            │
└─────────┼────────────────┼────────────────┼────────────────┼────────────────────┘
          │                │                │                │
          └────────────────┴────────────────┴────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           EDUGO-SHARED LIBRARY                                  │
│                                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                           BOOTSTRAP LAYER                                │   │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌─────────────┐  │   │
│  │  │   Logger     │  │  PostgreSQL  │  │   MongoDB    │  │  RabbitMQ   │  │   │
│  │  │   Factory    │  │   Factory    │  │   Factory    │  │   Factory   │  │   │
│  │  └──────────────┘  └──────────────┘  └──────────────┘  └─────────────┘  │   │
│  │                                                        ┌─────────────┐  │   │
│  │                                                        │  S3 Factory │  │   │
│  │                                                        └─────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                          CORE MODULES                                    │   │
│  │                                                                          │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │   │
│  │  │    auth     │  │   config    │  │   logger    │  │  lifecycle  │     │   │
│  │  │  JWT/Pass   │  │   Loader    │  │  Zap/Logrus │  │   Manager   │     │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘     │   │
│  │                                                                          │   │
│  │  ┌─────────────────────┐  ┌─────────────────────┐  ┌─────────────┐      │   │
│  │  │      database       │  │      messaging      │  │  middleware │      │   │
│  │  │  ┌───────┬───────┐  │  │  ┌───────────────┐  │  │  ┌───────┐  │      │   │
│  │  │  │Postgre│MongoDB│  │  │  │    RabbitMQ   │  │  │  │  Gin  │  │      │   │
│  │  │  └───────┴───────┘  │  │  │ Pub/Sub/DLQ   │  │  │  │ JWT   │  │      │   │
│  │  └─────────────────────┘  └─────────────────────┘  └─────────────┘      │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                          COMMON LAYER                                    │   │
│  │  ┌─────────────┐  ┌─────────────────────────┐  ┌─────────────┐          │   │
│  │  │   errors    │  │         types           │  │  validator  │          │   │
│  │  │  AppError   │  │  ┌────────┐ ┌────────┐  │  │  Input Val  │          │   │
│  │  │  Codes      │  │  │  enum  │ │  uuid  │  │  │             │          │   │
│  │  └─────────────┘  │  └────────┘ └────────┘  │  └─────────────┘          │   │
│  │                   └─────────────────────────┘                            │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                          TESTING LAYER                                   │   │
│  │  ┌───────────────────────────────────────────────────────────────────┐  │   │
│  │  │                        containers                                  │  │   │
│  │  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐               │  │   │
│  │  │  │ PostgreSQL  │  │   MongoDB   │  │  RabbitMQ   │               │  │   │
│  │  │  │  Container  │  │  Container  │  │  Container  │               │  │   │
│  │  │  └─────────────┘  └─────────────┘  └─────────────┘               │  │   │
│  │  └───────────────────────────────────────────────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────────┘
                                    │
          ┌─────────────────────────┼─────────────────────────┐
          │                         │                         │
          ▼                         ▼                         ▼
┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│   PostgreSQL    │       │     MongoDB     │       │    RabbitMQ     │
│    Database     │       │    Database     │       │  Message Broker │
└─────────────────┘       └─────────────────┘       └─────────────────┘
                                    │
                                    ▼
                          ┌─────────────────┐
                          │     AWS S3      │
                          │  Object Storage │
                          └─────────────────┘
```

---

## Capas del Sistema

### 1. Bootstrap Layer

La capa de bootstrap es el punto de entrada principal para inicializar todos los recursos de infraestructura. Implementa el **Factory Pattern** para crear conexiones a servicios externos.

**Responsabilidades:**
- Inicialización ordenada de recursos (Logger → DB → Messaging → Storage)
- Health checks automáticos
- Registro de cleanup en lifecycle manager
- Soporte para recursos requeridos vs opcionales

### 2. Core Modules

Módulos fundamentales que proporcionan funcionalidad específica:

| Módulo | Función |
|--------|---------|
| **auth** | Autenticación JWT y hashing de contraseñas |
| **config** | Carga de configuración YAML con validación |
| **logger** | Logging estructurado (Zap/Logrus) |
| **lifecycle** | Gestión de startup/shutdown de recursos |
| **database** | Conexiones PostgreSQL y MongoDB |
| **messaging** | Publisher/Consumer RabbitMQ |
| **middleware** | Middlewares HTTP para Gin |

### 3. Common Layer

Código compartido entre todos los módulos:

| Componente | Descripción |
|------------|-------------|
| **errors** | Sistema de errores estructurados con códigos HTTP |
| **types/enum** | Enumeraciones del dominio (roles, status, etc.) |
| **types/uuid** | Helpers para generación de UUIDs |
| **validator** | Validación de datos de entrada |

### 4. Testing Layer

Infraestructura de testing con testcontainers:

- Containers singleton reutilizables
- Cleanup automático entre tests
- Soporte para PostgreSQL, MongoDB y RabbitMQ

---

## Patrones de Diseño Utilizados

### Factory Pattern
```
Bootstrap → Factories → Resources
```
Los factories crean conexiones a servicios externos de manera consistente.

### Singleton Pattern
```go
manager, _ := containers.GetManager(t, nil)  // Siempre retorna la misma instancia
```
Usado en testing para reutilizar containers.

### LIFO Cleanup
```
Register: A → B → C
Cleanup:  C → B → A
```
El lifecycle manager limpia recursos en orden inverso al registro.

### Interface Segregation
```go
type Logger interface {
    Debug(msg string, fields ...interface{})
    Info(msg string, fields ...interface{})
    // ...
}
```
Interfaces pequeñas y específicas permiten múltiples implementaciones.

---

## Flujo de Inicialización

```
┌─────────────────────────────────────────────────────────────────┐
│                    BOOTSTRAP SEQUENCE                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  1. Load Configuration                                          │
│     └─► config.Load() → BaseConfig                              │
│                                                                  │
│  2. Create Lifecycle Manager                                     │
│     └─► lifecycle.NewManager(logger)                            │
│                                                                  │
│  3. Bootstrap Resources (ordered)                                │
│     │                                                            │
│     ├─► [1] Logger Factory → Logger                             │
│     │                                                            │
│     ├─► [2] PostgreSQL Factory → *gorm.DB                       │
│     │       └─► Register cleanup                                 │
│     │                                                            │
│     ├─► [3] MongoDB Factory → *mongo.Client                     │
│     │       └─► Register cleanup                                 │
│     │                                                            │
│     ├─► [4] RabbitMQ Factory → Connection + Channel             │
│     │       └─► Register cleanup                                 │
│     │                                                            │
│     └─► [5] S3 Factory → S3 Client                              │
│                                                                  │
│  4. Health Checks                                                │
│     ├─► PostgreSQL: SELECT 1                                    │
│     └─► MongoDB: Ping                                           │
│                                                                  │
│  5. Return Resources                                             │
│     └─► &Resources{PostgreSQL, MongoDB, MessagePublisher, ...}  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Comunicación entre Componentes

```
┌─────────────────┐     Publish      ┌─────────────────┐
│   Service A     │ ───────────────► │    RabbitMQ     │
│  (Producer)     │                  │    Exchange     │
└─────────────────┘                  └────────┬────────┘
                                              │
                                     Routing Key
                                              │
                           ┌──────────────────┼──────────────────┐
                           │                  │                  │
                           ▼                  ▼                  ▼
                    ┌───────────┐      ┌───────────┐      ┌───────────┐
                    │  Queue A  │      │  Queue B  │      │   DLQ     │
                    └─────┬─────┘      └─────┬─────┘      └───────────┘
                          │                  │               (errores)
                          ▼                  ▼
                    ┌───────────┐      ┌───────────┐
                    │ Consumer  │      │ Consumer  │
                    │     A     │      │     B     │
                    └───────────┘      └───────────┘
```

---

## Ambientes Soportados

| Ambiente | Descripción | Configuración |
|----------|-------------|---------------|
| `local` | Desarrollo local | Docker containers locales |
| `dev` | Desarrollo compartido | Infraestructura de desarrollo |
| `qa` | Testing/QA | Infraestructura de pruebas |
| `prod` | Producción | Infraestructura productiva |

La configuración se valida contra estos valores:
```go
Environment string `validate:"required,oneof=local dev qa prod"`
```
