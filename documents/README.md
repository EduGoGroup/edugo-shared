# EduGo Shared - Documentación Completa

> Librería compartida de infraestructura para el ecosistema EduGo

## Índice de Documentación

### Documentación Principal

| Documento | Descripción | Audiencia |
|-----------|-------------|-----------|
| [Arquitectura General](./ARCHITECTURE.md) | Diagrama y descripción de la arquitectura del proyecto, patrones de diseño, flujo de inicialización | Todos |
| [Diagrama de Base de Datos](./DATABASE.md) | Esquema y configuración de PostgreSQL y MongoDB, pool de conexiones, health checks | Backend |
| [Módulos del Sistema](./MODULES.md) | Descripción detallada de cada uno de los 10 módulos con ejemplos de código | Backend |
| [Flujo de Procesos](./PROCESSES.md) | Diagramas de flujo y secuencia para autenticación, bootstrap, mensajería, lifecycle | Todos |
| [Servicios Externos](./SERVICES.md) | Dependencias y servicios necesarios, Docker Compose completo, variables de entorno | DevOps/Backend |
| [API de Interfaces](./INTERFACES.md) | Interfaces públicas exportadas, métodos, parámetros, retornos | Backend |
| [Guía de Configuración](./CONFIGURATION.md) | Variables de entorno, archivos YAML, validación, configuración por ambiente | DevOps/Backend |
| [Testing](./TESTING.md) | Guía de testing con testcontainers, fixtures, cleanup, troubleshooting | QA/Backend |
| [Tipos y Enums](./TYPES.md) | Tipos compartidos, enumeraciones del dominio, códigos de error | Backend |

### Mejoras y Deuda Técnica

| Documento | Descripción | Prioridad |
|-----------|-------------|-----------|
| [Mejoras - Índice](./mejoras/README.md) | Resumen ejecutivo de todas las mejoras pendientes | - |
| [Código Incompleto](./mejoras/CODIGO_INCOMPLETO.md) | Funciones con TODOs y código sin implementar | Alta |
| [Tests Skipped](./mejoras/TESTS_SKIPPED.md) | Tests deshabilitados que necesitan arreglarse | Alta |
| [Malas Prácticas](./mejoras/MALAS_PRACTICAS.md) | Código con malas prácticas a corregir | Media |
| [Refactoring](./mejoras/REFACTORING.md) | Código que necesita reestructuración | Media |
| [Deuda Técnica](./mejoras/DEUDA_TECNICA.md) | Deuda técnica general del proyecto | Baja |

---

## Visión General

**edugo-shared** es una librería compartida escrita en Go que proporciona componentes de infraestructura reutilizables para los microservicios del ecosistema EduGo. El proyecto sigue una arquitectura modular con cada componente como un submódulo Go independiente.

### Propósito

- **Centralizar** código de infraestructura común
- **Estandarizar** conexiones a bases de datos, messaging y autenticación
- **Facilitar** el testing con containers reutilizables
- **Garantizar** consistencia en el manejo de errores y logging

### Stack Tecnológico

| Componente | Tecnología |
|------------|------------|
| Lenguaje | Go 1.25+ |
| Base de Datos Relacional | PostgreSQL 15+ |
| Base de Datos NoSQL | MongoDB 7.0+ |
| Message Broker | RabbitMQ 3.12+ |
| Object Storage | AWS S3 |
| Framework HTTP | Gin |
| ORM | GORM |
| Logging | Logrus / Zap |
| Testing | testcontainers-go |

---

## Estructura del Proyecto

```
edugo-shared/
├── auth/                    # Autenticación JWT
├── bootstrap/               # Inicialización de aplicaciones
├── common/                  # Código común compartido
│   ├── errors/             # Manejo de errores estandarizado
│   ├── types/              # Tipos compartidos
│   │   └── enum/           # Enumeraciones del sistema
│   └── validator/          # Validación de datos
├── config/                  # Carga y validación de configuración
├── database/                # Conexiones a bases de datos
│   ├── mongodb/            # Cliente MongoDB
│   └── postgres/           # Cliente PostgreSQL
├── lifecycle/               # Gestión de ciclo de vida
├── logger/                  # Logging estructurado
├── messaging/               # Sistema de mensajería
│   └── rabbit/             # Cliente RabbitMQ
├── middleware/              # Middlewares HTTP
│   └── gin/                # Middlewares para Gin
├── testing/                 # Infraestructura de testing
│   └── containers/         # Testcontainers
└── documents/               # Esta documentación
```

---

## Inicio Rápido

### Instalación

