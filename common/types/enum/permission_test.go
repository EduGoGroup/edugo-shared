package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPermission_String(t *testing.T) {
	tests := []struct {
		name       string
		permission Permission
		want       string
	}{
		{
			name:       "users:create retorna string correcto",
			permission: PermissionUsersCreate,
			want:       "users:create",
		},
		{
			name:       "schools:manage retorna string correcto",
			permission: PermissionSchoolsManage,
			want:       "schools:manage",
		},
		{
			name:       "assessments:grade retorna string correcto",
			permission: PermissionAssessmentsGrade,
			want:       "assessments:grade",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.permission.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPermission_IsValid(t *testing.T) {
	tests := []struct {
		name       string
		permission Permission
		want       bool
	}{
		{
			name:       "permiso válido: users:create",
			permission: PermissionUsersCreate,
			want:       true,
		},
		{
			name:       "permiso válido: schools:read",
			permission: PermissionSchoolsRead,
			want:       true,
		},
		{
			name:       "permiso válido: materials:publish",
			permission: PermissionMaterialsPublish,
			want:       true,
		},
		{
			name:       "permiso inválido: permiso no existente",
			permission: Permission("invalid:permission"),
			want:       false,
		},
		{
			name:       "permiso inválido: string vacío",
			permission: Permission(""),
			want:       false,
		},
		{
			name:       "permiso inválido: formato incorrecto",
			permission: Permission("not-a-valid-permission"),
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.permission.IsValid()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAllPermissions_MapIntegrity(t *testing.T) {
	t.Run("todas las constantes están en AllPermissions", func(t *testing.T) {
		// Lista de todas las constantes definidas
		allConstants := []Permission{
			// Usuarios
			PermissionUsersCreate,
			PermissionUsersRead,
			PermissionUsersUpdate,
			PermissionUsersDelete,
			PermissionUsersReadOwn,
			PermissionUsersUpdateOwn,
			// Escuelas
			PermissionSchoolsCreate,
			PermissionSchoolsRead,
			PermissionSchoolsUpdate,
			PermissionSchoolsDelete,
			PermissionSchoolsManage,
			// Unidades
			PermissionUnitsCreate,
			PermissionUnitsRead,
			PermissionUnitsUpdate,
			PermissionUnitsDelete,
			// Materiales
			PermissionMaterialsCreate,
			PermissionMaterialsRead,
			PermissionMaterialsUpdate,
			PermissionMaterialsDelete,
			PermissionMaterialsPublish,
			PermissionMaterialsDownload,
			// Evaluaciones
			PermissionAssessmentsCreate,
			PermissionAssessmentsRead,
			PermissionAssessmentsUpdate,
			PermissionAssessmentsDelete,
			PermissionAssessmentsPublish,
			PermissionAssessmentsGrade,
			PermissionAssessmentsAttempt,
			PermissionAssessmentsAssign,
			PermissionAssessmentsReview,
			PermissionNotificationsRead,
			PermissionMaterialsUpload,

			PermissionAssessmentsViewResults,
			// Progreso
			PermissionProgressRead,
			PermissionProgressUpdate,
			PermissionProgressReadOwn,
			// Estadísticas
			PermissionStatsGlobal,
			PermissionStatsSchool,
			PermissionStatsUnit,
			// Screen templates
			PermissionScreenTemplatesRead,
			PermissionScreenTemplatesCreate,
			PermissionScreenTemplatesUpdate,
			PermissionScreenTemplatesDelete,
			// Screen instances
			PermissionScreenInstancesRead,
			PermissionScreenInstancesCreate,
			PermissionScreenInstancesUpdate,
			PermissionScreenInstancesDelete,
			// Screens
			PermissionScreensRead,
			// Gestión de permisos
			PermissionPermissionsMgmtRead,
			PermissionPermissionsMgmtCreate,
			PermissionPermissionsMgmtUpdate,
			PermissionPermissionsMgmtDelete,
			// Roles
			PermissionRolesCreate,
			PermissionRolesRead,
			PermissionRolesUpdate,
			PermissionRolesDelete,
			// Membresías
			PermissionMembershipsCreate,
			PermissionMembershipsRead,
			PermissionMembershipsUpdate,
			PermissionMembershipsDelete,
			// Materias
			PermissionSubjectsCreate,
			PermissionSubjectsRead,
			PermissionSubjectsUpdate,
			PermissionSubjectsDelete,
			// Vínculos guardian-estudiante
			PermissionGuardianRelationsRead,
			PermissionGuardianRelationsApprove,
			PermissionGuardianRelationsRequest,
			PermissionGuardianRelationsManage,
			// Evaluaciones para estudiantes
			PermissionAssessmentsStudentRead,
			// Dashboard
			PermissionDashboardView,
			// Configuración del sistema
			PermissionSystemSettingsSettings,
			// Auditoría
			PermissionAuditRead,
			PermissionAuditExport,
			// Tipos de concepto
			PermissionConceptTypesCreate,
			PermissionConceptTypesRead,
			PermissionConceptTypesUpdate,
			PermissionConceptTypesDelete,
			// Horarios
			PermissionSchedulesCreate,
			PermissionSchedulesRead,
			PermissionSchedulesUpdate,
			PermissionSchedulesDelete,
			// Anuncios
			PermissionAnnouncementsCreate,
			PermissionAnnouncementsRead,
			PermissionAnnouncementsUpdate,
			PermissionAnnouncementsDelete,
			// Eventos de calendario
			PermissionCalendarEventsCreate,
			PermissionCalendarEventsRead,
			PermissionCalendarEventsUpdate,
			PermissionCalendarEventsDelete,
			// Asistencia
			PermissionAttendanceCreate,
			PermissionAttendanceRead,
			// Períodos académicos
			PermissionPeriodsCreate,
			PermissionPeriodsRead,
			PermissionPeriodsUpdate,
			PermissionPeriodsDelete,
			// Calificaciones
			PermissionGradesCreate,
			PermissionGradesRead,
			PermissionGradesUpdate,
			PermissionGradesFinalize,
			// Reportes
			PermissionReportsRead,
			// Contexto
			PermissionContextBrowseSchools,
			PermissionContextBrowseUnits,
			// Períodos especiales
			PermissionPeriodsActivate,
		}

		// Verificar que cada constante está en AllPermissions
		for _, perm := range allConstants {
			assert.True(t, AllPermissions[perm],
				"permiso %s no está en AllPermissions", perm)
		}

		// Verificar que AllPermissions tiene el mismo tamaño que allConstants
		assert.Equal(t, len(allConstants), len(AllPermissions),
			"AllPermissions debería tener exactamente %d permisos", len(allConstants))
	})

	t.Run("no hay permisos duplicados en AllPermissions", func(t *testing.T) {
		seen := make(map[string]bool)
		for perm := range AllPermissions {
			permStr := perm.String()
			assert.False(t, seen[permStr], "permiso duplicado encontrado: %s", permStr)
			seen[permStr] = true
		}
	})
}
