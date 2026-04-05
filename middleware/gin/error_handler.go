package gin

import (
	"fmt"
	"net/http"

	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/gin-gonic/gin"
)

// ErrorResponse es la estructura estandar de respuesta de error HTTP.
// Usada por ErrorHandler y HandleError para serializar errores a JSON.
type ErrorResponse struct {
	Error   string            `json:"error"`
	Code    string            `json:"code"`
	Details map[string]string `json:"details,omitempty"`
}

// ErrorHandler es un middleware que recupera panics y procesa errores de c.Errors.
// Combina panic recovery con procesamiento centralizado de AppError.
func ErrorHandler(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Error("panic recovered",
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
					"panic", r,
				)
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error: "internal server error",
					Code:  "INTERNAL_ERROR",
				})
				c.Abort()
			}
		}()

		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			HandleError(c, err)
		}
	}
}

// HandleError escribe la respuesta HTTP apropiada para un error.
// Puede usarse desde handlers directamente o es invocada por ErrorHandler.
// Usa el logger del contexto (inyectado por RequestLogging) para logs correlacionados.
func HandleError(c *gin.Context, err error) {
	reqLogger := GetLogger(c)

	if appErr, ok := errors.GetAppError(err); ok {
		resp := ErrorResponse{
			Error: appErr.Message,
			Code:  string(appErr.Code),
		}
		if len(appErr.Fields) > 0 {
			details := make(map[string]string, len(appErr.Fields))
			for k, v := range appErr.Fields {
				details[k] = fmt.Sprintf("%v", v)
			}
			resp.Details = details
		}
		reqLogger.Error("request failed",
			"error_code", string(appErr.Code),
			"status", appErr.StatusCode,
			"message", appErr.Message,
			logger.FieldPath, requestPath(c),
			logger.FieldMethod, requestMethod(c),
		)
		c.JSON(appErr.StatusCode, resp)
		return
	}

	reqLogger.Error("unexpected error",
		"status", http.StatusInternalServerError,
		"error_code", "INTERNAL_ERROR",
		"message", err.Error(),
		logger.FieldPath, requestPath(c),
		logger.FieldMethod, requestMethod(c),
	)
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: "internal server error",
		Code:  "INTERNAL_ERROR",
	})
}
