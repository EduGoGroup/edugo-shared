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

// Gestión de permisos
const (
	PermissionPermissionsMgmtRead   Permission = "permissions_mgmt:read"
	PermissionPermissionsMgmtUpdate Permission = "permissions_mgmt:update"

	// Deprecated: usar PermissionPermissionsMgmtRead
	PermissionResourcesRead = PermissionPermissionsMgmtRead
	// Deprecated: usar PermissionPermissionsMgmtUpdate
	PermissionResourcesUpdate = PermissionPermissionsMgmtUpdate
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
	// Gestión de permisos
	PermissionPermissionsMgmtRead:   true,
	PermissionPermissionsMgmtUpdate: true,
}

// AllPermissionsSlice retorna todos los permisos como slice
func AllPermissionsSlice() []Permission {
	perms := make([]Permission, 0, len(AllPermissions))
	for perm := range AllPermissions {
		perms = append(perms, perm)
	}
	return perms
}
