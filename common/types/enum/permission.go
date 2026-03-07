package enum

// Permission representa un permiso del sistema RBAC
type Permission string

// Permisos de usuarios
const (
	PermissionUsersCreate    Permission = "users:create"
	PermissionUsersRead      Permission = "users:read"
	PermissionUsersUpdate    Permission = "users:update"
	PermissionUsersDelete    Permission = "users:delete"
	PermissionUsersReadOwn   Permission = "users:read:own"
	PermissionUsersUpdateOwn Permission = "users:update:own"
)

// Permisos de escuelas
const (
	PermissionSchoolsCreate Permission = "schools:create"
	PermissionSchoolsRead   Permission = "schools:read"
	PermissionSchoolsUpdate Permission = "schools:update"
	PermissionSchoolsDelete Permission = "schools:delete"
	PermissionSchoolsManage Permission = "schools:manage"
)

// Permisos de unidades académicas
const (
	PermissionUnitsCreate Permission = "units:create"
	PermissionUnitsRead   Permission = "units:read"
	PermissionUnitsUpdate Permission = "units:update"
	PermissionUnitsDelete Permission = "units:delete"
)

// Permisos de materiales
const (
	PermissionMaterialsCreate   Permission = "materials:create"
	PermissionMaterialsRead     Permission = "materials:read"
	PermissionMaterialsUpdate   Permission = "materials:update"
	PermissionMaterialsDelete   Permission = "materials:delete"
	PermissionMaterialsPublish  Permission = "materials:publish"
	PermissionMaterialsDownload Permission = "materials:download"
)

// Permisos de evaluaciones
const (
	PermissionAssessmentsCreate      Permission = "assessments:create"
	PermissionAssessmentsRead        Permission = "assessments:read"
	PermissionAssessmentsUpdate      Permission = "assessments:update"
	PermissionAssessmentsDelete      Permission = "assessments:delete"
	PermissionAssessmentsPublish     Permission = "assessments:publish"
	PermissionAssessmentsGrade       Permission = "assessments:grade"
	PermissionAssessmentsAttempt     Permission = "assessments:attempt"
	PermissionAssessmentsViewResults Permission = "assessments:view_results"
)

// Permisos de progreso
const (
	PermissionProgressRead    Permission = "progress:read"
	PermissionProgressUpdate  Permission = "progress:update"
	PermissionProgressReadOwn Permission = "progress:read:own"
)

// Permisos de estadísticas
const (
	PermissionStatsGlobal Permission = "stats:global"
	PermissionStatsSchool Permission = "stats:school"
	PermissionStatsUnit   Permission = "stats:unit"
)

// Permisos de screen templates
const (
	PermissionScreenTemplatesRead   Permission = "screen_templates:read"
	PermissionScreenTemplatesCreate Permission = "screen_templates:create"
	PermissionScreenTemplatesUpdate Permission = "screen_templates:update"
	PermissionScreenTemplatesDelete Permission = "screen_templates:delete"
)

// Permisos de screen instances
const (
	PermissionScreenInstancesRead   Permission = "screen_instances:read"
	PermissionScreenInstancesCreate Permission = "screen_instances:create"
	PermissionScreenInstancesUpdate Permission = "screen_instances:update"
	PermissionScreenInstancesDelete Permission = "screen_instances:delete"
)

// Permisos de screens (lectura combinada)
const (
	PermissionScreensRead Permission = "screens:read"
)

// Permisos de roles
const (
	PermissionRolesCreate Permission = "roles:create"
	PermissionRolesRead   Permission = "roles:read"
	PermissionRolesUpdate Permission = "roles:update"
	PermissionRolesDelete Permission = "roles:delete"
)

// Gestión de permisos
const (
	PermissionPermissionsMgmtRead   Permission = "permissions_mgmt:read"
	PermissionPermissionsMgmtCreate Permission = "permissions_mgmt:create"
	PermissionPermissionsMgmtUpdate Permission = "permissions_mgmt:update"
	PermissionPermissionsMgmtDelete Permission = "permissions_mgmt:delete"

	// Deprecated: usar PermissionPermissionsMgmtRead
	PermissionResourcesRead = PermissionPermissionsMgmtRead
	// Deprecated: usar PermissionPermissionsMgmtUpdate
	PermissionResourcesUpdate = PermissionPermissionsMgmtUpdate
)

// Permisos de membresías
const (
	PermissionMembershipsCreate Permission = "memberships:create"
	PermissionMembershipsRead   Permission = "memberships:read"
	PermissionMembershipsUpdate Permission = "memberships:update"
	PermissionMembershipsDelete Permission = "memberships:delete"
)

// Permisos de materias
const (
	PermissionSubjectsCreate Permission = "subjects:create"
	PermissionSubjectsRead   Permission = "subjects:read"
	PermissionSubjectsUpdate Permission = "subjects:update"
	PermissionSubjectsDelete Permission = "subjects:delete"
)

