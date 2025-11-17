package evaluation_test

import (
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-shared/evaluation"
)

func TestAttempt_CalculatePercentage(t *testing.T) {
	tests := []struct {
		name           string
		totalScore     float64
		maxScore       int
		wantPercentage float64
	}{
		{"75 out of 100", 75.0, 100, 75.0},
		{"50 out of 100", 50.0, 100, 50.0},
		{"100 out of 100", 100.0, 100, 100.0},
		{"0 out of 100", 0.0, 100, 0.0},
		{"zero max score", 75.0, 0, 0.0},
		{"50 out of 50", 50.0, 50, 100.0},
		{"25 out of 50", 25.0, 50, 50.0},
		{"decimal score 85.5 out of 100", 85.5, 100, 85.5},
		{"decimal score 92.75 out of 100", 92.75, 100, 92.75},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attempt := evaluation.Attempt{
				TotalScore: tt.totalScore,
				MaxScore:   tt.maxScore,
			}

			attempt.CalculatePercentage()

			if attempt.Percentage != tt.wantPercentage {
				t.Errorf("CalculatePercentage() = %.2f, want %.2f", attempt.Percentage, tt.wantPercentage)
			}
		})
	}
}

func TestAttempt_CheckPassed(t *testing.T) {
	tests := []struct {
		name         string
		percentage   float64
		passingScore int
		wantPassed   bool
	}{
		{"passed with exact score", 70.0, 70, true},
		{"passed above score", 80.0, 70, true},
		{"failed below score", 65.0, 70, false},
		{"passed with 100%", 100.0, 70, true},
		{"failed with 0%", 0.0, 70, false},
		{"passed with 0% passing score", 50.0, 0, true},
		{"failed with 100% passing score", 99.0, 100, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attempt := evaluation.Attempt{
				Percentage: tt.percentage,
			}
			attempt.CheckPassed(tt.passingScore)

			if attempt.Passed != tt.wantPassed {
				t.Errorf("CheckPassed() set Passed=%v, want %v", attempt.Passed, tt.wantPassed)
			}
		})
	}
}

func TestAttempt_IsSubmitted(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		submittedAt *time.Time
		want        bool
	}{
		{"submitted", &now, true},
		{"not submitted", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attempt := evaluation.Attempt{
				SubmittedAt: tt.submittedAt,
			}

			if got := attempt.IsSubmitted(); got != tt.want {
				t.Errorf("IsSubmitted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttempt_CalculateAndCheckPassed(t *testing.T) {
	// Integration test: calculate percentage and check if passed
	attempt := evaluation.Attempt{
		TotalScore: 85.0,
		MaxScore:   100,
	}

	attempt.CalculatePercentage()
	attempt.CheckPassed(70)

	if attempt.Percentage != 85.0 {
		t.Errorf("Percentage = %.2f, want 85.0", attempt.Percentage)
	}
	if !attempt.Passed {
		t.Error("Expected attempt to be marked as passed")
	}
}
