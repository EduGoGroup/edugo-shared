# Módulos Nuevos a Crear

## 🎯 Objetivo

Especificar detalladamente los módulos que **NO existen** en edugo-shared pero son **necesarios** para los proyectos consumidores.

---

## 📦 evaluation/ (CRÍTICO - P0)

### Justificación

**Requerido por:**
- ✅ api-mobile (sistema de evaluaciones)
- ✅ worker (generación de quizzes con IA)

**Para qué:**
- Modelos compartidos de Assessment, Questions, Attempts
- Evitar duplicación de estructuras entre servicios
- Garantizar consistencia en MongoDB y PostgreSQL

**Bloqueante:** **SÍ** - Sin este módulo api-mobile NO puede implementar evaluaciones

---

### Especificación Técnica

#### go.mod

```go
module github.com/EduGoGroup/edugo-shared/evaluation

go 1.24

require (
    github.com/google/uuid v1.6.0
    github.com/EduGoGroup/edugo-shared/common v0.5.0
)
```

---

#### Estructuras a Exportar

**Archivo:** `assessment.go`

```go
package evaluation

import (
    "time"
    "github.com/google/uuid"
)

// Assessment representa un cuestionario generado o manual
type Assessment struct {
    ID                uuid.UUID `json:"id" bson:"_id"`
    MaterialID        int64     `json:"material_id" bson:"material_id"`
    MongoDocID        string    `json:"mongo_doc_id,omitempty" bson:"mongo_doc_id,omitempty"` // Referencia a documento en MongoDB
    Title             string    `json:"title" bson:"title"`
    Description       string    `json:"description,omitempty" bson:"description,omitempty"`
    Type              string    `json:"type" bson:"type"` // "manual", "generated"
    Status            string    `json:"status" bson:"status"` // "draft", "published", "archived"
    PassingScore      int       `json:"passing_score" bson:"passing_score"` // Porcentaje mínimo para aprobar (0-100)
    TotalQuestions    int       `json:"total_questions" bson:"total_questions"`
    TotalPoints       int       `json:"total_points" bson:"total_points"`
    CreatedBy         int64     `json:"created_by" bson:"created_by"` // User ID
    CreatedAt         time.Time `json:"created_at" bson:"created_at"`
    UpdatedAt         time.Time `json:"updated_at" bson:"updated_at"`
}

// Validate valida los campos del assessment
func (a *Assessment) Validate() error {
    if a.Title == "" {
        return errors.New("title is required")
    }
    if a.PassingScore < 0 || a.PassingScore > 100 {
        return errors.New("passing score must be between 0 and 100")
    }
    return nil
}

// IsPublished retorna si el assessment está publicado
func (a *Assessment) IsPublished() bool {
    return a.Status == "published"
}
```

---

**Archivo:** `question.go`

```go
package evaluation

import (
    "github.com/google/uuid"
)

// QuestionType define los tipos de preguntas soportados
type QuestionType string

const (
    QuestionTypeMultipleChoice QuestionType = "multiple_choice"
    QuestionTypeTrueFalse      QuestionType = "true_false"
    QuestionTypeShortAnswer    QuestionType = "short_answer"
)

// Question representa una pregunta dentro de un assessment
type Question struct {
    ID           uuid.UUID      `json:"id" bson:"_id"`
    AssessmentID uuid.UUID      `json:"assessment_id" bson:"assessment_id"`
    Type         QuestionType   `json:"type" bson:"type"`
    Text         string         `json:"text" bson:"text"`
    Options      []QuestionOption `json:"options,omitempty" bson:"options,omitempty"` // Solo para multiple_choice
    Position     int            `json:"position" bson:"position"` // Orden de la pregunta (1, 2, 3...)
    Points       int            `json:"points" bson:"points"` // Puntos que vale la pregunta
    Explanation  string         `json:"explanation,omitempty" bson:"explanation,omitempty"` // Feedback/explicación
}

// QuestionOption representa una opción de respuesta (para multiple_choice)
type QuestionOption struct {
    ID        uuid.UUID `json:"id" bson:"_id"`
    Text      string    `json:"text" bson:"text"`
    IsCorrect bool      `json:"is_correct" bson:"is_correct"`
    Position  int       `json:"position" bson:"position"` // Orden de la opción (A, B, C, D)
}

// Validate valida la pregunta
func (q *Question) Validate() error {
    if q.Text == "" {
        return errors.New("question text is required")
    }
    if q.Points < 0 {
        return errors.New("points must be non-negative")
    }
    if q.Type == QuestionTypeMultipleChoice && len(q.Options) < 2 {
        return errors.New("multiple choice questions must have at least 2 options")
    }
    return nil
}

// GetCorrectOptions retorna las opciones correctas
func (q *Question) GetCorrectOptions() []QuestionOption {
    var correct []QuestionOption
    for _, opt := range q.Options {
        if opt.IsCorrect {
            correct = append(correct, opt)
        }
    }
    return correct
}
```

