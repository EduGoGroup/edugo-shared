package gin

import "github.com/gin-gonic/gin"

// Header names para versión y build de la API.
// Se exponen como constantes para que consumidores y tests no hardcodeen strings.
const (
	HeaderAPIVersion = "X-Edugo-Api-Version"
	HeaderAPIBuild   = "X-Edugo-Api-Build"
)

// APIVersionHeader crea un middleware Gin que adjunta la versión semver y el
// commit de build de la API en cada respuesta.
//
//   - X-Edugo-Api-Version: versión semver (desde .github/version.txt, horneada por ldflags).
//   - X-Edugo-Api-Build:   git short sha del build (horneado por ldflags).
//
// En desarrollo local ambos valen "dev" porque las vars no se inyectan vía ldflags.
func APIVersionHeader(version, build string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set(HeaderAPIVersion, version)
		c.Writer.Header().Set(HeaderAPIBuild, build)
		c.Next()
	}
}
