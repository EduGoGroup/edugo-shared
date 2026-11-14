# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/cache/redis` deben registrarse aqui.

## [Unreleased]

## [0.900.0] - 2026-11-14

### Changed
- Directiva `go 1.27.4` (plan 054, homologación de Go y dependencias).
- Dependencias directas: `github.com/alicebob/miniredis/v2` `v2.37.0` → `v2.40.0`, `github.com/redis/go-redis/v9` `v9.18.0` → `v9.23.0`.
- 2 dependencias indirectas suben con `go mod tidy` (ver `go.mod`).
- El módulo entra a la banda `v0.900.x` (ADR 0010); el tag anterior `v0.1.x` queda como procedencia.

## [0.1.0] - 2026-05-28

### Added
- Reinicio de la versión del módulo a `v0.1.0` (borrón y cuenta nueva).
- Conservación del código de producción estable del módulo.

## [0.100.0] - 2026-04-02

### Changed

- Added `//go:build integration` tag to `cache_test.go`

### Added

- **Redis Connection**: Módulo de bajo nivel para conectar a Redis con validación de conectividad
- **RedisConfig**: Configuración centralizada con URL (soporte redis:// y rediss://+TLS)
- **ConnectRedis()**: Establece conexión a Redis con validación mediante Ping (5s timeout)
- **CacheService Interface**: Interfaz genérica agnóstica de Redis para operaciones de caché
- **NewCacheService()**: Crea servicio de caché respaldado por Redis client
- **Get()**: Obtiene y deserializa valores JSON del caché
- **Set()**: Serializa a JSON y guarda en Redis con TTL configurable
- **Delete()**: Borra una o múltiples claves específicas
- **DeleteByPattern()**: Borra claves por patrón usando SCAN (batch 100)
- **JSON Automático**: Serialización/deserialización transparente en Set/Get
- **TTL Flexible**: TTL configurable por operación, no global
- **SCAN Seguro**: DeleteByPattern usa SCAN para evitar bloqueos con KEYS
- **TLS Support**: rediss:// URL para conexiones a Upstash y similares
- **Concurrency Safe**: Client thread-safe, seguro compartir entre goroutines
- **Context Aware**: Todas las operaciones respetan contexto y timeouts
- **Error Handling**: Errores contextuales con wrapping ("parsing redis URL", "pinging redis")
- **Defensive Defaults**: Validación temprana en Connect, no lazy connection
- **Suite completa de tests unitarios e integración** sin dependencias externas para unitarios
- **Documentación técnica detallada** en docs/README.md con componentes, flujos comunes y patrones
- **Makefile** con targets: fmt, vet, lint, test, build, check

### Design Notes

- **Interfaz genérica**: CacheService abstrae Redis, usuario no interactúa directamente con client
- **Bajo nivel**: Solo conecta y provee caché, no define patrones de uso
- **JSON agnóstico**: Serialización automática, soporta cualquier tipo que sea JSON-serializable
- **TTL dinámico**: Configurable por operación, permite patrones de expiración complejos
- **SCAN eficiente**: Batch size de 100, no bloquea Redis en patrones amplios
- **TLS flexible**: rediss:// para producción, redis:// para desarrollo
- **Sin pooling adicional**: redis-go client maneja pool internamente (100 conexiones default)
- **Contexto sensible**: Cancellation y timeouts respetados en todas las operaciones
- **Errores transparentes**: Wrapped con contexto, redis.Nil distinguible de otros errores
- **Concurrencia**: Client thread-safe, seguro compartir en múltiples goroutines
