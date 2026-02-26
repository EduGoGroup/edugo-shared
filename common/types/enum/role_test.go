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
		{"Admin", SystemRoleAdmin, true},
		{"Teacher", SystemRoleTeacher, true},
		{"Student", SystemRoleStudent, true},
		{"Guardian", SystemRoleGuardian, true},
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
	assert.Equal(t, "admin", SystemRoleAdmin.String())
	assert.Equal(t, "teacher", SystemRoleTeacher.String())
}

func TestAllSystemRoles(t *testing.T) {
	roles := AllSystemRoles()
	assert.Len(t, roles, 4)
	assert.Contains(t, roles, SystemRoleAdmin)
	assert.Contains(t, roles, SystemRoleTeacher)
	assert.Contains(t, roles, SystemRoleStudent)
	assert.Contains(t, roles, SystemRoleGuardian)
}

func TestAllSystemRolesStrings(t *testing.T) {
	roles := AllSystemRolesStrings()
	assert.Len(t, roles, 4)
	assert.Contains(t, roles, "admin")
	assert.Contains(t, roles, "teacher")
	assert.Contains(t, roles, "student")
	assert.Contains(t, roles, "guardian")
}
