# API Contracts - spec-04
## Interfaces Públicas
```go
// logger/logger.go
func Init(config Config) error
func Info(msg string, fields ...interface{})
func Error(msg string, fields ...interface{})

// database/postgres.go
func NewPostgresConnection(config DBConfig) (*gorm.DB, error)

// middleware/auth.go
func JWTMiddleware(secretKey string) gin.HandlerFunc
```
