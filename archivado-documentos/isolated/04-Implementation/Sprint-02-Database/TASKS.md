# Tareas Sprint 02
## TASK-02-001: Database Helpers
```go
// shared/database/postgres.go
func NewPostgresConnection(cfg DBConfig) (*gorm.DB, error)
func NewMongoConnection(uri string) (*mongo.Client, error)
```
**Tiempo:** 3h
