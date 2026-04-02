package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaterialStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status MaterialStatus
		want   bool
	}{
		{"Draft", MaterialStatusDraft, true},
		{"Published", MaterialStatusPublished, true},
		{"Archived", MaterialStatusArchived, true},
		{"Invalid", "invalid_status", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.status.IsValid())
		})
	}
}

func TestMaterialStatus_String(t *testing.T) {
	assert.Equal(t, "draft", MaterialStatusDraft.String())
}

func TestProgressStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status ProgressStatus
		want   bool
	}{
		{"NotStarted", ProgressStatusNotStarted, true},
		{"InProgress", ProgressStatusInProgress, true},
		{"Completed", ProgressStatusCompleted, true},
		{"Invalid", "invalid_status", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.status.IsValid())
		})
	}
}

func TestProgressStatus_String(t *testing.T) {
	assert.Equal(t, "not_started", ProgressStatusNotStarted.String())
}

func TestProcessingStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status ProcessingStatus
		want   bool
	}{
		{"Pending", ProcessingStatusPending, true},
		{"Processing", ProcessingStatusProcessing, true},
		{"Completed", ProcessingStatusCompleted, true},
		{"Failed", ProcessingStatusFailed, true},
		{"Invalid", "invalid_status", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.status.IsValid())
		})
	}
}

func TestProcessingStatus_String(t *testing.T) {
	assert.Equal(t, "pending", ProcessingStatusPending.String())
}

