// Package postgres proporciona una implementación de audit.AuditLogger para PostgreSQL.
//
// Persiste eventos de auditoría en la tabla audit.events usando GORM.
//
// Ejemplo de uso:
//
//	db := gorm.Open(postgres.Open(dsn), &gorm.Config{})
//	logger := postgres.NewPostgresAuditLogger(db, "mi-servicio")
//
//	err := logger.Log(ctx, audit.AuditEvent{
//		Action:       "create",
//		ResourceType: "usuario",
//		ResourceID:   "123",
//		ActorID:      "user-456",
//	})
package postgres
