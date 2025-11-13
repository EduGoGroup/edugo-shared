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

// TestMongoDBIntegration es un test de integración para MongoDB
func TestMongoDBIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	config := NewConfig().
		WithMongoDB(nil).
		Build()

	manager, err := GetManager(t, config)
	if err != nil {
		t.Fatalf("Error creando manager: %v", err)
	}

	if manager.MongoDB() == nil {
		t.Error("MongoDB container debería estar disponible")
	}

	client := manager.MongoDB().Client()
	if client == nil {
		t.Error("Client connection debería estar disponible")
	}

	// Test básico de conexión
	ctx := context.Background()
	if err := client.Ping(ctx, nil); err != nil {
		t.Errorf("Error haciendo ping a MongoDB: %v", err)
	}

	// Cleanup
	if err := manager.Cleanup(ctx); err != nil {
		t.Errorf("Error en cleanup: %v", err)
	}
}

// TestRabbitMQIntegration es un test de integración para RabbitMQ
func TestRabbitMQIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	config := NewConfig().
		WithRabbitMQ(nil).
		Build()

	manager, err := GetManager(t, config)
	if err != nil {
		t.Fatalf("Error creando manager: %v", err)
	}

	if manager.RabbitMQ() == nil {
		t.Error("RabbitMQ container debería estar disponible")
	}

	conn := manager.RabbitMQ().Connection()
	if conn == nil {
		t.Error("Connection debería estar disponible")
	}

	// Test básico de canal
	ch, err := manager.RabbitMQ().Channel()
	if err != nil {
		t.Errorf("Error creando canal: %v", err)
	}
	if ch != nil {
		ch.Close()
	}

	// Cleanup
	ctx := context.Background()
	if err := manager.Cleanup(ctx); err != nil {
		t.Errorf("Error en cleanup: %v", err)
	}
}

// TestAllContainersIntegration es un test de integración con todos los containers
func TestAllContainersIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	config := NewConfig().
		WithPostgreSQL(nil).
		WithMongoDB(nil).
		WithRabbitMQ(nil).
		Build()

	manager, err := GetManager(t, config)
	if err != nil {
		t.Fatalf("Error creando manager: %v", err)
	}

	// Verificar que todos los containers estén disponibles
	if manager.PostgreSQL() == nil {
		t.Error("PostgreSQL container debería estar disponible")
	}

	if manager.MongoDB() == nil {
		t.Error("MongoDB container debería estar disponible")
	}

	if manager.RabbitMQ() == nil {
		t.Error("RabbitMQ container debería estar disponible")
	}

	// Test de conexión a cada uno
	ctx := context.Background()

	if err := manager.PostgreSQL().DB().Ping(); err != nil {
		t.Errorf("Error en PostgreSQL ping: %v", err)
	}

	if err := manager.MongoDB().Client().Ping(ctx, nil); err != nil {
		t.Errorf("Error en MongoDB ping: %v", err)
	}

	ch, err := manager.RabbitMQ().Channel()
	if err != nil {
		t.Errorf("Error creando canal RabbitMQ: %v", err)
	}
	if ch != nil {
		ch.Close()
	}

	// Cleanup
	if err := manager.Cleanup(ctx); err != nil {
		t.Errorf("Error en cleanup: %v", err)
	}
}
