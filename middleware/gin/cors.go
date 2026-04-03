package gin

import (
	"log"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig configuracion para el middleware CORS.
// Los consumidores mantienen sus propias config structs con env tags
// y convierten a esta struct al invocar CORSMiddleware.
type CORSConfig struct {
	AllowedOrigins string
	AllowedMethods string
	AllowedHeaders string
}

// CORSMiddleware crea un middleware Gin para CORS.
// Bloquea wildcard (*) en entornos que no sean development/local (fail-closed).
// En entornos no-development, si se detecta wildcard se llama log.Fatalf.
func CORSMiddleware(cfg CORSConfig, environment string) gin.HandlerFunc {
	allowedOrigins := parseCSV(cfg.AllowedOrigins)

	hasWildcard := slices.Contains(allowedOrigins, "*")

	// Bloquear wildcard CORS en entornos que no sean development (fail-closed: env vacio se trata como non-development)
	normalizedEnv := strings.ToLower(strings.TrimSpace(environment))

	if hasWildcard && normalizedEnv != "development" && normalizedEnv != "local" {
		envForLog := environment
		if strings.TrimSpace(environment) == "" {
			envForLog = "non-development"
		}
		log.Fatalf("CORS wildcard (*) is not allowed in %s environment. Set CORS_ALLOWED_ORIGINS explicitly.", envForLog)
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		originAllowed := origin != "" && isOriginAllowed(origin, allowedOrigins)
		if originAllowed {
			if hasWildcard {
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Set("Vary", "Origin")
			}
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length,ETag,X-Request-ID,X-Correlation-ID")
		}

		if c.Request.Method == "OPTIONS" {
			if originAllowed {
				c.Writer.Header().Set("Access-Control-Allow-Methods", cfg.AllowedMethods)
				c.Writer.Header().Set("Access-Control-Allow-Headers", cfg.AllowedHeaders)
				c.Writer.Header().Set("Access-Control-Max-Age", "86400")
			}
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func parseCSV(csv string) []string {
	if csv == "" {
		return []string{}
	}
	parts := strings.Split(csv, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if origin == "" {
		return false
	}
	for _, allowed := range allowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	return false
}
