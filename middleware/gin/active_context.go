package gin

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Context keys para active_context resuelto desde el JWT.
const (
	ContextKeyActiveSchoolID = "active_school_id"
	ContextKeyActiveUnitID   = "active_unit_id"
)

// Errores de contexto activo.
var (
	ErrNoActiveSchool = errors.New("active school not set in request context")
	ErrNoActiveUnit   = errors.New("active academic unit not set in request context")
)

// errorResponse refleja el shape estándar de errores de los handlers.
type noActiveContextError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// RequireActiveSchool exige que el JWT traiga active_context.school_id y lo
// publica en el gin.Context bajo ContextKeyActiveSchoolID. Si el token no lo
// trae, rechaza la request con 412 Precondition Required y código
// NO_ACTIVE_SCHOOL.
//
// Diseño: el school_id NUNCA debe llegar por body/query. Es parte de la huella
// de identidad del usuario autenticado; viene en el JWT y este middleware lo
// extrae para que los handlers lo lean vía MustActiveSchoolID(c).
//
// Pre-requisito: este middleware DEBE montarse después del JWT auth middleware
// (que pone claims en el context). Si no, retorna 412 igualmente — fail-safe.
func RequireActiveSchool() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := GetClaims(c)
		if err != nil || claims == nil || claims.ActiveContext == nil || claims.ActiveContext.SchoolID == "" {
			c.AbortWithStatusJSON(http.StatusPreconditionRequired, noActiveContextError{
				Error:   "precondition_required",
				Message: "active school is required; perform context selection before calling this endpoint",
				Code:    "NO_ACTIVE_SCHOOL",
			})
			return
		}
		c.Set(ContextKeyActiveSchoolID, claims.ActiveContext.SchoolID)
		if claims.ActiveContext.AcademicUnitID != "" {
			c.Set(ContextKeyActiveUnitID, claims.ActiveContext.AcademicUnitID)
		}
		c.Next()
	}
}

// RequireActiveContext exige que el JWT traiga active_context.school_id Y
// active_context.academic_unit_id. Es más estricto que RequireActiveSchool —
// úsalo en endpoints scope=unit que necesitan ambos.
func RequireActiveContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := GetClaims(c)
		if err != nil || claims == nil || claims.ActiveContext == nil {
			c.AbortWithStatusJSON(http.StatusPreconditionRequired, noActiveContextError{
				Error:   "precondition_required",
				Message: "active context is required; perform context selection before calling this endpoint",
				Code:    "NO_ACTIVE_CONTEXT",
			})
			return
		}
		if claims.ActiveContext.SchoolID == "" {
			c.AbortWithStatusJSON(http.StatusPreconditionRequired, noActiveContextError{
				Error:   "precondition_required",
				Message: "active school is required",
				Code:    "NO_ACTIVE_SCHOOL",
			})
			return
		}
		if claims.ActiveContext.AcademicUnitID == "" {
			c.AbortWithStatusJSON(http.StatusPreconditionRequired, noActiveContextError{
				Error:   "precondition_required",
				Message: "active academic unit is required",
				Code:    "NO_ACTIVE_UNIT",
			})
			return
		}
		c.Set(ContextKeyActiveSchoolID, claims.ActiveContext.SchoolID)
		c.Set(ContextKeyActiveUnitID, claims.ActiveContext.AcademicUnitID)
		c.Next()
	}
}

// GetActiveSchoolID lee el school_id activo del gin.Context. Devuelve
// ErrNoActiveSchool si el middleware RequireActiveSchool/RequireActiveContext
// no se montó o no encontró el valor.
func GetActiveSchoolID(c *gin.Context) (string, error) {
	v, exists := c.Get(ContextKeyActiveSchoolID)
	if !exists {
		return "", ErrNoActiveSchool
	}
	s, ok := v.(string)
	if !ok || s == "" {
		return "", ErrNoActiveSchool
	}
	return s, nil
}

// MustActiveSchoolID lee el school_id activo o entra en pánico. Es seguro
// usarlo en handlers protegidos por RequireActiveSchool/RequireActiveContext.
func MustActiveSchoolID(c *gin.Context) string {
	s, err := GetActiveSchoolID(c)
	if err != nil {
		panic(err)
	}
	return s
}

// GetActiveUnitID lee el academic_unit_id activo del gin.Context. Devuelve
// ErrNoActiveUnit si no fue seteado por el middleware (puede pasar si el
// endpoint solo monta RequireActiveSchool y el JWT no traía unit).
func GetActiveUnitID(c *gin.Context) (string, error) {
	v, exists := c.Get(ContextKeyActiveUnitID)
	if !exists {
		return "", ErrNoActiveUnit
	}
	s, ok := v.(string)
	if !ok || s == "" {
		return "", ErrNoActiveUnit
	}
	return s, nil
}

// MustActiveUnitID lee el academic_unit_id activo o entra en pánico. Solo
// úsalo en handlers protegidos por RequireActiveContext (estricto).
func MustActiveUnitID(c *gin.Context) string {
	s, err := GetActiveUnitID(c)
	if err != nil {
		panic(err)
	}
	return s
}
