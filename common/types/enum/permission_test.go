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

func TestAllPermissionsSlice(t *testing.T) {
	t.Run("retorna todos los permisos definidos", func(t *testing.T) {
		perms := AllPermissionsSlice()

		// Verificar que retorna la cantidad correcta de permisos
		assert.Len(t, perms, len(AllPermissions))

		// Verificar que todos los permisos retornados son válidos
		for _, perm := range perms {
			assert.True(t, perm.IsValid(), "permiso %s debería ser válido", perm)
		}
	})

	t.Run("contiene permisos de usuarios", func(t *testing.T) {
		perms := AllPermissionsSlice()
		permsMap := make(map[Permission]bool)
		for _, p := range perms {
			permsMap[p] = true
		}

		assert.True(t, permsMap[PermissionUsersCreate])
		assert.True(t, permsMap[PermissionUsersRead])
		assert.True(t, permsMap[PermissionUsersUpdate])
		assert.True(t, permsMap[PermissionUsersDelete])
	})

	t.Run("contiene permisos de escuelas", func(t *testing.T) {
		perms := AllPermissionsSlice()
		permsMap := make(map[Permission]bool)
		for _, p := range perms {
			permsMap[p] = true
		}

		assert.True(t, permsMap[PermissionSchoolsCreate])
		assert.True(t, permsMap[PermissionSchoolsRead])
		assert.True(t, permsMap[PermissionSchoolsManage])
	})

	t.Run("contiene permisos de materiales", func(t *testing.T) {
		perms := AllPermissionsSlice()
		permsMap := make(map[Permission]bool)
		for _, p := range perms {
			permsMap[p] = true
		}

		assert.True(t, permsMap[PermissionMaterialsCreate])
		assert.True(t, permsMap[PermissionMaterialsPublish])
		assert.True(t, permsMap[PermissionMaterialsDownload])
	})

	t.Run("contiene permisos de evaluaciones", func(t *testing.T) {
		perms := AllPermissionsSlice()
		permsMap := make(map[Permission]bool)
		for _, p := range perms {
			permsMap[p] = true
		}

		assert.True(t, permsMap[PermissionAssessmentsCreate])
		assert.True(t, permsMap[PermissionAssessmentsGrade])
		assert.True(t, permsMap[PermissionAssessmentsAttempt])
	})
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
			PermissionAssessmentsViewResults,
			// Progreso
			PermissionProgressRead,
			PermissionProgressUpdate,
			PermissionProgressReadOwn,
			// Estadísticas
			PermissionStatsGlobal,
			PermissionStatsSchool,
			PermissionStatsUnit,
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