---

**Archivo:** `attempt.go`

```go
package evaluation

import (
    "time"
    "github.com/google/uuid"
)

// Attempt representa un intento de un estudiante en un assessment
type Attempt struct {
    ID           uuid.UUID       `json:"id" bson:"_id"`
    AssessmentID uuid.UUID       `json:"assessment_id" bson:"assessment_id"`
    StudentID    int64           `json:"student_id" bson:"student_id"`
    Answers      []Answer        `json:"answers" bson:"answers"`
    TotalScore   int             `json:"total_score" bson:"total_score"` // Puntos obtenidos
    MaxScore     int             `json:"max_score" bson:"max_score"` // Puntos máximos posibles
    Percentage   float64         `json:"percentage" bson:"percentage"` // Porcentaje (0-100)
    Passed       bool            `json:"passed" bson:"passed"` // Si aprobó según passing_score
    StartedAt    time.Time       `json:"started_at" bson:"started_at"`
    SubmittedAt  *time.Time      `json:"submitted_at,omitempty" bson:"submitted_at,omitempty"`
    DurationSec  int             `json:"duration_sec,omitempty" bson:"duration_sec,omitempty"` // Duración en segundos
}

// Answer representa la respuesta a una pregunta
type Answer struct {
    QuestionID     uuid.UUID `json:"question_id" bson:"question_id"`
    AnswerText     string    `json:"answer_text,omitempty" bson:"answer_text,omitempty"` // Para short_answer
    SelectedOptions []uuid.UUID `json:"selected_options,omitempty" bson:"selected_options,omitempty"` // Para multiple_choice
    IsCorrect      bool      `json:"is_correct" bson:"is_correct"` // Si la respuesta fue correcta
    PointsEarned   int       `json:"points_earned" bson:"points_earned"` // Puntos ganados por esta pregunta
}

// CalculatePercentage calcula el porcentaje basado en score
func (a *Attempt) CalculatePercentage() {
    if a.MaxScore > 0 {
        a.Percentage = (float64(a.TotalScore) / float64(a.MaxScore)) * 100
    } else {
        a.Percentage = 0
    }
}

// CheckPassed verifica si el attempt pasó según el passing score
func (a *Attempt) CheckPassed(passingScore int) {
    a.Passed = a.Percentage >= float64(passingScore)
}

// IsSubmitted retorna si el attempt fue enviado
func (a *Attempt) IsSubmitted() bool {
    return a.SubmittedAt != nil
}
```

---

#### Tests Mínimos Requeridos

**Archivo:** `assessment_test.go`

```go
package evaluation_test

import (
    "testing"
    "github.com/EduGoGroup/edugo-shared/evaluation"
    "github.com/google/uuid"
)

func TestAssessment_Validate(t *testing.T) {
    tests := []struct {
        name    string
        assessment evaluation.Assessment
        wantErr bool
    }{
        {
            name: "valid assessment",
            assessment: evaluation.Assessment{
                ID:           uuid.New(),
                Title:        "Test Quiz",
                PassingScore: 70,
            },
            wantErr: false,
        },
        {
            name: "missing title",
            assessment: evaluation.Assessment{
                ID:           uuid.New(),
                PassingScore: 70,
            },
            wantErr: true,
        },
        {
            name: "invalid passing score",
            assessment: evaluation.Assessment{
                ID:           uuid.New(),
                Title:        "Test Quiz",
                PassingScore: 150, // > 100
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.assessment.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestAssessment_IsPublished(t *testing.T) {
    published := evaluation.Assessment{Status: "published"}
    draft := evaluation.Assessment{Status: "draft"}

    if !published.IsPublished() {
        t.Error("expected published assessment to return true")
    }
    if draft.IsPublished() {
        t.Error("expected draft assessment to return false")
    }
}
```

**Archivo:** `question_test.go`

```go
package evaluation_test

import (
    "testing"
    "github.com/EduGoGroup/edugo-shared/evaluation"
    "github.com/google/uuid"
)

func TestQuestion_Validate(t *testing.T) {
    tests := []struct {
        name     string
        question evaluation.Question
        wantErr  bool
    }{
        {
            name: "valid multiple choice",
            question: evaluation.Question{
                ID:   uuid.New(),
                Text: "What is 2+2?",
                Type: evaluation.QuestionTypeMultipleChoice,
                Options: []evaluation.QuestionOption{
                    {Text: "3", IsCorrect: false},
                    {Text: "4", IsCorrect: true},
                },
                Points: 5,
            },
            wantErr: false,
        },
        {
            name: "multiple choice with < 2 options",
            question: evaluation.Question{
                ID:   uuid.New(),
                Text: "What is 2+2?",
                Type: evaluation.QuestionTypeMultipleChoice,
                Options: []evaluation.QuestionOption{
                    {Text: "4", IsCorrect: true},
                },
                Points: 5,
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.question.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestQuestion_GetCorrectOptions(t *testing.T) {
    q := evaluation.Question{
        Options: []evaluation.QuestionOption{
            {Text: "A", IsCorrect: false},
            {Text: "B", IsCorrect: true},
            {Text: "C", IsCorrect: false},
            {Text: "D", IsCorrect: true},
        },
    }

    correct := q.GetCorrectOptions()
    if len(correct) != 2 {
        t.Errorf("expected 2 correct options, got %d", len(correct))
    }
}
```

