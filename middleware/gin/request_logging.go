package gin

import (
	"log/slog"
	"time"

	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// HeaderRequestID is the HTTP header for request tracing.
	HeaderRequestID = "X-Request-Id"
	// HeaderCorrelationID is the HTTP header for cross-service correlation.
	HeaderCorrelationID = "X-Correlation-Id"

	// ContextKeySlogLogger is the gin context key for the enriched slog.Logger.
	ContextKeySlogLogger = "slog_logger"
	// ContextKeyRequestID is the gin context key for the request ID.
	ContextKeyRequestID = "request_id"
)

// RequestLogging creates a middleware that:
//  1. Generates or extracts a request_id and correlation_id
//  2. Creates an enriched slog.Logger with contextual fields
//  3. After JWT auth, enriches with user_id and role
//  4. Injects the logger into gin.Context and context.Context
//  5. Logs the completed request with status, duration, and bytes
//
// Place this middleware BEFORE auth middleware in the chain so that
// request_id is available from the start. The user_id enrichment
// happens automatically after the JWT middleware sets context keys.
func RequestLogging(baseLogger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generate or extract request ID
		requestID := c.GetHeader(HeaderRequestID)
		if requestID == "" {
			requestID = uuid.NewString()
		}

		// Propagate correlation ID (or use request ID)
		correlationID := c.GetHeader(HeaderCorrelationID)
		if correlationID == "" {
			correlationID = requestID
		}

		// Set response headers for tracing
		c.Header(HeaderRequestID, requestID)
		c.Header(HeaderCorrelationID, correlationID)

		// Store request ID in gin context
		c.Set(ContextKeyRequestID, requestID)

		// Create enriched logger with request context
		reqLogger := baseLogger.With(
			slog.String(logger.FieldRequestID, requestID),
			slog.String(logger.FieldCorrelationID, correlationID),
			slog.String(logger.FieldMethod, c.Request.Method),
			slog.String(logger.FieldPath, c.FullPath()),
			slog.String(logger.FieldIP, c.ClientIP()),
		)

		// Store logger in gin context
		c.Set(ContextKeySlogLogger, reqLogger)

		// Store logger in Go context for service layer access
		ctx := logger.NewContext(c.Request.Context(), reqLogger)
		c.Request = c.Request.WithContext(ctx)

		// Process request
		c.Next()

		// Post-request: enrich with user_id if available (set by JWT middleware)
		if userID, exists := c.Get(ContextKeyUserID); exists {
			if uid, ok := userID.(string); ok {
				reqLogger = reqLogger.With(slog.String(logger.FieldUserID, uid))
			}
		}

		// Log completed request
		duration := time.Since(start)
		status := c.Writer.Status()

		attrs := []any{
			slog.Int(logger.FieldStatusCode, status),
			slog.String(logger.FieldDuration, duration.String()),
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

// GetLogger extracts the enriched slog.Logger from the gin context.
// Returns slog.Default() if no logger was stored.
// Use this in handlers to get a logger with request_id, user_id, etc.
func GetLogger(c *gin.Context) *slog.Logger {
	if l, exists := c.Get(ContextKeySlogLogger); exists {
		if sl, ok := l.(*slog.Logger); ok {
			return sl
		}
	}
	return slog.Default()
}

// GetRequestID extracts the request ID from the gin context.
func GetRequestID(c *gin.Context) string {
	if id, exists := c.Get(ContextKeyRequestID); exists {
		if s, ok := id.(string); ok {
			return s
		}
	}
	return ""
}
