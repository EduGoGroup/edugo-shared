# Auth — Documentación técnica

Servicios compartidos de autenticación: hashing de password, JWT de acceso, refresh tokens y blacklist de revocación.

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
Claims personalizados que heredan `jwt.RegisteredClaims`:
- `UserID`, `Email` — Identidad del usuario
- `ActiveContext` — UserContext embebido (obligatorio para access tokens)
- `TokenUse` — Tipo de token: access (vacío), refresh, etc.
- `SchoolID` — Se preserva en refresh tokens

### jwt_manager.go — Generación y validación de tokens

**JWTManager**
Encapsula issuer y secret para generación/validación.

Métodos principales:
- `NewJWTManager(secretKey, issuer) *JWTManager`
- `GenerateTokenWithContext(userID, email, activeContext, expiresIn) (string, time.Time, error)` — Access token con contexto
- `ValidateToken(token) (*Claims, error)` — Valida access token, requiere ActiveContext
- `GenerateMinimalToken(userID, email, schoolID, expiresIn) (string, time.Time, error)` — Token sin contexto (para refresh)
- `ValidateMinimalToken(token) (*Claims, error)` — Valida refresh token

### jwt_extract.go — Helpers para extracción

**ExtractUserID(token string) (string, error)**
Extrae userID de un token sin validar completamente. Útil solo para logging/debugging, NO para autenticación.

### refresh_token.go — Tokens criptográficos

**RefreshToken**
Representa un refresh token:
- `Token` — Texto plano (retorna al cliente)
- `TokenHash` — SHA-256 en hex (guarda en BD)
- `ExpiresAt` — Timestamp de expiración

**GenerateRefreshToken(ttl time.Duration) (*RefreshToken, error)**
Genera 32 bytes aleatorios con crypto/rand, codifica en base64 URL-safe, calcula SHA-256.

**HashToken(token string) string**
Calcula SHA-256 hex del token.

**VerifyTokenHash(token, hash string) bool**
Compara token contra su hash almacenado.

### blacklist.go — Revocación de tokens

**TokenBlacklist (interfaz)**
Contrato para revocación:
- `Revoke(jti string, expiresAt time.Time)` — Agregar a blacklist
- `IsRevoked(jti string) bool` — Verificar si está revocado

**InMemoryBlacklist**
Implementación en memoria usando sync.Map con TTL:
- `NewInMemoryBlacklist(ctx context.Context) *InMemoryBlacklist` — Constructor con cleanup goroutine
- Cleanup automático de entradas expiradas
- Seguro para concurrencia

Nota: Para producción con escalabilidad, reemplazar con Redis.

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

### 3. Generar access y refresh tokens

```go
manager := auth.NewJWTManager(issuer, secret)

// Access token con contexto
userCtx := &auth.UserContext{
    RoleID:      "admin",
    RoleName:    "Administrator",
    SchoolID:    "sch-123",
    Permissions: []string{"users:read", "users:write"},
}
accessToken, _ := manager.GenerateTokenWithContext(ctx, userID, email, userCtx)

// Refresh token minimal
refreshToken, _ := auth.GenerateRefreshToken(7 * 24 * time.Hour)
// Guardar refreshToken.TokenHash en BD
// Retornar accessToken y refreshToken.Token al cliente
```

### 4. Validar access token

```go
claims, err := manager.ValidateToken(accessToken)
if err != nil {
    return errors.Unauthorized
}
if claims.ActiveContext == nil {
    return errors.InvalidToken
}
// Usar claims.UserID, claims.ActiveContext.Permissions, etc.
```

### 5. Refresh token (obtener nuevo access token)

```go
// Cliente envía refresh token
claims, err := manager.ValidateMinimalToken(refreshToken)
if err != nil {
    return errors.Unauthorized
}

// Recuperar hash de BD y verificar
storedHash := getRefreshTokenHash(claims.UserID)
if !auth.VerifyTokenHash(refreshToken, storedHash) {
    return errors.Unauthorized
}

// Generar nuevo access token
newAccessToken, _ := manager.GenerateTokenWithContext(ctx, claims.UserID, claims.Email, newUserContext)
```

### 6. Revocar token (logout)

```go
blacklist := auth.NewInMemoryBlacklist(ctx)
// En middleware o durante logout
blacklist.Revoke(claims.ID, claims.ExpiresAt)

// En validación, verificar antes de usar
if blacklist.IsRevoked(claims.ID) {
    return errors.TokenRevoked
}
```

## Testing

El módulo incluye tests completos y benchmarks:
- Password hashing/verification
- Token generation/validation
- Refresh token flows
- Blacklist concurrency (race detector)
- Benchmarks de bcrypt y JWT

Ejecutar:
```bash
make test          # Tests básicos
make test-race     # Tests con race detector
make check         # Tests + linting + format
```

## Notas de diseño

- **Access tokens**: Requieren ActiveContext con rol y permisos. Ideales para autorización en handlers.
- **Refresh tokens**: Flujo mínimal separado, stateless. Ideales para renovación.
- **Password límite**: 72 bytes es restricción de bcrypt, no del módulo. Usar validación en cliente.
- **Blacklist**: En memoria para baja latencia. Redis recomendado para múltiples instancias.
- **Claims personalizados**: UserContext permite RBAC flexible sin JTI externo.
