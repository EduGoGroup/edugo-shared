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
manager := auth.NewJWTManager(secretKey, issuer)

// Generar access token
token, expiresAt, err := manager.GenerateTokenWithContext(userID, email, activeContext, expiresIn)

// Validar token
claims, err := manager.ValidateToken(token)
```

### Context Token ES256 (plan 053, convive con HS256)

Mismos claims que el access HS256 (`Claims`, con `active_context`), firmados con ES256 y cabecera
`kid`. HS256 sigue funcionando: sin keyring el `JWTManager` se comporta exactamente como antes.

```go
cfg, err := auth.ContextKeyConfigFromEnv(os.LookupEnv)
manager := auth.NewJWTManager(secretKey, issuer)
if !cfg.IsZero() {
    keyring, err := auth.NewContextKeyringFromConfig(cfg) // mal configurado = error de arranque
    manager.WithContextKeyring(keyring)
}

// Solo edugo-api-identity (tiene la clave privada):
token, expiresAt, err := manager.GenerateContextToken(userID, email, activeContext, expiresIn)

// Middleware y APIs: sin cambios, ValidateToken acepta HS256 (sin kid) y ES256 (kid conocido).
claims, err := manager.ValidateToken(token)
```

Variables de entorno (ninguna clave se versiona; vienen de Secret Manager o de un fichero fuera de git):

| Variable | Contenido | Quién la pone |
| --- | --- | --- |
| `AUTH_CONTEXT_SIGNING_KEY` | PEM de la clave privada P-256 activa (PKCS#8 o SEC1) | solo identity |
| `AUTH_CONTEXT_SIGNING_KID` | `kid` de la clave activa | solo identity |
| `AUTH_CONTEXT_PUBLIC_KEYS` | JSON `{"<kid>": "<PEM PKIX>"}` con las públicas aceptadas | todas las APIs (y la anterior al rotar) |

Cada una admite la variante `<NOMBRE>_FILE` con la ruta de un fichero (no las dos a la vez). En un
`.env` usa `AUTH_CONTEXT_SIGNING_KEY_FILE`: el PEM necesita saltos de línea reales y un valor de una
línea con `\n` escapados no se decodifica. Sin
ninguna, la API queda solo en HS256. Rotar: se publica la nueva pública en `AUTH_CONTEXT_PUBLIC_KEYS`,
se cambia la clave y el `kid` activos en identity, y la pública vieja se retira cuando caduquen sus
tokens.

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
