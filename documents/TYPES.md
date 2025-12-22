# Tipos y Enumeraciones

## Visión General

El módulo `common/types` proporciona tipos compartidos y enumeraciones del dominio que se utilizan consistentemente en todo el ecosistema EduGo.

---

## Ubicación de Tipos

```
common/
├── types/
│   ├── enum/
│   │   ├── role.go        # Roles del sistema
│   │   ├── status.go      # Estados de materiales y progreso
│   │   ├── assessment.go  # Tipos de evaluación
│   │   └── event.go       # Tipos de eventos
│   ├── uuid.go            # Helpers para UUIDs
│   └── uuid_test.go
└── errors/
    └── errors.go          # Códigos de error
```

---

## 1. Roles del Sistema (SystemRole)

**Archivo:** `common/types/enum/role.go`  
**Import:** `github.com/EduGoGroup/edugo-shared/common/types/enum`

### Definición

```go
type SystemRole string

const (
    // Administrador del sistema - acceso total
    SystemRoleAdmin SystemRole = "admin"
    
    // Profesor - puede crear materiales y evaluaciones
    SystemRoleTeacher SystemRole = "teacher"
    
    // Estudiante - consume contenido y realiza evaluaciones
    SystemRoleStudent SystemRole = "student"
    
    // Tutor/Padre - puede ver progreso de estudiantes
    SystemRoleGuardian SystemRole = "guardian"
)
```

### Métodos

```go
// IsValid verifica si el rol es válido
func (r SystemRole) IsValid() bool

// String retorna la representación string
func (r SystemRole) String() string

// AllSystemRoles retorna todos los roles válidos
func AllSystemRoles() []SystemRole

// AllSystemRolesStrings retorna roles como []string
func AllSystemRolesStrings() []string
```

### Uso

```go
import "github.com/EduGoGroup/edugo-shared/common/types/enum"

// Asignar rol
user.Role = enum.SystemRoleStudent

// Validar rol
if !user.Role.IsValid() {
    return errors.NewValidationError("invalid role")
}

// Comparar rol
if user.Role == enum.SystemRoleAdmin {
    // Acceso administrativo
}

// Obtener todos los roles para validación
validRoles := enum.AllSystemRolesStrings()
// ["admin", "teacher", "student", "guardian"]
```

### Tabla de Permisos Conceptual

| Rol | Crear Material | Ver Material | Crear Evaluación | Ver Progreso |
|-----|----------------|--------------|------------------|--------------|
| Admin | ✅ | ✅ | ✅ | ✅ |
| Teacher | ✅ | ✅ | ✅ | ✅ (sus estudiantes) |
| Student | ❌ | ✅ | ❌ | ✅ (propio) |
| Guardian | ❌ | ❌ | ❌ | ✅ (sus hijos) |

---

## 2. Estados de Material (MaterialStatus)

**Archivo:** `common/types/enum/status.go`

### Definición

```go
type MaterialStatus string

const (
    // Borrador - no visible para estudiantes
    MaterialStatusDraft MaterialStatus = "draft"
    
    // Publicado - visible y accesible
    MaterialStatusPublished MaterialStatus = "published"
    
    // Archivado - histórico, no accesible
    MaterialStatusArchived MaterialStatus = "archived"
)
```

### Métodos

```go
func (s MaterialStatus) IsValid() bool
func (s MaterialStatus) String() string
func AllMaterialStatuses() []MaterialStatus
```

### Diagrama de Estados

```
┌─────────┐                  ┌────────────┐
│  DRAFT  │ ──── publish ──► │ PUBLISHED  │
└─────────┘                  └─────┬──────┘
     ▲                             │
     │                         archive
  unpublish                        │
     │                             ▼
     └─────────────────────  ┌────────────┐
                             │  ARCHIVED  │
                             └────────────┘
```

### Uso

```go
// Crear material en borrador
material := &Material{
    Status: enum.MaterialStatusDraft,
}

// Publicar
func (m *Material) Publish() error {
    if m.Status != enum.MaterialStatusDraft {
        return errors.NewBusinessRuleError("can only publish draft materials")
    }
    m.Status = enum.MaterialStatusPublished
    return nil
}

// Filtrar materiales publicados
query := db.Where("status = ?", enum.MaterialStatusPublished)
```

---

## 3. Estados de Progreso (ProgressStatus)

**Archivo:** `common/types/enum/status.go`

### Definición

```go
type ProgressStatus string

const (
    // No iniciado - usuario no ha comenzado
    ProgressStatusNotStarted ProgressStatus = "not_started"
    
    // En progreso - usuario está consumiendo contenido
    ProgressStatusInProgress ProgressStatus = "in_progress"
    
    // Completado - usuario finalizó el contenido
    ProgressStatusCompleted ProgressStatus = "completed"
)
```

### Métodos

```go
func (p ProgressStatus) IsValid() bool
func (p ProgressStatus) String() string
func AllProgressStatuses() []ProgressStatus
```

### Diagrama de Estados

