package auth

import "testing"

// Tests unitarios complementarios al golden vector. El golden cubre la
// matriz completa cross-language; estos casos dejan explícitas las
// reglas de la extensión wildcard-first (*.suffix y prefix.*.suffix)
// con nombres legibles desde Go.
func TestPermissionMatches_LeadingWildcard(t *testing.T) {
	cases := []struct {
		name    string
		pattern string
		request string
		want    bool
	}{
		{"matchea root simple", "*.create", "users.create", true},
		{"matchea path multinivel", "*.create", "academic.units.create", true},
		{"matchea path tres niveles", "*.delete", "admin.academic.units.delete", true},
		{"no matchea sufijo distinto", "*.delete", "users.create", false},
		{"no matchea request con :own", "*.create", "users.create:own", false},
		{"no matchea segmento unico sin punto", "*.create", "create", false},
		{"no matchea sufijo embebido sin punto", "*.create", "usercreate", false},
		{"matchea pattern :own contra request :own", "*.read:own", "users.read:own", true},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if got := PermissionMatches(tc.pattern, tc.request); got != tc.want {
				t.Errorf("PermissionMatches(%q, %q) = %v, want %v",
					tc.pattern, tc.request, got, tc.want)
			}
		})
	}
}

func TestPermissionMatches_MiddleWildcard(t *testing.T) {
	cases := []struct {
		name    string
		pattern string
		request string
		want    bool
	}{
		{"matchea hijo directo", "academic.*.create", "academic.units.create", true},
		{"matchea segmentos intermedios extra", "academic.*.create", "academic.units.subitems.create", true},
		{"no matchea sin segmento intermedio", "academic.*.create", "academic.create", false},
		{"no matchea prefix distinto", "academic.*.create", "admin.units.create", false},
		{"no matchea sufijo distinto", "academic.*.create", "academic.units.delete", false},
		{"no matchea request con :own", "academic.*.create", "academic.units.create:own", false},
		{"matchea admin estrella delete", "admin.*.delete", "admin.users.delete", true},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if got := PermissionMatches(tc.pattern, tc.request); got != tc.want {
				t.Errorf("PermissionMatches(%q, %q) = %v, want %v",
					tc.pattern, tc.request, got, tc.want)
			}
		})
	}
}

// TestPermissionMatches_NoRegresionExisting verifica que la extensión
// wildcard-first no rompe ninguna de las semánticas históricas
// (`*`, literal, `prefix.*`).
func TestPermissionMatches_NoRegresionExisting(t *testing.T) {
	cases := []struct {
		name    string
		pattern string
		request string
		want    bool
	}{
		{"wildcard total", "*", "anything.goes", true},
		{"match exacto", "users.read", "users.read", true},
		{"match exacto fallido", "users.read", "users.write", false},
		{"subtree cubre prefix solo", "users.*", "users", true},
		{"subtree cubre hijo", "users.*", "users.read", true},
		{"subtree cubre nieto", "users.*", "users.read.detail", true},
		{"subtree no cruza root", "users.*", "schools.read", false},
		{"own exacto", "users.read:own", "users.read:own", true},
		{"own no matchea sin own", "users.read:own", "users.read", false},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if got := PermissionMatches(tc.pattern, tc.request); got != tc.want {
				t.Errorf("PermissionMatches(%q, %q) = %v, want %v",
					tc.pattern, tc.request, got, tc.want)
			}
		})
	}
}
