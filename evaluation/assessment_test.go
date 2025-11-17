package evaluation_test

import (
	"testing"

	"github.com/EduGoGroup/edugo-shared/evaluation"
	"github.com/google/uuid"
)

func TestAssessment_Validate(t *testing.T) {
	materialID := uuid.New()
	
	tests := []struct {
		name       string
		assessment evaluation.Assessment
		wantErr    bool
	}{
		{
			name: "valid assessment",
			assessment: evaluation.Assessment{
				ID:           uuid.New(),
				MaterialID:   materialID,
				Title:        "Test Quiz",
				PassingScore: 70,
			},
			wantErr: false,
		},
		{
			name: "missing title",
			assessment: evaluation.Assessment{
				ID:           uuid.New(),
				MaterialID:   materialID,
				PassingScore: 70,
			},
			wantErr: true,
		},
		{
			name: "invalid passing score - too high",
			assessment: evaluation.Assessment{
				ID:           uuid.New(),
				MaterialID:   materialID,
				Title:        "Test Quiz",
				PassingScore: 150, // > 100
			},
			wantErr: true,
		},
		{
			name: "invalid passing score - negative",
			assessment: evaluation.Assessment{
				ID:           uuid.New(),
				MaterialID:   materialID,
				Title:        "Test Quiz",
				PassingScore: -10,
			},
			wantErr: true,
		},
		{
			name: "valid assessment with 0 passing score",
			assessment: evaluation.Assessment{
				ID:           uuid.New(),
				MaterialID:   materialID,
				Title:        "Test Quiz",
				PassingScore: 0,
			},
			wantErr: false,
		},
		{
			name: "valid assessment with 100 passing score",
			assessment: evaluation.Assessment{
				ID:           uuid.New(),
				MaterialID:   materialID,
				Title:        "Test Quiz",
				PassingScore: 100,
			},
			wantErr: false,
		},
		{
			name: "valid assessment with max attempts",
			assessment: evaluation.Assessment{
				ID:           uuid.New(),
				MaterialID:   materialID,
				Title:        "Test Quiz",
				PassingScore: 70,
				MaxAttempts:  intPtr(3),
			},
			wantErr: false,
		},
		{
			name: "valid assessment with time limit",
			assessment: evaluation.Assessment{
				ID:               uuid.New(),
				MaterialID:       materialID,
				Title:            "Test Quiz",
				PassingScore:     70,
				TimeLimitMinutes: intPtr(30),
			},
			wantErr: false,
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
	tests := []struct {
		name       string
		status     string
		wantResult bool
	}{
		{"published assessment", "published", true},
		{"draft assessment", "draft", false},
		{"archived assessment", "archived", false},
		{"empty status", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assessment := evaluation.Assessment{Status: tt.status}
			if got := assessment.IsPublished(); got != tt.wantResult {
				t.Errorf("IsPublished() = %v, want %v", got, tt.wantResult)
			}
		})
	}
}

// Helper function for nullable int pointers
func intPtr(i int) *int {
	return &i
}

func TestAssessment_Validate_InvalidNullableFields(t *testing.T) {
	materialID := uuid.New()
	
	tests := []struct {
		name       string
		assessment evaluation.Assessment
		wantErr    bool
		errMsg     string
	}{
		{
			name: "invalid max attempts - zero",
			assessment: evaluation.Assessment{
				ID:           uuid.New(),
				MaterialID:   materialID,
				Title:        "Test Quiz",
				PassingScore: 70,
				MaxAttempts:  intPtr(0),
			},
			wantErr: true,
			errMsg:  "max attempts must be greater than 0",
		},
		{
			name: "invalid max attempts - negative",
			assessment: evaluation.Assessment{
				ID:           uuid.New(),
				MaterialID:   materialID,
				Title:        "Test Quiz",
				PassingScore: 70,
				MaxAttempts:  intPtr(-1),
			},
			wantErr: true,
			errMsg:  "max attempts must be greater than 0",
		},
		{
			name: "invalid time limit - zero",
			assessment: evaluation.Assessment{
				ID:               uuid.New(),
				MaterialID:       materialID,
				Title:            "Test Quiz",
				PassingScore:     70,
				TimeLimitMinutes: intPtr(0),
			},
			wantErr: true,
			errMsg:  "time limit minutes must be greater than 0",
		},
		{
			name: "invalid time limit - negative",
			assessment: evaluation.Assessment{
				ID:               uuid.New(),
				MaterialID:       materialID,
				Title:            "Test Quiz",
				PassingScore:     70,
				TimeLimitMinutes: intPtr(-5),
			},
			wantErr: true,
			errMsg:  "time limit minutes must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.assessment.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errMsg {
				t.Errorf("Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}