```bash
# Instalar módulo principal
go get github.com/EduGoGroup/edugo-shared@latest

# O módulos específicos
go get github.com/EduGoGroup/edugo-shared/auth@latest
go get github.com/EduGoGroup/edugo-shared/database/postgres@latest
```

### Ejemplo Básico de Uso

```go
package main

import (
    "context"
    "log"
    
    "github.com/EduGoGroup/edugo-shared/bootstrap"
    "github.com/EduGoGroup/edugo-shared/config"
    "github.com/EduGoGroup/edugo-shared/lifecycle"
)

func main() {
    ctx := context.Background()
    
    // Cargar configuración
    cfg, err := config.Load("config.yaml")
    if err != nil {
        log.Fatal(err)
    }
    
    // Crear lifecycle manager
    lm := lifecycle.NewManager(logger)
    defer lm.Cleanup()
    
    // Bootstrap de recursos
    resources, err := bootstrap.Bootstrap(ctx, cfg, factories, lm,
        bootstrap.WithRequiredResources("logger", "postgresql"),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Usar recursos...
    db := resources.PostgreSQL
}
```

---

## Comandos Make Disponibles

```bash
make help                    # Ver todos los comandos
make setup                   # Configurar entorno
make test                    # Ejecutar tests
make test-race               # Tests con race detection
make test-coverage           # Tests con cobertura
make lint                    # Ejecutar linter
make build                   # Verificar compilación
make test-all-modules        # Tests en todos los módulos
make ci                      # Pipeline CI completo
```

---

## Guía de Uso por Rol

### Para Nuevos Desarrolladores

1. **Primer día**: Lee `README.md` y `ARCHITECTURE.md`
2. **Segundo día**: Configura ambiente con `SERVICES.md` y `CONFIGURATION.md`
3. **Tercer día**: Explora `MODULES.md` para entender cada componente
4. **Ongoing**: Usa `INTERFACES.md` como referencia de API

### Para Tech Leads

1. Revisa `mejoras/README.md` para planificar sprints de deuda técnica
2. Usa `ARCHITECTURE.md` para onboarding de equipo
3. Consulta `TESTING.md` para estándares de QA

### Para DevOps

1. `SERVICES.md` tiene Docker Compose completo
2. `CONFIGURATION.md` documenta todas las variables de entorno
3. `DATABASE.md` para tuning de conexiones

---

## Métricas del Proyecto

### Estadísticas de Código

| Métrica | Valor |
|---------|-------|
| Módulos | 10 |
| Líneas de código | ~15,000 |
| Cobertura de tests | >80% |
| Archivos de documentación | 15 |

### Módulos por Categoría

```
CORE (5 módulos)
├── auth          → Autenticación JWT y passwords
├── config        → Carga de configuración
├── logger        → Logging estructurado
├── lifecycle     → Gestión de ciclo de vida
└── common        → Tipos, errores, validación

INFRASTRUCTURE (4 módulos)
├── database/postgres  → Cliente PostgreSQL
├── database/mongodb   → Cliente MongoDB
├── messaging/rabbit   → Publisher/Consumer RabbitMQ
└── bootstrap          → Inicialización de apps

SUPPORT (2 módulos)
├── middleware/gin     → Middlewares HTTP
└── testing/containers → Testcontainers
```

---

## Preguntas Frecuentes (FAQ)

### ¿Cómo agrego una nueva dependencia?

```bash
cd <módulo>
go get <dependencia>
go mod tidy
```

### ¿Cómo ejecuto tests de un módulo específico?

```bash
cd <módulo>
go test -v ./...
```

### ¿Cómo veo la cobertura de código?

```bash
make coverage-all-modules
```

### ¿Dónde reporto un bug o mejora?

1. Si es código incompleto → `mejoras/CODIGO_INCOMPLETO.md`
2. Si es mala práctica → `mejoras/MALAS_PRACTICAS.md`
3. Si necesita refactoring → `mejoras/REFACTORING.md`

### ¿Cómo contribuyo?

1. Crea branch desde `dev`
2. Implementa cambios
3. Asegura tests pasan: `make test-all-modules`
4. Crea PR hacia `dev`

---

## Contacto y Soporte

- **Repositorio**: github.com/EduGoGroup/edugo-shared
- **Issues**: Usar GitHub Issues
- **Documentación**: Este directorio `/documents`

---

## Historial de Cambios de Documentación

| Fecha | Cambio |
|-------|--------|
| 2024-12-06 | Creación inicial de documentación completa |
| 2024-12-06 | Agregada carpeta de mejoras con análisis de código |

---

## Licencia

Uso interno EduGo - Todos los derechos reservados
