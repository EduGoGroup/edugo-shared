package enum

// SystemRole representa los roles del sistema
type SystemRole string

const (
	SystemRoleAdmin    SystemRole = "admin"
	SystemRoleTeacher  SystemRole = "teacher"
	SystemRoleStudent  SystemRole = "student"
	SystemRoleGuardian SystemRole = "guardian"
)

// IsValid verifica si el rol es válido
func (r SystemRole) IsValid() bool {
	switch r {
	case SystemRoleAdmin, SystemRoleTeacher, SystemRoleStudent, SystemRoleGuardian:
		return true
	}
	return false
}

// String retorna la representación en string del rol
func (r SystemRole) String() string {
	return string(r)
}

// AllSystemRoles retorna todos los roles válidos
func AllSystemRoles() []SystemRole {
	return []SystemRole{
		SystemRoleAdmin,
		SystemRoleTeacher,
		SystemRoleStudent,
		SystemRoleGuardian,
	}
}

// AllSystemRolesStrings retorna todos los roles como strings (útil para validación)
func AllSystemRolesStrings() []string {
	return []string{
		string(SystemRoleAdmin),
		string(SystemRoleTeacher),
		string(SystemRoleStudent),
		string(SystemRoleGuardian),
	}
}
