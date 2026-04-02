package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemRole_IsValid(t *testing.T) {
	tests := []struct {
		name string
		role SystemRole
		want bool
	}{
		// Platform-level roles
		{"SuperAdmin", SystemRoleSuperAdmin, true},
		{"PlatformAdmin", SystemRolePlatformAdmin, true},

		// School-level administrative roles
		{"SchoolAdmin", SystemRoleSchoolAdmin, true},
		{"SchoolDirector", SystemRoleSchoolDirector, true},
		{"SchoolCoordinator", SystemRoleSchoolCoordinator, true},
		{"SchoolAssistant", SystemRoleSchoolAssistant, true},

		// Academic roles
		{"Teacher", SystemRoleTeacher, true},
		{"AssistantTeacher", SystemRoleAssistantTeacher, true},
		{"Student", SystemRoleStudent, true},
		{"Guardian", SystemRoleGuardian, true},

		// Observation roles
		{"Observer", SystemRoleObserver, true},
		{"ReadonlyAuditor", SystemRoleReadonlyAuditor, true},

		// Invalid roles
		{"OldAdmin", SystemRole("admin"), false},
		{"Invalid", "invalid_role", false},
		{"Empty", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.role.IsValid())
		})
	}
}

func TestSystemRole_String(t *testing.T) {
	assert.Equal(t, "super_admin", SystemRoleSuperAdmin.String())
	assert.Equal(t, "platform_admin", SystemRolePlatformAdmin.String())
	assert.Equal(t, "school_admin", SystemRoleSchoolAdmin.String())
	assert.Equal(t, "teacher", SystemRoleTeacher.String())
	assert.Equal(t, "student", SystemRoleStudent.String())
}

