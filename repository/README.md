# Repository

Adaptadores GORM para CRUD y listados seguros sobre entidades de PostgreSQL.

## Instalación

```bash
go get github.com/EduGoGroup/edugo-shared/repository
```

El módulo se versionan y consume de forma independiente gracias a su `go.mod` propio.

## Quick Start

### Crear repositorio de usuarios con GORM

```go
import (
    "context"
    "github.com/EduGoGroup/edugo-shared/repository"
    "gorm.io/gorm"
)

// Inicializar repositorio
userRepo := repository.NewPostgresUserRepository(db)

// Crear usuario
err := userRepo.Create(ctx, &User{
    Email:    "user@example.com",
    Name:     "John Doe",
})
```

### Búsqueda segura con ListFilters

```go
// Crear filtros con búsqueda segura
filters := &repository.ListFilters{
    Search:       "john",           // Búsqueda segura (previene SQL injection)
    SearchFields: []string{"name", "email"}, // Campos permitidos
    Limit:        10,
    Offset:       0,
}

// Listar usuarios con filtros aplicados
users, total, err := userRepo.List(ctx, filters)
// total = número total de registros (antes de límite)
// users = resultados paginados
```

### Paginación y consultas filtradas

```go
// Aplicar múltiples filtros
filters := &repository.ListFilters{
    Search:       "active",
    SearchFields: []string{"status"},
    FieldFilters: map[string]interface{}{
        "school_id": "sch-123",
        "is_active": true,
    },
    Limit:  20,
    Offset: 40, // Página 3
}

users, total, err := userRepo.List(ctx, filters)
// Iterar resultados paginados
for _, user := range users {
    log.Printf("Usuario: %s", user.Email)
}
```

### Operaciones CRUD básicas

```go
// FindByID - Obtener usuario por ID
user, err := userRepo.FindByID(ctx, "user-123")

// FindByEmail - Obtener usuario por email
user, err := userRepo.FindByEmail(ctx, "user@example.com")

// Update - Actualizar usuario existente
err := userRepo.Update(ctx, &User{
    ID:    "user-123",
    Name:  "Jane Doe",
})

// Delete - Eliminar usuario
err := userRepo.Delete(ctx, "user-123")

// ExistsByEmail - Verificar si existe
exists, err := userRepo.ExistsByEmail(ctx, "user@example.com")
```

## Componentes principales

- **ListFilters**: Estructura para filtros seguros con búsqueda y paginación
- **UserRepository**: Interfaz para operaciones CRUD sobre usuarios
- **SchoolRepository**: Interfaz para operaciones CRUD sobre escuelas
- **MembershipRepository**: Interfaz para operaciones CRUD sobre membresías
- **MembershipAdminRepository**: Extensión con consultas de administración
- **AppError**: Errores tipados (ErrNotFound, etc.)

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

- **Seguridad en búsqueda**: ListFilters valida nombres de campos y escapa patrones SQL
- **Sin políticas de negocio**: Módulo proporciona CRUD genérico, no lógica específica
- **Múltiples repositorios**: UserRepository, SchoolRepository, MembershipRepository intercambiables
- **Context-aware**: Todas las operaciones requieren context.Context para trazabilidad
