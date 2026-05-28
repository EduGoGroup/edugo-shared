# Consumidores por modulo

Esta vista recorre `edugo-shared` desde el lado opuesto: para cada modulo, indica que componentes del ecosistema lo consumen directamente y cual es su papel dentro de esas integraciones.

| Modulo | Consumidores directos verificados | Rol en el ecosistema |
| --- | --- | --- |
| `audit` | IAM, Admin, Mobile | Contrato comun de auditoria para servicios HTTP |
| `audit/postgres` | IAM, Admin, Mobile | Persistencia de auditoria en PostgreSQL |
| `auth` | IAM, Admin, Mobile | JWT, claims y soporte de validacion/autenticacion |
| `bootstrap` | Worker | Construccion de recursos de infraestructura |
| `cache/redis` | Mobile | Cache de datos de aplicacion |
| `common` | IAM, Admin, Mobile, Worker | Errores, enums, UUIDs y helpers comunes |
| `config` | Sin consumidor backend directo detectado | Modulo disponible, no cableado de forma directa en los servicios escaneados |
| `database/mongodb` | Sin consumidor backend directo detectado | Modulo disponible, mientras los servicios actuales usan drivers/repos locales |
| `database/postgres` | Worker | Conexion y utilidades de Postgres en runtime del worker |
| `lifecycle` | Worker | Coordinacion de startup/shutdown |
| `logger` | IAM, Admin, Mobile, Worker | Logging estructurado transversal |
| `messaging/events` | Mobile | Contratos de eventos que luego viajan por RabbitMQ |
| `messaging/rabbit` | Mobile | Publicacion de mensajes RabbitMQ |
| `middleware/gin` | IAM, Admin, Mobile | Middleware compartido de permisos y pipeline HTTP |
| `repository` | IAM, Admin, Mobile | Repositorios GORM compartidos sobre entidades de infraestructura |
| `screenconfig` | Sin consumidor backend directo detectado | Modulo disponible para UI dinamica, sin wiring directo verificado |
| `testing` | Worker | Infraestructura de testcontainers para integracion |

## Notas por modulo

### Modulos HTTP compartidos

- `auth`, `audit`, `audit/postgres`, `middleware/gin`, `logger`, `common` y `repository` forman el nucleo comun de IAM, Admin y Mobile.
- En la practica, ese conjunto define el estilo tecnico transversal del backend HTTP de EduGo.

### Modulos orientados a worker

- `bootstrap`, `database/postgres`, `lifecycle` y `testing` aparecen concentrados en `edugo-worker`.
- Esto indica que el worker esta usando `edugo-shared` mas como runtime de infraestructura que como kit HTTP.

### Modulos con integracion puntual

- `cache/redis`, `messaging/events` y `messaging/rabbit` aparecen solo en `edugo-api-mobile-new` dentro del escaneo actual.
- Esa concentracion encaja con lo descrito en `ecosistema.md`: Mobile API es el servicio con mas dependencias de infraestructura y mas carga funcional.

### Modulos disponibles pero no conectados directamente

- `config`, `database/mongodb` y `screenconfig` no aparecieron como dependencias directas en los servicios backend revisados.
- `screenconfig` tiene sentido funcional en el ecosistema, pero hoy esa logica parece resolverse desde IAM mediante implementaciones locales y repos propios, no mediante import directo del modulo compartido.
- `database/mongodb` existe como libreria reusable, pero Mobile y Worker hoy usan otros caminos de conexion/repositorio en sus propios repositorios.
