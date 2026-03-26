package gin

import (
	"log/slog"
	"time"

	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// HeaderRequestID es el header HTTP para trazabilidad de peticiones.
	HeaderRequestID = "X-Request-ID"
	// HeaderCorrelationID es el header HTTP para correlación entre servicios.
	HeaderCorrelationID = "X-Correlation-ID"

	// ContextKeySlogLogger es la clave en gin.Context para el slog.Logger enriquecido.
	ContextKeySlogLogger = "slog_logger"
	// ContextKeyRequestID es la clave en gin.Context para el ID de petición.
	ContextKeyRequestID = "request_id"
)

// RequestLogging crea un middleware que:
//  1. Genera o extrae un request_id y correlation_id
//  2. Crea un slog.Logger enriquecido con campos contextuales
//  3. Después de la autenticación JWT, enriquece con user_id y role
//  4. Inyecta el logger en gin.Context y context.Context
//  5. Registra la petición completada con status, duración y bytes
//
// Coloca este middleware ANTES del middleware de autenticación para que
// el request_id esté disponible desde el inicio. El enriquecimiento con
// user_id y role ocurre automáticamente después de que el middleware JWT
// establece las claves en el contexto.
func RequestLogging(baseLogger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generar o extraer el ID de petición
		requestID := c.GetHeader(HeaderRequestID)
		if requestID == "" {
			requestID = uuid.NewString()
		}

		// Propagar el ID de correlación (o usar el ID de petición)
		correlationID := c.GetHeader(HeaderCorrelationID)
		if correlationID == "" {
			correlationID = requestID
		}

		// Establecer headers de respuesta para trazabilidad
		c.Header(HeaderRequestID, requestID)
		c.Header(HeaderCorrelationID, correlationID)

		// Almacenar el ID de petición en gin.Context
		c.Set(ContextKeyRequestID, requestID)

		// Crear logger enriquecido con contexto de la petición
		reqLogger := baseLogger.With(
			slog.String(logger.FieldRequestID, requestID),
			slog.String(logger.FieldCorrelationID, correlationID),
			slog.String(logger.FieldMethod, c.Request.Method),
			slog.String(logger.FieldPath, c.FullPath()),
			slog.String(logger.FieldIP, c.ClientIP()),
		)

		// Almacenar el logger en gin.Context
		c.Set(ContextKeySlogLogger, reqLogger)

		// Almacenar el logger en context.Context para acceso desde la capa de servicio
		ctx := logger.NewContext(c.Request.Context(), reqLogger)
		c.Request = c.Request.WithContext(ctx)

		// Procesar la petición
		c.Next()

		// Post-petición: enriquecer con user_id si está disponible (establecido por JWT middleware)
		if userID, exists := c.Get(ContextKeyUserID); exists {
			if uid, ok := userID.(string); ok {
				reqLogger = reqLogger.With(slog.String(logger.FieldUserID, uid))
			}
		}

		// Post-petición: enriquecer con role si está disponible (establecido por JWT middleware)
		if role, exists := c.Get(ContextKeyRole); exists {
			if r, ok := role.(string); ok {
				reqLogger = reqLogger.With(slog.String("role", r))
			}
		}

		// Registrar la petición completada
		duration := time.Since(start)
		status := c.Writer.Status()

		attrs := []any{
			slog.Int(logger.FieldStatusCode, status),
			slog.Int64(logger.FieldDuration, duration.Milliseconds()),
			slog.Int("bytes", c.Writer.Size()),
		}

		if len(c.Errors) > 0 {
			attrs = append(attrs, slog.String(logger.FieldError, c.Errors.String()))
		}

		switch {
		case status >= 500:
			reqLogger.Error("request completed", attrs...)
		case status >= 400:
			reqLogger.Warn("request completed", attrs...)
		default:
			reqLogger.Info("request completed", attrs...)
		}
	}
}

// GetLogger extrae el slog.Logger enriquecido del gin.Context.
// Retorna slog.Default() si no se almacenó ningún logger.
// Úsalo en handlers para obtener un logger con request_id, user_id, etc.
func GetLogger(c *gin.Context) *slog.Logger {
	if l, exists := c.Get(ContextKeySlogLogger); exists {
		if sl, ok := l.(*slog.Logger); ok {
			return sl
		}
	}
	return slog.Default()
}

// GetRequestID extrae el ID de petición del gin.Context.
func GetRequestID(c *gin.Context) string {
	if id, exists := c.Get(ContextKeyRequestID); exists {
		if s, ok := id.(string); ok {
			return s
		}
	}
	return ""
}
