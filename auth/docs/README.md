# Auth — Documentación técnica

Servicios compartidos de autenticación: hashing de password, el Context Token ES256 de usuario y el JWT M2M entre servicios. El HS256 de usuario, el refresh token propio y la blacklist en memoria se retiraron en el plan 053 (T6.2).

## Propósito

Proporcionar primitivos de autenticación y autorización reutilizables para los servicios de EduGo.

## Componentes principales

### password.go — Hashing seguro

**Constantes**
- `bcryptCost = 12` → ~250ms por hash (balance seguridad/UX)
- `maxPasswordLength = 72` → límite de bcrypt

**HashPassword(password string) (string, error)**
Genera hash bcrypt con salt aleatorio. Valida límite de 72 bytes.

**VerifyPassword(hash, password string) error**
Verifica password contra su hash usando bcrypt.Compare.

### jwt_claims.go — Tipos de datos para JWT

**UserContext**
Contexto activo del usuario embebido en JWT:
- `RoleID`, `RoleName` — Rol actual
- `SchoolID`, `SchoolName` — Escuela (opcional)
- `AcademicUnitID`, `AcademicUnitName` — Unidad académica (opcional)
- `Permissions` — Lista de permisos

**Claims**
Claims del Context Token, que heredan `jwt.RegisteredClaims`:
- `UserID`, `Email` — Identidad del usuario
- `ActiveContext` — UserContext embebido (obligatorio)

### jwt_manager.go y context_token.go — Context Token ES256

**JWTManager**
Encapsula issuer y keyring ES256. Solo acepta ES256 con un `kid` del keyring.

Métodos principales:
- `NewContextTokenManager(issuer, keyring) (*JWTManager, error)`
- `NewContextTokenManagerFromConfig(issuer, cfg) (*JWTManager, error)`
- `GenerateContextToken(userID, email, activeContext, expiresIn) (string, time.Time, error)` — requiere clave de firma
- `GenerateContextTokenUntil(userID, email, activeContext, expiresAt) (string, time.Time, error)` — para el canje
- `ValidateToken(token) (*Claims, error)` — exige ES256, `kid` conocido, issuer, `exp` y ActiveContext

### context_keys_config.go — Carga de claves

`ContextKeyConfigFromEnv`, `NewContextKeyringFromConfig` y los parsers PEM. Variables
`AUTH_CONTEXT_SIGNING_KEY`, `AUTH_CONTEXT_SIGNING_KID`, `AUTH_CONTEXT_PUBLIC_KEYS` (y sus `_FILE`).

### jwt_extract.go — Helpers para extracción

**ExtractUserID(token string) (string, error)**
Extrae userID de un token sin validar completamente. Útil solo para logging/debugging, NO para autenticación.

### authtest — Helper de tests

Par ES256 en memoria para acuñar tokens y configurar el verificador en tests. Uso en el
[README del módulo](../README.md).

### service_claims.go — JWT M2M

`ServiceJWTManager`, HS256 entre servicios (`SERVICE_JWT_SECRET` y afines). No es el token de usuario.

## Flujos comunes

### 1. Registro de usuario

```go
// HashPassword genera hash seguro para BD
hash, err := auth.HashPassword(userPassword)
// Guardar hash en BD
```

### 2. Login del usuario

```go
// Recuperar hash de BD
storedHash := getUserHash(userID)
// Verificar password
err := auth.VerifyPassword(storedHash, providedPassword)
if err == nil {
    // Login exitoso, generar tokens
}
```

### 3. Emitir un Context Token (solo identity)

```go
manager, err := auth.NewContextTokenManagerFromConfig(issuer, cfg) // cfg con clave de firma
userCtx := &auth.UserContext{RoleID: "admin", RoleName: "Administrator", SchoolID: "sch-123"}
token, expiresAt, err := manager.GenerateContextToken(userID, email, userCtx, 15*time.Minute)
```

### 4. Validar el token (cualquier API)

```go
claims, err := manager.ValidateToken(token)
if err != nil {
    return errors.Unauthorized
}
// Usar claims.UserID, claims.ActiveContext.Grants, etc.
```

La renovación no pasa por este módulo: el cliente re-canjea en edugo-api-identity.

## Testing

El módulo incluye tests completos y benchmarks:
- Password hashing/verification
- Token generation/validation
- Alg confusion, `kid` y rotación de claves
- Benchmarks de bcrypt y JWT

Ejecutar:
```bash
make test          # Tests básicos
make test-race     # Tests con race detector
make check         # Tests + linting + format
```

## Notas de diseño

- **Context Token**: requiere ActiveContext con rol y grants. Ideal para autorización en handlers.
- **Sin revocación local**: la vida corta del token y el re-canje sustituyen a la blacklist.
- **Password límite**: 72 bytes es restricción de bcrypt, no del módulo. Usar validación en cliente.
