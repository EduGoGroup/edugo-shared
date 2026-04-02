package enum

// SystemRole representa los roles del sistema, alineados 1:1 con la tabla iam.roles
type SystemRole string

const (
	// SystemRoleSuperAdmin es el rol de super administrador de la plataforma
	SystemRoleSuperAdmin SystemRole = "super_admin"
	// SystemRolePlatformAdmin es el rol de administrador de la plataforma
	SystemRolePlatformAdmin SystemRole = "platform_admin"

	// SystemRoleSchoolAdmin es el rol de administrador de escuela
	SystemRoleSchoolAdmin SystemRole = "school_admin"
	// SystemRoleSchoolDirector es el rol de director de escuela
	SystemRoleSchoolDirector SystemRole = "school_director"
	// SystemRoleSchoolCoordinator es el rol de coordinador de escuela
	SystemRoleSchoolCoordinator SystemRole = "school_coordinator"
	// SystemRoleSchoolAssistant es el rol de asistente administrativo de escuela
	SystemRoleSchoolAssistant SystemRole = "school_assistant"

	// SystemRoleTeacher es el rol de profesor
	SystemRoleTeacher SystemRole = "teacher"
	// SystemRoleAssistantTeacher es el rol de profesor asistente
	SystemRoleAssistantTeacher SystemRole = "assistant_teacher"
	// SystemRoleStudent es el rol de estudiante
	SystemRoleStudent SystemRole = "student"
	// SystemRoleGuardian es el rol de tutor/padre de familia
	SystemRoleGuardian SystemRole = "guardian"

	// SystemRoleObserver es el rol de observador
	SystemRoleObserver SystemRole = "observer"
	// SystemRoleReadonlyAuditor es el rol de auditor de solo lectura
	SystemRoleReadonlyAuditor SystemRole = "readonly_auditor"
)

// IsValid verifica si el rol es valido
func (r SystemRole) IsValid() bool {
	switch r {
	case SystemRoleSuperAdmin, SystemRolePlatformAdmin,
		SystemRoleSchoolAdmin, SystemRoleSchoolDirector, SystemRoleSchoolCoordinator, SystemRoleSchoolAssistant,
		SystemRoleTeacher, SystemRoleAssistantTeacher, SystemRoleStudent, SystemRoleGuardian,
		SystemRoleObserver, SystemRoleReadonlyAuditor:
		return true
	}
	return false
}

// String retorna la representacion en string del rol
func (r SystemRole) String() string {
	return string(r)
}

