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

func TestAllSystemRoles(t *testing.T) {
	roles := AllSystemRoles()
	assert.Len(t, roles, 12)
	assert.Contains(t, roles, SystemRoleSuperAdmin)
	assert.Contains(t, roles, SystemRolePlatformAdmin)
	assert.Contains(t, roles, SystemRoleSchoolAdmin)
	assert.Contains(t, roles, SystemRoleSchoolDirector)
	assert.Contains(t, roles, SystemRoleSchoolCoordinator)
	assert.Contains(t, roles, SystemRoleSchoolAssistant)
	assert.Contains(t, roles, SystemRoleTeacher)
	assert.Contains(t, roles, SystemRoleAssistantTeacher)
	assert.Contains(t, roles, SystemRoleStudent)
	assert.Contains(t, roles, SystemRoleGuardian)
	assert.Contains(t, roles, SystemRoleObserver)
	assert.Contains(t, roles, SystemRoleReadonlyAuditor)
}

func TestAllSystemRolesStrings(t *testing.T) {
	roles := AllSystemRolesStrings()
	assert.Len(t, roles, 12)
	assert.Contains(t, roles, "super_admin")
	assert.Contains(t, roles, "platform_admin")
	assert.Contains(t, roles, "school_admin")
	assert.Contains(t, roles, "school_director")
	assert.Contains(t, roles, "school_coordinator")
	assert.Contains(t, roles, "school_assistant")
	assert.Contains(t, roles, "teacher")
	assert.Contains(t, roles, "assistant_teacher")
	assert.Contains(t, roles, "student")
	assert.Contains(t, roles, "guardian")
	assert.Contains(t, roles, "observer")
	assert.Contains(t, roles, "readonly_auditor")
}
