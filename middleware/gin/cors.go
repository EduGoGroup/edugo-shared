package gin

import (
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
// En development/local permite wildcard (*), en otros entornos requiere orígenes explícitos (fail-closed).
func CORSMiddleware(cfg CORSConfig, environment string) gin.HandlerFunc {
	allowedOrigins := parseCSV(cfg.AllowedOrigins)

	hasWildcard := slices.Contains(allowedOrigins, "*")

	// Normalizar ambiente: env vacío se trata como non-development (fail-closed)
	normalizedEnv := strings.ToLower(strings.TrimSpace(environment))
	allowsWildcard := normalizedEnv == "development" || normalizedEnv == "local"

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		originAllowed := origin != "" && isOriginAllowed(origin, allowedOrigins, allowsWildcard)
		if originAllowed {
			// Wildcard solo en development/local; otros entornos reflejan el origin
			if hasWildcard && allowsWildcard {
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				// Reflejar origin para evitar combinación inválida: * + credentials
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				appendVaryHeader(c, "Origin")
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

func isOriginAllowed(origin string, allowedOrigins []string, allowsWildcard bool) bool {
	if origin == "" {
		return false
	}
	for _, allowed := range allowedOrigins {
		switch allowed {
		case "*":
			// Solo permito wildcard en dev/local
			if allowsWildcard {
				return true
			}
		case origin:
			return true
		}
	}
	return false
}

// appendVaryHeader agrega "Origin" al header Vary, evitando sobrescribir valores existentes
func appendVaryHeader(c *gin.Context, value string) {
	existing := c.Writer.Header().Get("Vary")
	if existing == "" {
		c.Writer.Header().Set("Vary", value)
	} else if !strings.Contains(existing, value) {
		c.Writer.Header().Set("Vary", existing+","+value)
	}
}
