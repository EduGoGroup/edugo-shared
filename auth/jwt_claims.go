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
	// Landing es el screen_key de la pantalla de aterrizaje (landing) del
	// contexto activo, resuelto data-driven en el backend (ADR 0024):
	// role.landing_screen_key ?? school.default_landing_screen_key ??
	// "dashboard-home". Viaja en el claim active_context del JWT y en el body
	// HTTP; el cliente lo usa para decidir el destino inicial al activar el
	// contexto.
	Landing string `json:"landing,omitempty"`
	Grants  Grants `json:"grants"`
	// Dimensión "sujeto" + "modo de actor" del representante/guardián (ADR 0026).
	// El guardián que "actúa" sigue siendo el UserID (raíz del token); estos
	// campos marcan al acudido que se está viendo y el modo, para auditoría
	// (ADR 0026 DEC-R-A.1). En contexto propio van omitidos (ActorModeSelf).
	SubjectStudentID   string `json:"subject_student_id,omitempty"`
	SubjectStudentName string `json:"subject_student_name,omitempty"`
	ActorMode          string `json:"actor_mode,omitempty"`
}

// ActorMode describe el modo de actuación del usuario sobre el contexto activo
// (ADR 0026 DEC-R-A.1). ActorModeSelf se OMITE en el JWT por omitempty.
const (
	ActorModeSelf = "self" // contexto propio; se OMITE en el JWT (omitempty)
	ActorModeWard = "ward" // el representante ve a un acudido
)

// Grants es el wire format D2: lista de patterns allow + lista de
// patterns deny. Idéntica en estructura a domain.Grants (cada módulo
// la declara para evitar dependencia cruzada hacia edugo-api-identity).
type Grants struct {
	Allow []string `json:"allow"`
	Deny  []string `json:"deny"`
}

// Claims representa los claims del Context Token ES256.
//
// `ActiveContext` lleva la estructura completa con nombres y grants, que el
// middleware lee en cada request sin pegarle a la BD.
type Claims struct {
	UserID        string       `json:"user_id"`
	Email         string       `json:"email"`
	ActiveContext *UserContext `json:"active_context"`
	jwt.RegisteredClaims
}
