package evaluation

import (
	"time"

	"github.com/google/uuid"
)

// Attempt representa un intento de un estudiante en un assessment
type Attempt struct {
	ID             uuid.UUID  `json:"id" bson:"_id"`
	AssessmentID   uuid.UUID  `json:"assessment_id" bson:"assessment_id"`
	UserID         int64      `json:"user_id" bson:"user_id"`                                   // BREAKING CHANGE: StudentID → UserID
	Answers        []Answer   `json:"answers" bson:"answers"`
	TotalScore     float64    `json:"total_score" bson:"total_score"`                           // BREAKING CHANGE: int → float64 para scores decimales
	MaxScore       int        `json:"max_score" bson:"max_score"`                               // Puntos máximos posibles
	Percentage     float64    `json:"percentage" bson:"percentage"`                             // Porcentaje (0-100)
	Passed         bool       `json:"passed" bson:"passed"`                                     // Si aprobó según passing_score
	StartedAt      time.Time  `json:"started_at" bson:"started_at"`
	SubmittedAt    *time.Time `json:"submitted_at,omitempty" bson:"submitted_at,omitempty"`
	DurationSec    int        `json:"duration_sec,omitempty" bson:"duration_sec,omitempty"`         // Duración en segundos
	IdempotencyKey string     `json:"idempotency_key,omitempty" bson:"idempotency_key,omitempty"` // Para prevenir duplicados
}

// Answer representa la respuesta a una pregunta
type Answer struct {
	QuestionID      uuid.UUID   `json:"question_id" bson:"question_id"`
	AnswerText      string      `json:"answer_text,omitempty" bson:"answer_text,omitempty"`             // Para short_answer
	SelectedOptions []uuid.UUID `json:"selected_options,omitempty" bson:"selected_options,omitempty"` // Para multiple_choice
	IsCorrect       bool        `json:"is_correct" bson:"is_correct"`                                   // Si la respuesta fue correcta
	PointsEarned    int         `json:"points_earned" bson:"points_earned"`                             // Puntos ganados por esta pregunta
}

// CalculatePercentage calcula el porcentaje basado en score
func (a *Attempt) CalculatePercentage() {
	if a.MaxScore > 0 {
		a.Percentage = (a.TotalScore / float64(a.MaxScore)) * 100
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
