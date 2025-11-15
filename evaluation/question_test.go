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
			name: "valid true/false",
			question: evaluation.Question{
				ID:     uuid.New(),
				Text:   "The sky is blue",
				Type:   evaluation.QuestionTypeTrueFalse,
				Points: 5,
			},
			wantErr: false,
		},
		{
			name: "valid short answer",
			question: evaluation.Question{
				ID:     uuid.New(),
				Text:   "What is the capital of France?",
				Type:   evaluation.QuestionTypeShortAnswer,
				Points: 10,
			},
			wantErr: false,
		},
		{
			name: "missing text",
			question: evaluation.Question{
				ID:   uuid.New(),
				Type: evaluation.QuestionTypeMultipleChoice,
				Options: []evaluation.QuestionOption{
					{Text: "A", IsCorrect: true},
					{Text: "B", IsCorrect: false},
				},
				Points: 5,
			},
			wantErr: true,
		},
		{
			name: "negative points",
			question: evaluation.Question{
				ID:     uuid.New(),
				Text:   "Test question",
				Type:   evaluation.QuestionTypeShortAnswer,
				Points: -5,
			},
			wantErr: true,
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
		{
			name: "multiple choice with 0 options",
			question: evaluation.Question{
				ID:      uuid.New(),
				Text:    "What is 2+2?",
				Type:    evaluation.QuestionTypeMultipleChoice,
				Options: []evaluation.QuestionOption{},
				Points:  5,
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
	tests := []struct {
		name          string
		question      evaluation.Question
		wantCorrectCount int
	}{
		{
			name: "two correct options",
			question: evaluation.Question{
				Options: []evaluation.QuestionOption{
					{Text: "A", IsCorrect: false},
					{Text: "B", IsCorrect: true},
					{Text: "C", IsCorrect: false},
					{Text: "D", IsCorrect: true},
				},
			},
			wantCorrectCount: 2,
		},
		{
			name: "one correct option",
			question: evaluation.Question{
				Options: []evaluation.QuestionOption{
					{Text: "A", IsCorrect: false},
					{Text: "B", IsCorrect: true},
					{Text: "C", IsCorrect: false},
				},
			},
			wantCorrectCount: 1,
		},
		{
			name: "no correct options",
			question: evaluation.Question{
				Options: []evaluation.QuestionOption{
					{Text: "A", IsCorrect: false},
					{Text: "B", IsCorrect: false},
				},
			},
			wantCorrectCount: 0,
		},
		{
			name: "all correct options",
			question: evaluation.Question{
				Options: []evaluation.QuestionOption{
					{Text: "A", IsCorrect: true},
					{Text: "B", IsCorrect: true},
				},
			},
			wantCorrectCount: 2,
		},
		{
			name:             "empty options",
			question:         evaluation.Question{Options: []evaluation.QuestionOption{}},
			wantCorrectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			correct := tt.question.GetCorrectOptions()
			if len(correct) != tt.wantCorrectCount {
				t.Errorf("GetCorrectOptions() returned %d options, want %d", len(correct), tt.wantCorrectCount)
			}
			// Verify all returned options are actually correct
			for _, opt := range correct {
				if !opt.IsCorrect {
					t.Errorf("GetCorrectOptions() returned option %q which is not correct", opt.Text)
				}
			}
		})
	}
}
