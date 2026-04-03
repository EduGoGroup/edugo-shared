# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/middleware/gin` deben registrarse aqui.

## [Unreleased]

## [0.101.0] - 2026-04-02

### Added

- Middleware CORS compartido (`CORSMiddleware`) con `CORSConfig` para unificar comportamiento entre servicios consumidores.
- Helper `BindJSON` para binding + validación con errores por campo usando `common/errors`.
- Nuevas pruebas para `auth_client`, `bind` y `cors` para reforzar cobertura del módulo.

### Changed

- CORS: en entornos no `development/local`, el wildcard se maneja en modo fail-closed sin terminar el proceso.
- CORS: se evita la combinación inválida `Access-Control-Allow-Origin: *` con `Access-Control-Allow-Credentials: true` reflejando `Origin` cuando corresponde.
- CORS: `Vary` ahora se fusiona (append) en lugar de sobrescribirse.
- Bind: los errores de validación priorizan el nombre real del campo desde el tag `json` (fallback a `snake_case`).
- Cobertura del módulo actualizada para cumplir umbral del pipeline (>=95%).

## [0.100.0] - 2026-04-02

### Added

- **JWT Authentication**: Validación segura de tokens JWT con `JWTAuthMiddleware` y variante con blacklist (`JWTAuthMiddlewareWithBlacklist`)
- **Permission Authorization**: Validación granular con `RequirePermission` (individual), `RequireAnyPermission` (OR lógico), `RequireAllPermissions` (AND lógico)
- **Request Logging**: Middleware `RequestLogging` con generación automática de request_id y correlation_id, enriquecimiento de contexto con logger estructurado
- **Post-Auth Logging**: Middleware `PostAuthLogging` que enriquece logs con user_id, role, school_id después de validación JWT
- **Audit Logging**: Middleware `AuditMiddleware` que registra automáticamente operaciones mutantes (POST, PUT, PATCH, DELETE), extrae resource_type e resource_id del path
- **Context Helpers**: Extractores seguros `GetUserID`, `GetEmail`, `GetRole`, `GetClaims` y variantes `Must*` para acceso a datos de autenticación
- **List Filters**: Función `ParseListFilters` para parseo defensivo de parámetros de paginación (page, limit), búsqueda (search, search_fields), filtrado booleano (is_active) y campos extra
- **Logger Integration**: `GetLogger`, `GetRequestID` y helpers internos para inyección de logger en gin.Context y context.Context
- **Constants**: Context keys (`ContextKeyUserID`, `ContextKeyEmail`, `ContextKeyRole`, `ContextKeyClaims`, `ContextKeySlogLogger`, `ContextKeyRequestID`) y headers HTTP (`HeaderRequestID`, `HeaderCorrelationID`)
- **Error Types**: Errores contextuales (`ErrUserIDNotFound`, `ErrEmailNotFound`, `ErrRoleNotFound`, `ErrClaimsNotFound`, `ErrInvalidType`)
- **Resource Path Extraction**: Helper `extractResourceFromPath()` para extraer resource_type (singular) e resource_id del path API REST
- **Singularization**: Helper `singularize()` para convertir nombres de recursos plurales a singulares (roles → role, categories → category)
- **Defensive Defaults**: Paginación con limit default 50, máximo 200; is_active nil = "mostrar todos"; validación de parámetros positivos
- **HTTP Status Logging**: Log levels automáticos según status code (5xx=Error, 4xx=Warn, 2xx/3xx=Info) con duración en milliseconds y bytes transferidos
- **Middleware Chain Recommendation**: Documentación clara del orden correcto: Recovery → RequestLogging → CORS → JWT → PostAuthLogging → Audit → handlers
- **Suite completa de tests unitarios** sin dependencias externas
- **Documentación técnica detallada** en docs/README.md con componentes, flujos comunes, arquitectura y ejemplos de integración
- **Makefile** con targets: fmt, vet, lint, test, build, check

### Design Notes

