package evaluation

import (
	"errors"

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
	ID           uuid.UUID        `json:"id" bson:"_id"`
	AssessmentID uuid.UUID        `json:"assessment_id" bson:"assessment_id"`
	Type         QuestionType     `json:"type" bson:"type"`
	Text         string           `json:"text" bson:"text"`
	Options      []QuestionOption `json:"options,omitempty" bson:"options,omitempty"` // Solo para multiple_choice
	Position     int              `json:"position" bson:"position"`                   // Orden de la pregunta (1, 2, 3...)
	Points       int              `json:"points" bson:"points"`                       // Puntos que vale la pregunta
	Explanation  string           `json:"explanation,omitempty" bson:"explanation,omitempty"` // Feedback/explicación
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
