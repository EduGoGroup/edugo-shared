package auth

import "github.com/golang-jwt/jwt/v5"

// UserContext representa el contexto activo del usuario en el JWT.
// Encapsula el rol actual, la escuela y unidad académica asociadas, y los permisos
// específicos del usuario en ese contexto.
//
// Campos opcionales (omitempty):
//   - SchoolID, SchoolName: Solo para usuarios con contexto de escuela
//   - AcademicUnitID, AcademicUnitName: Solo para usuarios con contexto de unidad académica
type UserContext struct {
	RoleID           string   `json:"role_id"`
	RoleName         string   `json:"role_name"`
	SchoolID         string   `json:"school_id,omitempty"`
	SchoolName       string   `json:"school_name,omitempty"`
	AcademicUnitID   string   `json:"academic_unit_id,omitempty"`
	AcademicUnitName string   `json:"academic_unit_name,omitempty"`
	Permissions      []string `json:"permissions"`
}

// Claims representa los claims personalizados del JWT
type Claims struct {
	UserID        string       `json:"user_id"`
	Email         string       `json:"email"`
	ActiveContext *UserContext `json:"active_context"`
	TokenUse      string       `json:"token_use,omitempty"`
	// SchoolID se incluye en refresh tokens para preservar el contexto de escuela
	// seleccionado por el usuario a través de toda la vida útil del token.
	SchoolID string `json:"school_id,omitempty"`
	jwt.RegisteredClaims
}
