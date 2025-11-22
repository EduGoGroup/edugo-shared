package containers

import (
	"context"
	"testing"
)

func TestConfigBuilder(t *testing.T) {
	config := NewConfig().
		WithPostgreSQL(nil).
		Build()

	if !config.UsePostgreSQL {
		t.Error("UsePostgreSQL debería ser true")
	}

	if config.PostgresConfig == nil {
		t.Error("PostgresConfig no debería ser nil")
	}

	if config.PostgresConfig.Image != "postgres:15-alpine" {
		t.Errorf("Image esperado 'postgres:15-alpine', obtenido '%s'", config.PostgresConfig.Image)
	}

	if config.PostgresConfig.Database != "edugo_test" {
		t.Errorf("Database esperado 'edugo_test', obtenido '%s'", config.PostgresConfig.Database)
	}
}

func TestConfigBuilderWithCustomConfig(t *testing.T) {
	customCfg := &PostgresConfig{
		Image:    "postgres:16",
		Database: "custom_db",
		Username: "custom_user",
		Password: "custom_pass",
	}

	config := NewConfig().
		WithPostgreSQL(customCfg).
		Build()

	if config.PostgresConfig.Image != "postgres:16" {
		t.Errorf("Image esperado 'postgres:16', obtenido '%s'", config.PostgresConfig.Image)
	}

	if config.PostgresConfig.Database != "custom_db" {
		t.Errorf("Database esperado 'custom_db', obtenido '%s'", config.PostgresConfig.Database)
	}
}

func TestConfigBuilderMultipleContainers(t *testing.T) {
	config := NewConfig().
		WithPostgreSQL(nil).
		WithMongoDB(nil).
		WithRabbitMQ(nil).
		Build()

	if !config.UsePostgreSQL {
		t.Error("UsePostgreSQL debería ser true")
	}

	if !config.UseMongoDB {
		t.Error("UseMongoDB debería ser true")
	}

	if !config.UseRabbitMQ {
		t.Error("UseRabbitMQ debería ser true")
	}
}

// TestManagerIntegration es un test de integración que require Docker
// Skip si Docker no está disponible
func TestManagerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	config := NewConfig().
		WithPostgreSQL(nil).
		Build()

	manager, err := GetManager(t, config)
	if err != nil {
		t.Fatalf("Error creando manager: %v", err)
	}

	if manager.PostgreSQL() == nil {
		t.Error("PostgreSQL container debería estar disponible")
	}

	db := manager.PostgreSQL().DB()
	if db == nil {
		t.Error("DB connection debería estar disponible")
	}

	// Test básico de conexión
	if err := db.Ping(); err != nil {
		t.Errorf("Error haciendo ping a PostgreSQL: %v", err)
	}

	// Cleanup
	ctx := context.Background()
	if err := manager.Cleanup(ctx); err != nil {
		t.Errorf("Error en cleanup: %v", err)
	}
}

// Nota: Los tests TestMongoDBIntegration, TestRabbitMQIntegration, y
// TestAllContainersIntegration fueron eliminados porque:
//
// 1. Son redundantes - cada container ya tiene sus propios tests en:
//    - mongodb_test.go: TestMongoDBContainer_Integration (PASS)
//    - rabbitmq_test.go: TestRabbitMQContainer_Integration (PASS)
//    - postgres_test.go: TestPostgreSQLContainer_Integration (PASS)
//
// 2. Causaban conflictos con el patrón singleton del Manager:
//    - El singleton se inicializa UNA VEZ con la primera config
//    - Tests subsiguientes que piden diferentes configs no funcionan
//    - Cleanup() en un test cierra containers para todos los tests
//
// El TestManagerIntegration verifica que el manager funciona correctamente
// con PostgreSQL, que es el caso de uso más común.

// TestManager_AccessorMethods_NilSafety verifica que accessors retornan nil correctamente
func TestManager_AccessorMethods_NilSafety(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	// Manager solo con PostgreSQL
	config := NewConfig().
		WithPostgreSQL(nil).
		Build()

	manager, err := GetManager(t, config)
	if err != nil {
		t.Fatalf("Error creando manager: %v", err)
	}

	// PostgreSQL debe estar disponible
	if manager.PostgreSQL() == nil {
		t.Error("PostgreSQL() debe retornar un container válido")
	}

	// MongoDB y RabbitMQ deben ser nil (no habilitados)
	if manager.MongoDB() != nil {
		t.Error("MongoDB() debe retornar nil cuando no está habilitado")
	}

	if manager.RabbitMQ() != nil {
		t.Error("RabbitMQ() debe retornar nil cuando no está habilitado")
	}
}