```
┌──────────────┐              ┌──────────────┐              ┌───────────┐
│ NOT_STARTED  │ ── start ──► │ IN_PROGRESS  │ ── finish ─► │ COMPLETED │
└──────────────┘              └──────────────┘              └───────────┘
                                    │                            │
                                    └────────── reset ───────────┘
                                                │
                                                ▼
                                         ┌──────────────┐
                                         │ NOT_STARTED  │
                                         └──────────────┘
```

### Uso

```go
// Iniciar progreso
func (p *Progress) Start() error {
    if p.Status != enum.ProgressStatusNotStarted {
        return errors.NewBusinessRuleError("already started")
    }
    p.Status = enum.ProgressStatusInProgress
    p.StartedAt = time.Now()
    return nil
}

// Completar
func (p *Progress) Complete() error {
    if p.Status != enum.ProgressStatusInProgress {
        return errors.NewBusinessRuleError("must be in progress to complete")
    }
    p.Status = enum.ProgressStatusCompleted
    p.CompletedAt = time.Now()
    return nil
}
```

---

## 4. Estados de Procesamiento (ProcessingStatus)

**Archivo:** `common/types/enum/status.go`

### Definición

```go
type ProcessingStatus string

const (
    // Pendiente - en cola para procesar
    ProcessingStatusPending ProcessingStatus = "pending"
    
    // Procesando - siendo procesado activamente
    ProcessingStatusProcessing ProcessingStatus = "processing"
    
    // Completado - procesamiento exitoso
    ProcessingStatusCompleted ProcessingStatus = "completed"
    
    // Fallido - procesamiento con error
    ProcessingStatusFailed ProcessingStatus = "failed"
)
```

### Métodos

```go
func (p ProcessingStatus) IsValid() bool
func (p ProcessingStatus) String() string
func AllProcessingStatuses() []ProcessingStatus
```

### Diagrama de Estados

```
┌──────────┐              ┌─────────────┐
│ PENDING  │ ── pick ───► │ PROCESSING  │
└──────────┘              └──────┬──────┘
     ▲                           │
     │               ┌───────────┴───────────┐
   retry             │                       │
     │          success                   failure
     │               │                       │
     │               ▼                       ▼
     │        ┌───────────┐           ┌──────────┐
     └─────── │ COMPLETED │           │  FAILED  │
              └───────────┘           └──────────┘
```

### Uso

```go
// Proceso de upload de archivo
type FileUpload struct {
    ID              string
    Status          enum.ProcessingStatus
    ErrorMessage    string
    RetryCount      int
}

// Worker de procesamiento
func (w *Worker) Process(upload *FileUpload) {
    upload.Status = enum.ProcessingStatusProcessing
    db.Save(upload)
    
    err := w.processFile(upload)
    if err != nil {
        upload.Status = enum.ProcessingStatusFailed
        upload.ErrorMessage = err.Error()
        upload.RetryCount++
    } else {
        upload.Status = enum.ProcessingStatusCompleted
    }
    db.Save(upload)
}
```

---

## 5. Tipos de UUID

**Archivo:** `common/types/uuid.go`

### Funciones

```go
// NewUUID genera un nuevo UUID v4
func NewUUID() string

// IsValidUUID verifica si un string es un UUID válido
func IsValidUUID(s string) bool

// MustParseUUID parsea un UUID, panic si es inválido
func MustParseUUID(s string) uuid.UUID

// ParseUUID parsea un UUID con manejo de error
func ParseUUID(s string) (uuid.UUID, error)
```

### Uso

```go
import "github.com/EduGoGroup/edugo-shared/common/types"

// Generar nuevo ID
userID := types.NewUUID()  // "550e8400-e29b-41d4-a716-446655440000"

// Validar ID de request
if !types.IsValidUUID(req.UserID) {
    return errors.NewValidationError("invalid user_id format")
}

// Parsear para operaciones
id, err := types.ParseUUID(req.UserID)
if err != nil {
    return errors.NewValidationError("invalid uuid")
}
```

---

## 6. Códigos de Error (ErrorCode)

**Archivo:** `common/errors/errors.go`

### Definición

```go
type ErrorCode string

const (
    // Errores de Validación (4xx)
    ErrorCodeValidation    ErrorCode = "VALIDATION_ERROR"
    ErrorCodeInvalidInput  ErrorCode = "INVALID_INPUT"
    
    // Errores de Recurso (4xx)
    ErrorCodeNotFound      ErrorCode = "NOT_FOUND"
    ErrorCodeAlreadyExists ErrorCode = "ALREADY_EXISTS"
    ErrorCodeConflict      ErrorCode = "CONFLICT"
    
    // Errores de Autenticación (4xx)
    ErrorCodeUnauthorized  ErrorCode = "UNAUTHORIZED"
    ErrorCodeForbidden     ErrorCode = "FORBIDDEN"
    ErrorCodeInvalidToken  ErrorCode = "INVALID_TOKEN"
    ErrorCodeTokenExpired  ErrorCode = "TOKEN_EXPIRED"
    
    // Errores de Negocio (4xx)
    ErrorCodeBusinessRule  ErrorCode = "BUSINESS_RULE_VIOLATION"
    ErrorCodeInvalidState  ErrorCode = "INVALID_STATE"
    
    // Errores de Servidor (5xx)
    ErrorCodeInternal        ErrorCode = "INTERNAL_ERROR"
    ErrorCodeDatabaseError   ErrorCode = "DATABASE_ERROR"
    ErrorCodeExternalService ErrorCode = "EXTERNAL_SERVICE_ERROR"
    ErrorCodeTimeout         ErrorCode = "TIMEOUT"
    
    // Errores de Límites
    ErrorCodeRateLimit     ErrorCode = "RATE_LIMIT_EXCEEDED"
    ErrorCodeQuotaExceeded ErrorCode = "QUOTA_EXCEEDED"
)
```