// Permisos de vínculos guardian-estudiante
const (
	PermissionGuardianRelationsRead    Permission = "guardian_relations:read"
	PermissionGuardianRelationsApprove Permission = "guardian_relations:approve"
	PermissionGuardianRelationsRequest Permission = "guardian_relations:request"
	PermissionGuardianRelationsManage  Permission = "guardian_relations:manage"
)

// Permisos adicionales de evaluaciones
const (
	PermissionAssessmentsTake    Permission = "assessments:take"
	PermissionAssessmentsArchive Permission = "assessments:archive"
)

// Permisos de tipos de concepto
const (
	PermissionConceptTypesCreate Permission = "concept_types:create"
	PermissionConceptTypesRead   Permission = "concept_types:read"
	PermissionConceptTypesUpdate Permission = "concept_types:update"
	PermissionConceptTypesDelete Permission = "concept_types:delete"
)

// Permisos de auditoría
const (
	PermissionAuditRead   Permission = "audit:read"
	PermissionAuditExport Permission = "audit:export"
)

// String retorna la representación en string del permiso
func (p Permission) String() string {
	return string(p)
}

// IsValid verifica si el permiso es válido
func (p Permission) IsValid() bool {
	return AllPermissions[p]
}

// AllPermissions es un mapa de todos los permisos válidos
var AllPermissions = map[Permission]bool{
	// Usuarios
	PermissionUsersCreate:    true,
	PermissionUsersRead:      true,
	PermissionUsersUpdate:    true,
	PermissionUsersDelete:    true,
	PermissionUsersReadOwn:   true,
	PermissionUsersUpdateOwn: true,
	// Escuelas
	PermissionSchoolsCreate: true,
	PermissionSchoolsRead:   true,
	PermissionSchoolsUpdate: true,
	PermissionSchoolsDelete: true,
	PermissionSchoolsManage: true,
	// Unidades
	PermissionUnitsCreate: true,
	PermissionUnitsRead:   true,
	PermissionUnitsUpdate: true,
	PermissionUnitsDelete: true,
	// Materiales
	PermissionMaterialsCreate:   true,
	PermissionMaterialsRead:     true,
	PermissionMaterialsUpdate:   true,
	PermissionMaterialsDelete:   true,
	PermissionMaterialsPublish:  true,
	PermissionMaterialsDownload: true,
	// Evaluaciones
	PermissionAssessmentsCreate:      true,
	PermissionAssessmentsRead:        true,
	PermissionAssessmentsUpdate:      true,
	PermissionAssessmentsDelete:      true,
	PermissionAssessmentsPublish:     true,
	PermissionAssessmentsGrade:       true,
	PermissionAssessmentsAttempt:     true,
	PermissionAssessmentsViewResults: true,
	// Progreso
	PermissionProgressRead:    true,
	PermissionProgressUpdate:  true,
	PermissionProgressReadOwn: true,
	// Estadísticas
	PermissionStatsGlobal: true,
	PermissionStatsSchool: true,
	PermissionStatsUnit:   true,
	// Screen templates
	PermissionScreenTemplatesRead:   true,
	PermissionScreenTemplatesCreate: true,
	PermissionScreenTemplatesUpdate: true,
	PermissionScreenTemplatesDelete: true,
	// Screen instances
	PermissionScreenInstancesRead:   true,
	PermissionScreenInstancesCreate: true,
	PermissionScreenInstancesUpdate: true,
	PermissionScreenInstancesDelete: true,
	// Screens
	PermissionScreensRead: true,
	// Roles
	PermissionRolesCreate: true,
	PermissionRolesRead:   true,
	PermissionRolesUpdate: true,
	PermissionRolesDelete: true,
	// Gestión de permisos
	PermissionPermissionsMgmtRead:   true,
	PermissionPermissionsMgmtCreate: true,
	PermissionPermissionsMgmtUpdate: true,
	PermissionPermissionsMgmtDelete: true,
	// Membresías
	PermissionMembershipsCreate: true,
	PermissionMembershipsRead:   true,
	PermissionMembershipsUpdate: true,
	PermissionMembershipsDelete: true,
	// Materias
	PermissionSubjectsCreate: true,
	PermissionSubjectsRead:   true,
	PermissionSubjectsUpdate: true,
	PermissionSubjectsDelete: true,
	// Vínculos guardian-estudiante
	PermissionGuardianRelationsRead:    true,
	PermissionGuardianRelationsApprove: true,
	PermissionGuardianRelationsRequest: true,
	PermissionGuardianRelationsManage:  true,
	// Evaluaciones adicionales
	PermissionAssessmentsTake:    true,
	PermissionAssessmentsArchive: true,
	// Tipos de concepto
	PermissionConceptTypesCreate: true,
	PermissionConceptTypesRead:   true,
	PermissionConceptTypesUpdate: true,
	PermissionConceptTypesDelete: true,
	// Auditoría
	PermissionAuditRead:   true,
	PermissionAuditExport: true,
}

// AllPermissionsSlice retorna todos los permisos como slice
func AllPermissionsSlice() []Permission {
	perms := make([]Permission, 0, len(AllPermissions))
	for perm := range AllPermissions {
		perms = append(perms, perm)
	}
	return perms
}
