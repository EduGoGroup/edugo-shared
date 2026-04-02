# Repository — Documentación técnica

Adaptadores GORM para CRUD y listados seguros sobre entidades de PostgreSQL.

## Propósito

Proporcionar capa abstracta de acceso a datos (data access layer) que:
- Define contratos claros (interfaces) para operaciones CRUD
- Implementa búsqueda segura previniendo SQL injection
- Maneja paginación y filtrado de forma consistente
- Proporciona adaptadores concretos para diferentes entidades

## Componentes principales

### ListFilters — Filtros seguros

Estructura que especifica criterios de búsqueda, paginación y filtrado con protección contra SQL injection.

**Estructura:**
```go
type ListFilters struct {
    Search       string                 // Término de búsqueda (aplicado con ILIKE)
    SearchFields []string              // Campos donde buscar (validados contra whitelist)
    FieldFilters map[string]interface{} // Filtros por campo exacto
    Limit        int                    // Registros por página
    Offset       int                    // Desplazamiento (para paginación)
}
```

**Métodos:**
- `ApplySearch(query *gorm.DB) *gorm.DB` — Aplicar búsqueda segura con ILIKE y escaping de patrones
- `ApplyPagination(query *gorm.DB) *gorm.DB` — Aplicar LIMIT y OFFSET
- `GetOffset() int` — Obtener desplazamiento (útil para cálculo de página)

**Características:**
- Validación de nombres de campo usando regex `^[a-zA-Z_][a-zA-Z0-9_]*$`
- Escaping de caracteres especiales en patrones de búsqueda (%, _)
- Protección contra field names inyectados SQL
- Soporte para búsqueda case-insensitive con ILIKE

### UserRepository — Interfaz de usuarios

Define operaciones CRUD sobre la entidad User.

**Métodos:**
- `Create(ctx context.Context, user *User) error` — Crear nuevo usuario
- `FindByID(ctx context.Context, id string) (*User, error)` — Obtener usuario por ID
- `FindByEmail(ctx context.Context, email string) (*User, error)` — Obtener por email
- `ExistsByEmail(ctx context.Context, email string) (bool, error)` — Verificar existencia
- `Update(ctx context.Context, user *User) error` — Actualizar usuario
- `Delete(ctx context.Context, id string) error` — Eliminar usuario
- `List(ctx context.Context, filters *ListFilters) ([]*User, int64, error)` — Listar con filtros

**Características:**
- Operaciones context-aware para trazabilidad
- Búsqueda y paginación integradas
- Manejo de errores tipado (ErrNotFound)

### SchoolRepository — Interfaz de escuelas

Define operaciones CRUD sobre la entidad School.

**Métodos:**
- `Create(ctx context.Context, school *School) error`
- `FindByID(ctx context.Context, id string) (*School, error)`
- `Update(ctx context.Context, school *School) error`
- `Delete(ctx context.Context, id string) error`
- `List(ctx context.Context, filters *ListFilters) ([]*School, int64, error)`

### MembershipRepository — Interfaz de membresías

Define operaciones CRUD sobre la entidad Membership (relación usuario-escuela).

**Métodos:**
- `Create(ctx context.Context, membership *Membership) error`
- `FindByID(ctx context.Context, id string) (*Membership, error)`
- `Update(ctx context.Context, membership *Membership) error`
- `Delete(ctx context.Context, id string) error`
- `List(ctx context.Context, filters *ListFilters) ([]*Membership, int64, error)`

### MembershipAdminRepository — Interfaz extendida de membresías

Extiende MembershipRepository con operaciones administrativas.

**Métodos adicionales:**
- `FindBySchool(ctx context.Context, schoolID string, filters *ListFilters) ([]*Membership, int64, error)` — Listar membresías de una escuela

## Flujos comunes

### 1. Crear repositorio y ejecutar CRUD básico

```go
import (
    "context"
    "github.com/EduGoGroup/edugo-shared/repository"
    "gorm.io/gorm"
)

func initializeUser(ctx context.Context, db *gorm.DB) error {
    userRepo := repository.NewPostgresUserRepository(db)

    // Crear usuario
    user := &User{
        ID:    "user-123",
        Email: "john@example.com",
        Name:  "John Doe",
    }
    if err := userRepo.Create(ctx, user); err != nil {
        return err
    }

    // Recuperar usuario
    retrieved, err := userRepo.FindByID(ctx, "user-123")
    if err != nil {
        return err
    }

    // Actualizar
    retrieved.Name = "John Smith"
    if err := userRepo.Update(ctx, retrieved); err != nil {
        return err
    }

    return nil
}
```

### 2. Búsqueda segura con validación de campos

```go
func searchUsersBy(ctx context.Context, userRepo repository.UserRepository) ([]*User, error) {
    // Filtros con búsqueda segura
    filters := &repository.ListFilters{
        Search:       "john",
        SearchFields: []string{"name", "email"}, // Solo estos campos se buscan
        Limit:        10,
        Offset:       0,
    }

    // ListFilters valida:
    // - Que SearchFields contenga solo nombres válidos
    // - Que Search sea escapado correctamente para prevenir ILIKE injection
    users, total, err := userRepo.List(ctx, filters)
    if err != nil {
        return nil, err
    }

    log.Printf("Encontrados %d usuarios de %d totales", len(users), total)
    return users, nil
}
```

