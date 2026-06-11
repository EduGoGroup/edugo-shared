# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/auth` se registran aquí.

## [Unreleased]

## [0.900.0] - 2026-06-11

Migración a banda `0.900.x` (ADR 0015 / bug 0022). Misma base que `v0.2.0`.

### Added
- `ServiceClaims`: claims de un service JWT (M2M / B2B) con `token_use="service"`, `client_id`, `scopes` y `aud`; SIN `active_context` ni `user_id`. Método `HasScope(scope)`.
- `ServiceJWTManager` (firma/valida HS256 con su PROPIO secret `SERVICE_JWT_SECRET`, distinto del de usuarios): `NewServiceJWTManager(secret, issuer, audience)`, `GenerateServiceToken`, `ParseServiceToken`, `ValidateServiceToken`.
- Constante `TokenUseService = "service"`.

Habilita la autenticación M2M del plan 020 N5 (D14/D17).

## [0.2.0] - 2026-06-11

### Added
- `ServiceClaims`: claims de un service JWT (M2M / B2B) con `token_use="service"`, `client_id`, `scopes` y `aud`; SIN `active_context` ni `user_id` (el caller ya resolvió destinatarios). Método `HasScope(scope)`.
- `ServiceJWTManager` (firma/valida HS256 con su PROPIO secret `SERVICE_JWT_SECRET`, distinto del de usuarios): `NewServiceJWTManager(secret, issuer, audience)`, `GenerateServiceToken`, `ParseServiceToken` (valida firma + iss + aud + exp) y `ValidateServiceToken` (además exige `token_use == "service"`).
- Constante `TokenUseService = "service"`.

Habilita la autenticación M2M del plan 020 N5 (D14/D17): platform valida el token en `/api/v1/internal/*` y los callers (worker, learning) lo firman. Un JWT de usuario NO valida como service token y viceversa.

## [0.1.0] - 2026-05-28

### Added
- Reinicio de la versión del módulo a `v0.1.0` (borrón y cuenta nueva).
- Conservación del código de producción estable del módulo.

## [0.100.0] - 2026-04-02

### Changed

- Removed trivial `TestNoOpBlacklist` test (covered by middleware/gin/jwt_auth_test.go)

### Added

- Hashing seguro de passwords con bcrypt (costo 12, ~250ms, límite 72 bytes).
- Función `HashPassword(password)` y `VerifyPassword(hash, password)`.
- `JWTManager` para generación y validación de JWT con contexto RBAC embebido.
- `UserContext` con rol, permisos, escuela y unidad académica.
- `Claims` personalizados que heredan `jwt.RegisteredClaims`.
- Métodos JWT: `GenerateTokenWithContext`, `ValidateToken`, `GenerateMinimalToken`, `ValidateMinimalToken`.
- Generación criptográfica de refresh tokens con 32 bytes aleatorios (crypto/rand).
- `RefreshToken` con token plaintext y hash SHA-256 para almacenamiento seguro.
- Funciones: `GenerateRefreshToken`, `HashToken`, `VerifyTokenHash`.
- `TokenBlacklist` interfaz para revocación de tokens.
- `InMemoryBlacklist` con TTL automático y cleanup concurrente.
- Suite completa de tests unitarios con race detector.
- Benchmarks de bcrypt y JWT.
- Documentación en README.md y docs/README.md.
- Makefile con targets: build, test, test-race, check, lint, fmt, vet, tidy, deps, release.

### Design Notes

- Access tokens requieren `ActiveContext` para autorización.
- Refresh tokens usan flujo mínimal separado, stateless.
- Blacklist en memoria; Redis recomendado para producción con múltiples instancias.
- Versión v0.100.0 marca estabilización del contrato de autenticación.
