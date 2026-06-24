package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPermission_String(t *testing.T) {
	cases := []struct {
		perm Permission
		want string
	}{
		{PermissionUsersCreate, "admin.users.create"},
		{PermissionSchoolsManage, "admin.schools.manage"},
		{PermissionAssessmentsGrade, "content.assessments.grade"},
		{PermissionUsersReadOwn, "admin.users.read:own"},
		{PermissionDashboardView, "dashboard.view"},
	}
	for _, tc := range cases {
		t.Run(string(tc.perm), func(t *testing.T) {
			assert.Equal(t, tc.want, tc.perm.String())
		})
	}
}

func TestPermission_IsValid(t *testing.T) {
	cases := []struct {
		name string
		perm Permission
		want bool
	}{
		{"válido path-based exacto", PermissionUsersCreate, true},
		{"válido sufijo :own", PermissionUsersReadOwn, true},
		{"válido raíz 2 segmentos", PermissionDashboardView, true},
		{"inválido formato legacy", Permission("users:create"), false},
		{"inválido string vacío", Permission(""), false},
		{"inválido string libre", Permission("not-a-valid-permission"), false},
		{"inválido pattern (no exacto)", Permission("admin.*"), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.perm.IsValid())
		})
	}
}

func TestAllPermissions_MapIntegrity(t *testing.T) {
	// Lista de TODAS las constantes declaradas. Cualquier `Permission*`
	// nuevo en permission.go debe agregarse acá (y al map AllPermissions).
	all := []Permission{
		PermissionUsersCreate, PermissionUsersRead, PermissionUsersUpdate,
		PermissionUsersDelete, PermissionUsersReadOwn, PermissionUsersUpdateOwn,
		PermissionUsersGrantsManage,
		PermissionSchoolsCreate, PermissionSchoolsRead, PermissionSchoolsUpdate,
		PermissionSchoolsDelete, PermissionSchoolsManage,
		PermissionRolesCreate, PermissionRolesRead, PermissionRolesUpdate, PermissionRolesDelete,
		PermissionPermissionsMgmtCreate, PermissionPermissionsMgmtRead,
		PermissionPermissionsMgmtUpdate, PermissionPermissionsMgmtDelete,
		PermissionScreenTemplatesCreate, PermissionScreenTemplatesRead,
		PermissionScreenTemplatesUpdate, PermissionScreenTemplatesDelete,
		PermissionScreenInstancesCreate, PermissionScreenInstancesRead,
		PermissionScreenInstancesUpdate, PermissionScreenInstancesDelete,
		PermissionAuditRead, PermissionAuditExport,
		PermissionConceptTypesCreate, PermissionConceptTypesRead,
		PermissionConceptTypesUpdate, PermissionConceptTypesDelete,
		PermissionSystemSettingsSettings, PermissionSystemSettingsRead, PermissionSystemSettingsUpdate,
		PermissionUnitsCreate, PermissionUnitsRead, PermissionUnitsUpdate, PermissionUnitsDelete,
		PermissionMembershipsCreate, PermissionMembershipsRead,
		PermissionMembershipsUpdate, PermissionMembershipsDelete,
		PermissionMyMembershipsReadOwn,
		PermissionMyGradesReadOwn,
		PermissionMyTeachingReadOwn,
		PermissionMyAttendanceReadOwn,
		PermissionMyWardsGradesReadOwn, PermissionMyWardsAttendanceReadOwn,
		PermissionMyWardsAnnouncementsReadOwn, PermissionMyWardsMaterialsReadOwn,
		PermissionMyWardsAssessmentsReadOwn,
		PermissionSubjectsCreate, PermissionSubjectsRead, PermissionSubjectsUpdate, PermissionSubjectsDelete,
		PermissionSubjectOfferingsCreate, PermissionSubjectOfferingsRead, PermissionSubjectOfferingsUpdate,
		PermissionSubjectOfferingsDelete, PermissionSubjectOfferingsEnroll,
		PermissionGuardianRelationsRead, PermissionGuardianRelationsApprove,
		PermissionGuardianRelationsRequest, PermissionGuardianRelationsManage,
		PermissionInvitationsCreate, PermissionInvitationsRead, PermissionInvitationsRevoke,
		PermissionJoinRequestsRead, PermissionJoinRequestsReject,
		PermissionPeriodsCreate, PermissionPeriodsRead, PermissionPeriodsUpdate,
		PermissionPeriodsDelete, PermissionPeriodsActivate,
		PermissionGradesCreate, PermissionGradesRead, PermissionGradesUpdate, PermissionGradesFinalize,
		PermissionAttendanceCreate, PermissionAttendanceRead, PermissionAttendanceUpdate,
		PermissionAnnouncementsCreate, PermissionAnnouncementsRead,
		PermissionAnnouncementsUpdate, PermissionAnnouncementsDelete,
		PermissionMaterialsCreate, PermissionMaterialsRead, PermissionMaterialsUpdate,
		PermissionMaterialsDelete, PermissionMaterialsPublish,
		PermissionMaterialsDownload, PermissionMaterialsUpload,
		PermissionAssessmentsCreate, PermissionAssessmentsRead, PermissionAssessmentsUpdate,
		PermissionAssessmentsDelete, PermissionAssessmentsPublish, PermissionAssessmentsGrade,
		PermissionAssessmentsAttempt, PermissionAssessmentsViewResults,
		PermissionAssessmentsAssign, PermissionAssessmentsReview,
		PermissionAssessmentsStudentRead,
		PermissionProgressRead, PermissionProgressUpdate, PermissionProgressReadOwn,
		PermissionStatsGlobal, PermissionStatsSchool, PermissionStatsUnit,
		PermissionDashboardView, PermissionMenuRead, PermissionMenuFullRead,
		PermissionNotificationsRead, PermissionScreensRead,
		PermissionContextBrowseSchools, PermissionContextBrowseUnits,
		PermissionReportsRead,
	}

	t.Run("toda constante está en AllPermissions", func(t *testing.T) {
		for _, perm := range all {
			assert.Truef(t, AllPermissions[perm], "permiso %s no está en AllPermissions", perm)
		}
		assert.Equal(t, len(all), len(AllPermissions),
			"AllPermissions debería tener exactamente %d permisos", len(all))
	})

	t.Run("todo permiso del map cumple IsPathFormat", func(t *testing.T) {
		for perm := range AllPermissions {
			assert.Truef(t, IsPathFormat(string(perm)),
				"permiso %s no cumple PathPermissionRegex", perm)
		}
	})

	t.Run("no hay duplicados en AllPermissions", func(t *testing.T) {
		seen := make(map[string]bool)
		for perm := range AllPermissions {
			permStr := perm.String()
			assert.Falsef(t, seen[permStr], "permiso duplicado: %s", permStr)
			seen[permStr] = true
		}
	})
}