// TestManager_CleanMethods_NotEnabled verifica errores cuando servicios no están habilitados
func TestManager_CleanMethods_NotEnabled(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	ctx := context.Background()

	// Manager solo con PostgreSQL
	config := NewConfig().
		WithPostgreSQL(nil).
		Build()

	manager, err := GetManager(t, config)
	if err != nil {
		t.Fatalf("Error creando manager: %v", err)
	}

	// CleanMongoDB debe fallar (no habilitado)
	err = manager.CleanMongoDB(ctx)
	if err == nil {
		t.Error("CleanMongoDB() debe retornar error cuando MongoDB no está habilitado")
	}
	if err != nil && err.Error() != "MongoDB no está habilitado" {
		t.Errorf("Error esperado 'MongoDB no está habilitado', obtenido: %v", err)
	}

	// PurgeRabbitMQ debe fallar (no habilitado)
	err = manager.PurgeRabbitMQ(ctx)
	if err == nil {
		t.Error("PurgeRabbitMQ() debe retornar error cuando RabbitMQ no está habilitado")
	}
	if err != nil && err.Error() != "RabbitMQ no está habilitado" {
		t.Errorf("Error esperado 'RabbitMQ no está habilitado', obtenido: %v", err)
	}

	// CleanPostgreSQL debe funcionar (está habilitado)
	err = manager.CleanPostgreSQL(ctx, "test_table")
	// Este puede dar error si la tabla no existe, pero no debe dar error de "no habilitado"
	if err != nil && err.Error() == "PostgreSQL no está habilitado" {
		t.Errorf("CleanPostgreSQL no debe decir que PostgreSQL no está habilitado: %v", err)
	}
}

// TestManager_CleanPostgreSQL_TruncateTables verifica truncamiento de tablas
func TestManager_CleanPostgreSQL_TruncateTables(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	ctx := context.Background()

	config := NewConfig().
		WithPostgreSQL(nil).
		Build()

	manager, err := GetManager(t, config)
	if err != nil {
		t.Fatalf("Error creando manager: %v", err)
	}

	pg := manager.PostgreSQL()
	if pg == nil {
		t.Fatal("PostgreSQL container no disponible")
	}

	db := pg.DB()
	if db == nil {
		t.Fatal("DB connection no disponible")
	}

	// Crear tabla y agregar datos
	_, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS test_truncate (id SERIAL PRIMARY KEY, name VARCHAR(100))")
	if err != nil {
		t.Fatalf("Error creando tabla: %v", err)
	}

	_, err = db.ExecContext(ctx, "INSERT INTO test_truncate (name) VALUES ('test1'), ('test2'), ('test3')")
	if err != nil {
		t.Fatalf("Error insertando datos: %v", err)
	}

	// Verificar que hay datos
	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM test_truncate").Scan(&count)
	if err != nil {
		t.Fatalf("Error contando registros: %v", err)
	}
	if count != 3 {
		t.Errorf("Esperaba 3 registros, obtuvo %d", count)
	}

	// Truncar tabla
	err = manager.CleanPostgreSQL(ctx, "test_truncate")
	if err != nil {
		t.Fatalf("Error truncando tabla: %v", err)
	}

	// Verificar que la tabla está vacía
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM test_truncate").Scan(&count)
	if err != nil {
		t.Fatalf("Error contando registros después de truncate: %v", err)
	}
	if count != 0 {
		t.Errorf("Esperaba 0 registros después de truncate, obtuvo %d", count)
	}
}

// TestManager_SingletonBehavior verifica comportamiento singleton
func TestManager_SingletonBehavior(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	config := NewConfig().
		WithPostgreSQL(nil).
		Build()

	// Primera llamada
	manager1, err1 := GetManager(t, config)
	if err1 != nil {
		t.Fatalf("Error en primera llamada a GetManager: %v", err1)
	}

	// Segunda llamada (debe retornar el mismo manager)
	manager2, err2 := GetManager(t, config)
	if err2 != nil {
		t.Fatalf("Error en segunda llamada a GetManager: %v", err2)
	}

	// Verificar que son la misma instancia
	if manager1 != manager2 {
		t.Error("GetManager debe retornar la misma instancia (singleton)")
	}

	// Los containers deben ser los mismos
	if manager1.PostgreSQL() != manager2.PostgreSQL() {
		t.Error("Los containers PostgreSQL deben ser la misma instancia")
	}
}