**Archivo:** `attempt_test.go`

```go
package evaluation_test

import (
    "testing"
    "github.com/EduGoGroup/edugo-shared/evaluation"
)

func TestAttempt_CalculatePercentage(t *testing.T) {
    attempt := evaluation.Attempt{
        TotalScore: 75,
        MaxScore:   100,
    }

    attempt.CalculatePercentage()

    if attempt.Percentage != 75.0 {
        t.Errorf("expected 75.0, got %.2f", attempt.Percentage)
    }
}

func TestAttempt_CheckPassed(t *testing.T) {
    tests := []struct {
        name         string
        percentage   float64
        passingScore int
        wantPassed   bool
    }{
        {"passed", 80.0, 70, true},
        {"failed", 65.0, 70, false},
        {"exact pass", 70.0, 70, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            attempt := evaluation.Attempt{
                Percentage: tt.percentage,
            }
            attempt.CheckPassed(tt.passingScore)

            if attempt.Passed != tt.wantPassed {
                t.Errorf("expected passed=%v, got %v", tt.wantPassed, attempt.Passed)
            }
        })
    }
}
```

---

#### Versión Inicial

**Tag:** `evaluation/v0.1.0`

**Changelog:**
```
# Changelog - evaluation module

## [0.1.0] - 2025-11-XX

### Added
- Assessment struct with validation
- Question struct with QuestionType enum
- QuestionOption struct for multiple choice
- Attempt struct with scoring logic
- Answer struct for student responses
- Helper methods: Validate, IsPublished, GetCorrectOptions, CalculatePercentage, CheckPassed
- Comprehensive unit tests (>90% coverage)

### Dependencies
- github.com/google/uuid v1.6.0
- github.com/EduGoGroup/edugo-shared/common v0.5.0
```

---

#### Tiempo Estimado

**Implementación:** 3-4 horas
- Estructuras: 1 hora
- Helper methods: 1 hora
- Tests: 1-2 horas

**Documentación:** 1 hora
- Godoc comments
- README.md del módulo
- Ejemplos de uso

**Total:** 4-5 horas

---

#### Checklist de Implementación

- [ ] Crear carpeta `evaluation/`
- [ ] Crear `go.mod` con dependencias
- [ ] Implementar `assessment.go` con struct y métodos
- [ ] Implementar `question.go` con struct y enum
- [ ] Implementar `attempt.go` con struct y scoring
- [ ] Crear `assessment_test.go` con >80% coverage
- [ ] Crear `question_test.go` con >80% coverage
- [ ] Crear `attempt_test.go` con >80% coverage
- [ ] Crear `README.md` del módulo
- [ ] Ejecutar `go test -v -cover ./...`
- [ ] Ejecutar `go mod tidy`
- [ ] Commit y push
- [ ] Crear tag `evaluation/v0.1.0`
- [ ] Publicar en GitHub

---

## 📊 Resumen de Módulos Nuevos

| Módulo | Prioridad | Requerido Por | Tiempo Est. | Complejidad |
|--------|-----------|---------------|-------------|-------------|
| evaluation/ | P0 | api-mobile, worker | 4-5 horas | Media |

**Total módulos nuevos:** 1

**Total tiempo estimado:** 4-5 horas (~1 día)

---

## ✅ Criterios de Aceptación

### Para considerar `evaluation/` completo:

- ✅ Todas las estructuras implementadas (Assessment, Question, QuestionOption, Attempt, Answer)
- ✅ Todos los helper methods implementados
- ✅ Coverage de tests >80%
- ✅ go test pasa sin errores
- ✅ go mod tidy ejecutado
- ✅ README.md con ejemplos de uso
- ✅ Tag evaluation/v0.1.0 publicado
- ✅ api-mobile puede importar y usar: `import "github.com/EduGoGroup/edugo-shared/evaluation"`
- ✅ worker puede importar y usar el módulo

---

## 🚀 Próximos Pasos

1. Crear módulo `evaluation/` en Sprint 1
2. Validar que api-mobile y worker pueden compilar con el módulo
3. Documentar en `02-NECESIDADES_CONSOLIDADAS.md` como ✅ Implementado
4. Actualizar `01-ESTADO_ACTUAL.md` con nuevo módulo

---

**Documento generado:** 15 de Noviembre, 2025  
**Próximo documento:** `04-FEATURES_FALTANTES.md`