### 3. Paginación con offset/limit

```go
func listUsersWithPagination(ctx context.Context, userRepo repository.UserRepository, pageNum, pageSize int) ([]*User, int64, error) {
    offset := (pageNum - 1) * pageSize

    filters := &repository.ListFilters{
        Limit:  pageSize,
        Offset: offset,
    }

    users, total, err := userRepo.List(ctx, filters)
    if err != nil {
        return nil, 0, err
    }

    totalPages := (total + int64(pageSize) - 1) / int64(pageSize)
    log.Printf("Página %d de %d (total: %d registros)", pageNum, totalPages, total)

    return users, total, nil
}
```

### 4. Operaciones relacionales: membresías por escuela

```go
func listMembershipsForSchool(ctx context.Context, membershipAdminRepo repository.MembershipAdminRepository, schoolID string) ([]*Membership, error) {
    filters := &repository.ListFilters{
        FieldFilters: map[string]interface{}{
            "is_active": true,
        },
        Limit: 50,
    }

    // MembershipAdminRepository proporciona FindBySchool
    memberships, total, err := membershipAdminRepo.FindBySchool(ctx, schoolID, filters)
    if err != nil {
        return nil, err
    }

    log.Printf("Escuela %s tiene %d miembros activos de %d totales", schoolID, len(memberships), total)
    return memberships, nil
}
```

### 5. Manejo de errores tipado

```go
func findUserOrHandle(ctx context.Context, userRepo repository.UserRepository, userID string) (*User, error) {
    user, err := userRepo.FindByID(ctx, userID)

    // Manejo de error ErrNotFound
    if err != nil {
        if errors.Is(err, repository.ErrNotFound) {
            log.Printf("Usuario %s no existe", userID)
            return nil, fmt.Errorf("usuario no encontrado")
        }
        return nil, fmt.Errorf("error al buscar usuario: %w", err)
    }

    return user, nil
}
```

## Arquitectura

Flujo de operación típico:

```
1. Crear Manager de ciclo de vida (lifecycle)
   ↓
2. Inicializar GORM DB (gorm.io/gorm)
   ↓
3. Crear adaptadores concretos (NewPostgresUserRepository)
   ├─ UserRepository
   ├─ SchoolRepository
   ├─ MembershipRepository
   └─ MembershipAdminRepository
   ↓
4. Usar repositorios en lógica de negocio
   ├─ CRUD básico (Create, FindByID, Update, Delete)
   ├─ Búsqueda con ListFilters (segura contra SQL injection)
   ├─ Paginación (Limit, Offset)
   └─ Consultas especializadas (FindByEmail, FindBySchool)
   ↓
5. Propagar errores tipados (ErrNotFound)
```

Protecciones en ListFilters:

```
Input: {"search": "'; DROP TABLE--", "searchFields": ["name"]}
              ↓
         Validar SearchFields vs whitelist
              ↓
         Escapar caracteres SQL en search (%, _)
              ↓
         Generar: db.Where("name ILIKE ?", "%'; DROP TABLE--%")
              ↓
         Output: búsqueda segura, sin inyección
```

## Dependencias

- **Internas**: Ninguna (módulo autónomo)
- **Externas**:
  - `gorm.io/gorm` (GORM ORM)
  - `github.com/google/uuid` (generación de UUID)
  - `github.com/EduGoGroup/edugo-infrastructure/postgres/entities` (definiciones de entidades)

## Testing

Suite de tests completa:

- Creación de repositorios (NewPostgresUserRepository, etc.)
- Operaciones CRUD (Create, FindByID, Update, Delete)
- Búsqueda segura con ListFilters (validación de campos, escaping)
- Paginación (Limit, Offset, cálculo de total)
- Errores tipados (ErrNotFound)
- Operaciones relacionales (FindBySchool)
- SQL injection prevention (caracteres maliciosos en campos)

Ejecutar:
```bash
make test          # Tests básicos
make test-race     # Tests con race detector
make check         # Tests + linting + format
```

## Notas de diseño

- **Sin lógica de negocio**: Repositorios proporcionan acceso genérico, no reglas específicas del dominio
- **Seguridad en primer lugar**: Búsqueda con ILIKE requiere validación de campos y escaping de patrones
- **Múltiples adaptadores**: Cada entidad (User, School, Membership) tiene su propio repositorio
- **Context-aware**: Todas las operaciones requieren context.Context para trazabilidad
- **Errores tipados**: ErrNotFound permite manejar ausencia de registros explícitamente
- **Paginación integrada**: ListFilters abstrae Limit/Offset
- **GORM agnóstico**: Adaptadores usan *gorm.DB sin exponerlo en interfaces públicas
