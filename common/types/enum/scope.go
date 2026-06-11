package enum

// Scope representa un permiso M2M (machine-to-machine) que viaja en el
// service JWT, NO en el JWT de usuario. A diferencia de Permission (RBAC
// path-based del usuario), un Scope autoriza a un cliente de servicio
// (`auth.service_clients`) a invocar un endpoint interno (`/api/v1/internal/*`).
//
// Vive en el mismo paquete que Permission para ser la única fuente de verdad
// del catálogo de autorizaciones del backend (D14/D17 del plan 020).
type Scope string

// notifications — scopes del Notification Gateway (plan 020 N5).
const (
	// ScopeNotificationsDispatch autoriza a un cliente M2M (ej. edugo-worker,
	// edugo-api-learning) a invocar el Notification Gateway:
	// POST /api/v1/internal/notifications/dispatch. Es el scope que valida
	// ServiceJWTAuthMiddleware en platform y el que se siembra en
	// auth.service_clients (D15/D16).
	ScopeNotificationsDispatch Scope = "notifications.dispatch"
)

// String retorna la representación en string del scope.
func (s Scope) String() string {
	return string(s)
}

// IsValid verifica si el scope es uno de los conocidos por el sistema.
func (s Scope) IsValid() bool {
	return AllScopes[s]
}

// AllScopes es el catálogo cerrado de scopes M2M conocidos. Cualquier
// `Scope*` declarado arriba debe aparecer acá; el test
// `TestAllScopes_MapIntegrity` lo verifica.
var AllScopes = map[Scope]bool{
	ScopeNotificationsDispatch: true,
}
