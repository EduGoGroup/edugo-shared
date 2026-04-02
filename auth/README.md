# Auth

Servicios compartidos de autenticación: hashing de password, JWT de acceso, refresh tokens y blacklist.

## Instalación

```bash
go get github.com/EduGoGroup/edugo-shared/auth
```

## Quick Start

### Password Hashing

```go
// Hash password
hashed, err := auth.HashPassword("my-password")

// Verify password
err := auth.VerifyPassword(hashed, "my-password")
```

### JWT Tokens

```go
// Crear JWTManager
manager := auth.NewJWTManager(issuer, secret)

// Generar access token
token, err := manager.GenerateTokenWithContext(ctx, userID, email, activeContext)

// Validar token
claims, err := manager.ValidateToken(token)
```

### Refresh Tokens

```go
// Generar refresh token (retorna token en plaintext y hash para BD)
refreshToken, err := auth.GenerateRefreshToken(7 * 24 * time.Hour)

// Validar refresh token
claims, err := manager.ValidateMinimalToken(refreshToken.Token)
```

### Token Blacklist (Revocación)

```go
// Crear blacklist
blacklist := auth.NewInMemoryBlacklist(ctx)

// Revocar un token
blacklist.Revoke(jti, expiresAt)

// Verificar si está revocado
if blacklist.IsRevoked(jti) {
    // Token revocado
}
```

## Características principales

- **HashPassword**: Bcrypt con costo 12 (~250ms), límite 72 bytes
- **JWTManager**: Generación y validación con contexto RBAC embebido
- **UserContext**: Rol, permisos, escuela y unidad académica
- **RefreshToken**: Tokens criptográficos con hash SHA-256 para almacenamiento seguro
- **TokenBlacklist**: Revocación en memoria con TTL automático

## Documentación

- [Documentación técnica](docs/README.md)
- [Changelog](CHANGELOG.md)

## Operación local

```bash
make build    # Compilar módulo
make test     # Ejecutar tests
make test-race # Tests con race detector
make check    # Validar (fmt, vet, lint, test)
```

## Notas de diseño

- Access tokens requieren `ActiveContext` con rol y permisos.
- Refresh tokens usan flujo mínimal separado para escalabilidad.
- TokenBlacklist es en memoria; adapta a Redis para producción.
- Límite de password de 72 bytes es restricción de bcrypt, no del módulo.
