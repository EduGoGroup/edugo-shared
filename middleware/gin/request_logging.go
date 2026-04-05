package gin

import (
	"log/slog"
	"time"

	"github.com/EduGoGroup/edugo-shared/auth"
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
//  2. Crea un Logger enriquecido con campos contextuales
//  3. Inyecta el logger en gin.Context y context.Context
//  4. Registra la petición completada con status, duración y bytes
//
// Coloca este middleware ANTES del middleware de autenticación.
// Usa PostAuthLogging() DESPUÉS del middleware JWT para enriquecer
// el logger con user_id, role y school_id.
func RequestLogging(baseLogger logger.Logger) gin.HandlerFunc {
	if baseLogger == nil {
		baseLogger = logger.NewSlogAdapter(slog.Default())
	}
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

		// Usar c.FullPath() con fallback a la URL real para rutas no registradas (404)
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		// Crear logger enriquecido con contexto de la petición
		reqLogger := baseLogger.With(
			logger.FieldRequestID, requestID,
			logger.FieldCorrelationID, correlationID,
			logger.FieldMethod, c.Request.Method,
			logger.FieldPath, path,
			logger.FieldIP, c.ClientIP(),
		)

		// Inyectar logger en gin.Context y context.Context
		setLogger(c, reqLogger)

		// Procesar la petición
		c.Next()

		// Post-petición: leer user_id/role/school_id si fueron inyectados por PostAuthLogging
		// para incluirlos en el log de resumen final
		finalLogger := GetLogger(c)

		// Registrar la petición completada
		duration := time.Since(start)
		status := c.Writer.Status()

		attrs := []any{
			logger.FieldStatusCode, status,
			logger.FieldDuration, duration.Milliseconds(),
			logger.FieldBytes, c.Writer.Size(),
		}

		if len(c.Errors) > 0 {
			attrs = append(attrs, logger.FieldError, c.Errors.String())
		}

		switch {
		case status >= 500:
			finalLogger.Error("request completed", attrs...)
		case status >= 400:
			finalLogger.Warn("request completed", attrs...)
		default:
			finalLogger.Info("request completed", attrs...)
		}
	}
}

// PostAuthLogging crea un middleware que enriquece el logger del contexto
// con user_id, role y school_id después de que el middleware JWT los establece.
//
// Cadena de middleware recomendada:
//
//	Recovery -> RequestLogging -> CORS -> JWT -> PostAuthLogging -> handlers
//
// Sin este middleware, los logs de servicios (via logger.FromContext) no
// contendrán información de autenticación.
func PostAuthLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqLogger := GetLogger(c)

		if userID, exists := c.Get(ContextKeyUserID); exists {
			if uid, ok := userID.(string); ok && uid != "" {
				reqLogger = reqLogger.With(logger.FieldUserID, uid)
			}
		}

		if role, exists := c.Get(ContextKeyRole); exists {
			if r, ok := role.(string); ok && r != "" {
				reqLogger = reqLogger.With(logger.FieldRole, r)
			}
		}

		if claims, exists := c.Get(ContextKeyClaims); exists {
			if jwtClaims, ok := claims.(*auth.Claims); ok && jwtClaims.ActiveContext != nil {
				if jwtClaims.ActiveContext.SchoolID != "" {
					reqLogger = reqLogger.With(logger.FieldSchoolID, jwtClaims.ActiveContext.SchoolID)
				}
			}
		}

		// Re-inyectar el logger enriquecido en ambos contextos
		setLogger(c, reqLogger)

		c.Next()
	}
}

// setLogger almacena el logger en gin.Context y context.Context.
func setLogger(c *gin.Context, l logger.Logger) {
	c.Set(ContextKeySlogLogger, l)
	ctx := logger.NewContext(c.Request.Context(), l)
	c.Request = c.Request.WithContext(ctx)
}

// GetLogger extrae el Logger enriquecido del gin.Context.
// Retorna un adapter de slog.Default() si no se almacenó ningún logger.
func GetLogger(c *gin.Context) logger.Logger {
	if l, exists := c.Get(ContextKeySlogLogger); exists {
		if sl, ok := l.(logger.Logger); ok {
			return sl
		}
	}
	return logger.NewSlogAdapter(slog.Default())
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
