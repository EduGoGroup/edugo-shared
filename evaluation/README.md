# evaluation - EduGo Shared

Módulo que define modelos compartidos para el sistema de evaluaciones de EduGo.

## Instalación

```bash
go get github.com/EduGoGroup/edugo-shared/evaluation@v0.7.0
```

## Uso

### Assessment (Cuestionario)

```go
import "github.com/EduGoGroup/edugo-shared/evaluation"

assessment := evaluation.Assessment{
    ID:           uuid.New(),
    MaterialID:   123,
    Title:        "Quiz de Matemáticas",
    Description:  "Evaluación del capítulo 1",
    Type:         "generated",
    Status:       "published",
    PassingScore: 70,
    TotalQuestions: 10,
    TotalPoints:  100,
    CreatedBy:    456,
    CreatedAt:    time.Now(),
    UpdatedAt:    time.Now(),
}

// Validar
if err := assessment.Validate(); err != nil {
    log.Fatal(err)
}

// Verificar si está publicado
if assessment.IsPublished() {
    fmt.Println("Assessment está publicado")
}
```

### Question (Pregunta)

```go
// Pregunta de opción múltiple
question := evaluation.Question{
    ID:           uuid.New(),
    AssessmentID: assessmentID,
    Type:         evaluation.QuestionTypeMultipleChoice,
    Text:         "¿Cuánto es 2 + 2?",
    Options: []evaluation.QuestionOption{
        {ID: uuid.New(), Text: "3", IsCorrect: false, Position: 1},
        {ID: uuid.New(), Text: "4", IsCorrect: true, Position: 2},
        {ID: uuid.New(), Text: "5", IsCorrect: false, Position: 3},
    },
    Position:    1,
    Points:      10,
    Explanation: "2 + 2 = 4",
}

// Validar
if err := question.Validate(); err != nil {
    log.Fatal(err)
}

// Obtener opciones correctas
correctOptions := question.GetCorrectOptions()
fmt.Printf("Hay %d opciones correctas\n", len(correctOptions))
```

### Attempt (Intento de estudiante)

```go
attempt := evaluation.Attempt{
    ID:           uuid.New(),
    AssessmentID: assessmentID,
    StudentID:    789,
    Answers: []evaluation.Answer{
        {
            QuestionID:      questionID,
            SelectedOptions: []uuid.UUID{optionID},
            IsCorrect:       true,
            PointsEarned:    10,
        },
    },
    TotalScore: 85,
    MaxScore:   100,
    StartedAt:  time.Now().Add(-30 * time.Minute),
}

// Calcular porcentaje
attempt.CalculatePercentage()
fmt.Printf("Porcentaje: %.2f%%\n", attempt.Percentage)

// Verificar si aprobó (passing score = 70%)
attempt.CheckPassed(70)
if attempt.Passed {
    fmt.Println("¡Aprobado!")
}

// Marcar como enviado
now := time.Now()
attempt.SubmittedAt = &now
if attempt.IsSubmitted() {
    fmt.Println("Attempt enviado")
}
```

## Tipos de Preguntas

El módulo soporta 3 tipos de preguntas:

- `QuestionTypeMultipleChoice`: Opción múltiple (requiere al menos 2 opciones)
- `QuestionTypeTrueFalse`: Verdadero/Falso
- `QuestionTypeShortAnswer`: Respuesta corta

## Modelos

### Assessment
Representa un cuestionario (generado por IA o manual).

**Campos principales:**
- `ID`: UUID del assessment
- `MaterialID`: ID del material educativo asociado
- `Title`: Título del cuestionario
- `PassingScore`: Porcentaje mínimo para aprobar (0-100)
- `Status`: "draft", "published", "archived"

### Question
Representa una pregunta dentro de un assessment.

**Campos principales:**
- `AssessmentID`: UUID del assessment padre
- `Type`: Tipo de pregunta
- `Text`: Texto de la pregunta
- `Options`: Opciones de respuesta (solo para multiple_choice)
- `Points`: Puntos que vale la pregunta

### Attempt
Representa un intento de un estudiante en un assessment.

**Campos principales:**
- `AssessmentID`: UUID del assessment
- `StudentID`: ID del estudiante
- `Answers`: Lista de respuestas
- `TotalScore`: Puntos obtenidos
- `Percentage`: Porcentaje de acierto
- `Passed`: Si aprobó o no

## Testing

El módulo tiene **100% de cobertura de tests**.

```bash
go test -v -cover ./...
```

## Compatibilidad

- Go 1.24+
- Compatible con MongoDB (tags bson)
- Compatible con JSON marshaling

## Licencia

Uso interno de EduGo
