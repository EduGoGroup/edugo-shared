# Tareas Sprint 01
## TASK-01-001: Implementar shared/logger
**Estimación:** 3h
```go
// shared/logger/logger.go
package logger

type Config struct {
    Level  string // debug|info|warn|error
    Format string // json|text
}

func Init(cfg Config) error
func Debug(msg string, fields ...interface{})
func Info(msg string, fields ...interface{})
func Warn(msg string, fields ...interface{})
func Error(msg string, fields ...interface{})
```
**Tiempo:** 3h
