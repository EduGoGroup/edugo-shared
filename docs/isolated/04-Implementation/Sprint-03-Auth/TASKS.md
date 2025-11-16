# Tareas Sprint 03
## TASK-03-001: JWT Middleware
```go
// shared/middleware/jwt.go
func JWTMiddleware(secretKey string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Validar JWT
        // Extraer claims
        // Set en context
    }
}
```
**Tiempo:** 3h
