// Package enum provides enumeration types and constants for various
// domain entities in the EduGo shared library.
package enum

// AssessmentType representa el tipo de pregunta en un assessment
type AssessmentType string

const (
	// AssessmentTypeMultipleChoice represents a multiple choice question
	AssessmentTypeMultipleChoice AssessmentType = "multiple_choice"
	// AssessmentTypeTrueFalse represents a true/false question
	AssessmentTypeTrueFalse AssessmentType = "true_false"
	// AssessmentTypeShortAnswer represents a short answer question
	AssessmentTypeShortAnswer AssessmentType = "short_answer"
	// AssessmentTypeOpenEnded represents an open ended question
	AssessmentTypeOpenEnded AssessmentType = "open_ended"
)

// IsValid verifica si el tipo es válido
func (a AssessmentType) IsValid() bool {
	switch a {
	case AssessmentTypeMultipleChoice, AssessmentTypeTrueFalse, AssessmentTypeShortAnswer, AssessmentTypeOpenEnded:
		return true
	}
	return false
}

// String retorna la representación en string del tipo
func (a AssessmentType) String() string {
	return string(a)
}
