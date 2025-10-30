package enum

// AssessmentType representa el tipo de pregunta en un assessment
type AssessmentType string

const (
	AssessmentTypeMultipleChoice AssessmentType = "multiple_choice"
	AssessmentTypeTrueFalse      AssessmentType = "true_false"
	AssessmentTypeShortAnswer    AssessmentType = "short_answer"
)

// IsValid verifica si el tipo es válido
func (a AssessmentType) IsValid() bool {
	switch a {
	case AssessmentTypeMultipleChoice, AssessmentTypeTrueFalse, AssessmentTypeShortAnswer:
		return true
	}
	return false
}

// String retorna la representación en string del tipo
func (a AssessmentType) String() string {
	return string(a)
}

// AllAssessmentTypes retorna todos los tipos válidos
func AllAssessmentTypes() []AssessmentType {
	return []AssessmentType{
		AssessmentTypeMultipleChoice,
		AssessmentTypeTrueFalse,
		AssessmentTypeShortAnswer,
	}
}