// TestPathPermissionRegex_WildcardExtension cubre la extensión
// wildcard-first del regex: `*.suffix` y `prefix.*.suffix`, además
// de garantizar que los patterns inválidos siguen rechazándose.
func TestPathPermissionRegex_WildcardExtension(t *testing.T) {
	validCases := []string{
		// Soportados previamente — paridad histórica.
		"*",
		"academic.units.read",
		"academic.*",
		"admin.users.*",
		"admin.users.read:own",
		// Extensión wildcard-first.
		"*.create",
		"*.delete",
		"*.update:own",
		"academic.*.create",
		"admin.*.delete",
		"academic.*.create:own",
	}
	for _, p := range validCases {
		t.Run("válido "+p, func(t *testing.T) {
			assert.Truef(t, IsPathFormat(p),
				"pattern %q debería ser válido", p)
		})
	}

	invalidCases := []string{
		// No se aceptan wildcards arbitrarios en cualquier posición.
		"*.*",
		"*.*.create",
		"academic.*.*",
		"academic.*.*.create",
		"a.b.*.c.d", // demasiados segmentos antes y después
		// Sintaxis general inválida.
		"Academic.units.create", // mayúsculas
		"academic..create",      // doble punto
		".create",               // punto inicial sin segmento
		"academic.",             // punto final sin segmento
		"academic-units.create", // guion medio
	}
	for _, p := range invalidCases {
		t.Run("inválido "+p, func(t *testing.T) {
			assert.Falsef(t, IsPathFormat(p),
				"pattern %q NO debería ser válido", p)
		})
	}
}