### Mapeo a HTTP Status

| ErrorCode | HTTP Status | Descripción |
|-----------|-------------|-------------|
| `VALIDATION_ERROR` | 400 Bad Request | Datos de entrada inválidos |
| `INVALID_INPUT` | 400 Bad Request | Formato de request incorrecto |
| `NOT_FOUND` | 404 Not Found | Recurso no existe |
| `ALREADY_EXISTS` | 409 Conflict | Recurso ya existe |
| `CONFLICT` | 409 Conflict | Conflicto de estado |
| `UNAUTHORIZED` | 401 Unauthorized | No autenticado |
| `FORBIDDEN` | 403 Forbidden | Sin permisos |
| `INVALID_TOKEN` | 401 Unauthorized | Token JWT inválido |
| `TOKEN_EXPIRED` | 401 Unauthorized | Token expirado |
| `BUSINESS_RULE_VIOLATION` | 422 Unprocessable Entity | Regla de negocio violada |
| `INVALID_STATE` | 422 Unprocessable Entity | Transición de estado inválida |
| `INTERNAL_ERROR` | 500 Internal Server Error | Error interno |
| `DATABASE_ERROR` | 500 Internal Server Error | Error de BD |
| `EXTERNAL_SERVICE_ERROR` | 500 Internal Server Error | Servicio externo falló |
| `TIMEOUT` | 408 Request Timeout | Timeout |
| `RATE_LIMIT_EXCEEDED` | 429 Too Many Requests | Rate limit |
| `QUOTA_EXCEEDED` | 429 Too Many Requests | Cuota excedida |

---

## 7. Tipos de Assessment (Evaluación)

**Archivo:** `common/types/enum/assessment.go`

### Definición

```go
type AssessmentType string

const (
    // Quiz - evaluación corta con preguntas múltiples
    AssessmentTypeQuiz AssessmentType = "quiz"
    
    // Examen - evaluación formal más extensa
    AssessmentTypeExam AssessmentType = "exam"
    
    // Tarea - trabajo para entregar
    AssessmentTypeAssignment AssessmentType = "assignment"
    
    // Práctica - ejercicio sin calificación formal
    AssessmentTypePractice AssessmentType = "practice"
)
```

---

## 8. Tipos de Evento (EventType)

**Archivo:** `common/types/enum/event.go`

### Definición

```go
type EventType string

const (
    // Eventos de Usuario
    EventTypeUserCreated   EventType = "user.created"
    EventTypeUserUpdated   EventType = "user.updated"
    EventTypeUserDeleted   EventType = "user.deleted"
    
    // Eventos de Material
    EventTypeMaterialCreated   EventType = "material.created"
    EventTypeMaterialPublished EventType = "material.published"
    EventTypeMaterialArchived  EventType = "material.archived"
    
    // Eventos de Progreso
    EventTypeProgressStarted   EventType = "progress.started"
    EventTypeProgressCompleted EventType = "progress.completed"
    
    // Eventos de Assessment
    EventTypeAssessmentSubmitted EventType = "assessment.submitted"
    EventTypeAssessmentGraded    EventType = "assessment.graded"
)
```

### Uso para Event Sourcing / Audit

```go
// Publicar evento
type Event struct {
    Type      enum.EventType `json:"type"`
    UserID    string         `json:"user_id"`
    Resource  string         `json:"resource"`
    Data      interface{}    `json:"data"`
    Timestamp time.Time      `json:"timestamp"`
}

event := Event{
    Type:      enum.EventTypeMaterialPublished,
    UserID:    "user-123",
    Resource:  "material-456",
    Data:      material,
    Timestamp: time.Now(),
}

publisher.Publish(ctx, "events", "material.published", event)
```

---

## Resumen de Imports

```go
// Enums del dominio
import "github.com/EduGoGroup/edugo-shared/common/types/enum"

// Helpers de UUID
import "github.com/EduGoGroup/edugo-shared/common/types"

// Errores estructurados
import "github.com/EduGoGroup/edugo-shared/common/errors"
```

---

## Validación con Tags

Los tipos se pueden usar con validadores:

```go
type CreateUserRequest struct {
    Email string         `json:"email" validate:"required,email"`
    Role  enum.SystemRole `json:"role" validate:"required,oneof=admin teacher student guardian"`
}

// Validar
if !req.Role.IsValid() {
    return errors.NewValidationError("invalid role")
}
```
