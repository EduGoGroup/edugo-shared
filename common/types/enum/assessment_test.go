package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssessmentType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		typ  AssessmentType
		want bool
	}{
		{"MultipleChoice", AssessmentTypeMultipleChoice, true},
		{"TrueFalse", AssessmentTypeTrueFalse, true},
		{"ShortAnswer", AssessmentTypeShortAnswer, true},
		{"OpenEnded", AssessmentTypeOpenEnded, true},
		{"Invalid", "invalid_type", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.typ.IsValid())
		})
	}
}

func TestAssessmentType_String(t *testing.T) {
	assert.Equal(t, "multiple_choice", AssessmentTypeMultipleChoice.String())
}

func TestAllAssessmentTypes(t *testing.T) {
	types := AllAssessmentTypes()
	// Should contain OpenEnded, so length should be 4
	assert.Contains(t, types, AssessmentTypeOpenEnded)
}
