# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/auth` se registran aquí.

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