- **Agnóstico de backend**: Soporta múltiples backends de auditoría mediante interfaz `audit.AuditLogger`
- **Zero-panic por defecto**: Métodos Safe `Get*` retornan errores; variantes `Must*` entran en pánico (uso post-JWT)
- **Contexto enriquecido**: RequestLogging + PostAuthLogging garantizan que todos los logs incluyan request_id, correlation_id, user_id
- **Auditoría post-handler**: AuditMiddleware se ejecuta después de `c.Next()` para capturar status final de la operación
- **Seguridad de concurrencia**: Cada petición obtiene su propio logger; sin estado compartido entre requests
- **Validación defensiva**: ParseListFilters tiene defaults sensatos y capas a valores inválidos
- **Logging de seguridad**: Registra accesos denegados (permission denied) con contexto completo para auditoría

## 0.57.0 - 2026-03-31

### Changed
- fix varios

## 0.56.3 - 2026-03-30

### Changed

- Actualización de dependencia `github.com/EduGoGroup/edugo-shared/repository`

## 0.56.2 - 2026-03-28

### Changed

- Actualización de dependencia `github.com/EduGoGroup/edugo-shared/repository`

## 0.56.1 - 2026-03-28

### Changed

- Actualización de dependencia `github.com/EduGoGroup/edugo-shared/repository`

## 0.56.0 - 2026-03-27

### Changed

- Actualización de dependencia `github.com/EduGoGroup/edugo-shared/repository` de `v0.4.6` a `v0.4.7`.
- Actualización de dependencia indirecta `github.com/EduGoGroup/edugo-infrastructure/postgres` 

## [0.53.0] - 2026-03-26

### Added

- Logging de permission denied en `RequirePermission`, `RequireAnyPermission`, `RequireAllPermissions` con required_permission, missing_permissions, path, method.
- Logging de fallos de autenticacion en `JWTAuthMiddleware`: missing header, invalid format, validation failed con path, ip, method, error.
- Helper `requestPath()` usa `c.FullPath()` para evitar alta cardinalidad en logs/metrics.
- Helper `requestMethod()` con nil guard para robustez en tests.

### Changed

- JWT middleware usa `GetLogger(c)` (context logger) en vez de `slog.Default()` para incluir request_id/correlation_id.
- Permission logs usan constantes `logger.FieldPath`, `logger.FieldMethod` en vez de strings literales.
- `RequireAnyPermission` usa `p.String()` en vez de `string(p)` para consistencia.

## [0.52.0] - 2026-03-26

### Added

- `RequestLogging`: middleware de logging estructurado por request con request_id, correlation_id, method, path, client_ip.
- `PostAuthLogging`: middleware post-JWT que enriquece el logger del contexto con user_id, role y school_id.
- `GetLogger`/`GetRequestID`: helpers para extraer logger y request_id del gin.Context.
- `setLogger`: helper interno para inyectar logger en gin.Context y context.Context.
- Nil guard en `RequestLogging` cuando `baseLogger` es nil (fallback a `slog.Default()`).
- Fallback a `c.Request.URL.Path` cuando `c.FullPath()` retorna vacio (rutas no registradas / 404).
- Log level automatico segun status code: 5xx=Error, 4xx=Warn, 2xx/3xx=Info.

### Changed

- Actualización de dependencia `github.com/EduGoGroup/edugo-shared/logger` de `v0.50.1` a `v0.51.0`.
- Eliminado replace directive local para logger (ahora usa version publicada).

## [0.51.6] - 2026-03-23

### Changed

- Actualización de dependencia `github.com/EduGoGroup/edugo-shared/repository` de `v0.4.5` a `v0.4.6`.
- Actualización de dependencia indirecta `github.com/EduGoGroup/edugo-infrastructure/postgres` de `v0.65.0` a `v0.66.0`.

## [0.51.4] - 2026-03-17

### Added

- Baseline de documentacion de fase 1 con `README.md`, `docs/README.md` y organizacion local por modulo.
