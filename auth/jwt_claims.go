package auth

import "github.com/golang-jwt/jwt/v5"

// UserContext representa el contexto activo del usuario en el JWT.
// Encapsula el rol actual, la escuela/unidad académica asociadas y los
// grants efectivos (patterns allow + deny) para evaluación stateless
// del middleware.
//
// Nota sobre D2: el spec original pedía que los grants viajaran solo
// en el HTTP body. Acá los llevamos también en el JWT por pragmatismo:
// el middleware de las APIs evalúa cada request contra estos patterns
// usando el matcher con glob — meterlos en el JWT evita ~1 query DB
// por request. El cliente sigue leyendo grants del body (es la fuente
// para él), pero el backend usa el claim.
type UserContext struct {
	RoleID           string `json:"role_id"`
	RoleName         string `json:"role_name"`
	SchoolID         string `json:"school_id,omitempty"`
	SchoolName       string `json:"school_name,omitempty"`
	AcademicUnitID   string `json:"academic_unit_id,omitempty"`
	AcademicUnitName string `json:"academic_unit_name,omitempty"`
	Grants           Grants `json:"grants"`
}

// Grants es el wire format D2: lista de patterns allow + lista de
// patterns deny. Idéntica en estructura a domain.Grants (cada módulo
// la declara para evitar dependencia cruzada hacia edugo-api-identity).
type Grants struct {
	Allow []string `json:"allow"`
	Deny  []string `json:"deny"`
}

// Claims representa los claims personalizados del JWT.
//
// El access_token usa `ActiveContext` (estructura completa con nombres y grants,
// que el middleware lee en cada request sin pegarle a la BD).
//
// El refresh_token (`token_use = "refresh"`) usa `ActiveContext = nil` y solo
// lleva el snapshot mínimo del contexto activo en `SchoolID`, `AcademicUnitID`
// y `RoleID`. Eso le alcanza al refresh use case para re-emitir el siguiente
// access_token con los mismos claims sin reconstruir nada arbitrariamente.
// Los grants y nombres NO se guardan en el refresh — pueden cambiar entre
// rotaciones y el access_token nuevo los recompone consultando IAM.
type Claims struct {
	UserID         string       `json:"user_id"`
	Email          string       `json:"email"`
	ActiveContext  *UserContext `json:"active_context"`
	TokenUse       string       `json:"token_use,omitempty"`
	SchoolID       string       `json:"school_id,omitempty"`
	AcademicUnitID string       `json:"academic_unit_id,omitempty"`
	RoleID         string       `json:"role_id,omitempty"`
	jwt.RegisteredClaims
}
