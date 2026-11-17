# Auth

Servicios compartidos de autenticación: hashing de password, el Context Token ES256 de usuario
(con su contexto activo) y el JWT HS256 entre servicios (M2M).

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

### Context Token ES256 (el token de usuario)

Es el único token de sesión de usuario del ecosistema: ES256 con cabecera `kid` y claims `Claims`
(con `active_context`). El HS256 de usuario, el refresh token propio y la blacklist en memoria se
retiraron en el plan 053 (T6.2); la renovación es un re-canje en edugo-api-identity.
`ValidateToken` rechaza un token HS256, uno sin `kid` y uno con un `kid` que el keyring no conoce.

```go
cfg, err := auth.ContextKeyConfigFromEnv(os.LookupEnv)
// Sin claves configuradas es un error: la API no arranca.
manager, err := auth.NewContextTokenManagerFromConfig(issuer, cfg)

// Solo edugo-api-identity (tiene la clave privada):
token, expiresAt, err := manager.GenerateContextToken(userID, email, activeContext, expiresIn)

// Todas las APIs (y el middleware de Gin):
claims, err := manager.ValidateToken(token)
```

Si ya tienes el keyring: `auth.NewContextTokenManager(issuer, keyring)`.

Variables de entorno (ninguna clave se versiona; vienen de Secret Manager o de un fichero fuera de git):

| Variable | Contenido | Quién la pone |
| --- | --- | --- |
| `AUTH_CONTEXT_SIGNING_KEY` | PEM de la clave privada P-256 activa (PKCS#8 o SEC1) | solo identity |
| `AUTH_CONTEXT_SIGNING_KID` | `kid` de la clave activa | solo identity |
| `AUTH_CONTEXT_PUBLIC_KEYS` | JSON `{"<kid>": "<PEM PKIX>"}` con las públicas aceptadas | todas las APIs (y la anterior al rotar) |

Cada una admite la variante `<NOMBRE>_FILE` con la ruta de un fichero (no las dos a la vez). En un
`.env` usa `AUTH_CONTEXT_SIGNING_KEY_FILE`: el PEM necesita saltos de línea reales y un valor de una
línea con `\n` escapados no se decodifica. Rotar: se publica la nueva pública en
`AUTH_CONTEXT_PUBLIC_KEYS`, se cambia la clave y el `kid` activos en identity, y la pública vieja se
retira cuando caduquen sus tokens.

### Tests: `auth/authtest`

Es el helper único del ecosistema para tokens de usuario en tests (reemplaza a la constante
`TestJWTSecret` de los consumidores). Genera el par ES256 **en memoria** en cada corrida; ninguna
clave se versiona.

```go
import "github.com/EduGoGroup/edugo-shared/auth/authtest"

keys := authtest.Default() // un par por proceso, issuer "edugo-central"

// (a) Acuñar un Context Token válido:
token, err := keys.Token(userID, email, &auth.UserContext{RoleName: "teacher", SchoolID: schoolID}, 15*time.Minute)

// (b) Configurar el verificador como en producción:
t.Setenv(auth.EnvContextPublicKeys, keys.PublicKeysJSON()) // o keys.Env() para un proceso hijo
verifier := keys.Verifier()                                // o *auth.JWTManager listo, solo verifica
```

- `authtest.New(issuer)` / `MustNew(issuer)`: un par nuevo (p. ej. para un token que el servidor no
  debe aceptar).
- `keys.Sign(auth.Claims{...})`: firma claims arbitrarios sin validar (caducado, otro issuer, sin
  `active_context`).
- `keys.Manager()` firma y verifica; `keys.PrivateKeyPEM()` y `keys.KID()` alimentan
  `AUTH_CONTEXT_SIGNING_KEY`/`_KID` de un proceso que deba firmar.
- No importa `testing`: sirve desde `TestMain` y helpers de paquete.

### JWT entre servicios (M2M, HS256)

`ServiceJWTManager` (`NewServiceJWTManager(secret, issuer, audience)`) firma y valida los tokens
servicio→servicio. No es el token de usuario y no cambia con el plan 053.

## Características principales

- **HashPassword**: Bcrypt con costo 12 (~250ms), límite 72 bytes
- **JWTManager**: emisión y verificación del Context Token ES256 con contexto RBAC embebido
- **ContextKeyring**: clave de firma opcional + públicas aceptadas por `kid` (rotación)
- **UserContext**: Rol, grants, escuela, unidad académica y acudido
- **ServiceJWTManager**: JWT M2M HS256

## Documentación

- [Documentación técnica](./docs/README.md)
- [Changelog](./CHANGELOG.md)

## Operación local

```bash
make build    # Compilar módulo
make test     # Ejecutar tests
make test-race # Tests con race detector
make check    # Validar (fmt, vet, lint, test)
```

## Notas de diseño

- El Context Token requiere `ActiveContext` con rol y grants.
- No hay revocación local: la vida corta del token y el re-canje en identity la sustituyen.
- Límite de password de 72 bytes es restricción de bcrypt, no del módulo.
