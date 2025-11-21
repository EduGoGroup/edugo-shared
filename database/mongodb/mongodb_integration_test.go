package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-shared/database/mongodb"
	"github.com/EduGoGroup/edugo-shared/testing/containers"
	"go.mongodb.org/mongo-driver/bson"
)

func TestConnect_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	config := containers.NewConfig().
		WithMongoDB(nil).
		Build()

	manager, err := containers.GetManager(t, config)
	if err != nil {
		t.Fatalf("Error creando manager: %v", err)
	}

	mongoContainer := manager.MongoDB()
	if mongoContainer == nil {
		t.Fatal("MongoDB container no disponible")
	}

	// Test de conexión básica
	ctx := context.Background()
	connStr, err := mongoContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Error obteniendo connection string: %v", err)
	}

	cfg := mongodb.Config{
		URI:      connStr,
		Database: "test_db",
		Timeout:  10 * time.Second,
	}

	client, err := mongodb.Connect(cfg)
	if err != nil {
		t.Fatalf("Error conectando a MongoDB: %v", err)
	}

	// Verificar que client no es nil
	if client == nil {
		t.Error("Client es nil")
	}

	// Test de ping
	if err := client.Ping(ctx, nil); err != nil {
		t.Errorf("Error en ping: %v", err)
	}
}

func TestHealthCheck_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	config := containers.NewConfig().
		WithMongoDB(nil).
		Build()

	manager, err := containers.GetManager(t, config)
	if err != nil {
		t.Fatalf("Error creando manager: %v", err)
	}

	mongoContainer := manager.MongoDB()
	client := mongoContainer.Client()

	if err := mongodb.HealthCheck(client); err != nil {
		t.Errorf("HealthCheck falló: %v", err)
	}
}

func TestBasicOperations_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	config := containers.NewConfig().
		WithMongoDB(nil).
		Build()

	manager, err := containers.GetManager(t, config)
	if err != nil {
		t.Fatalf("Error creando manager: %v", err)
	}

	mongoContainer := manager.MongoDB()
	ctx := context.Background()
	db := mongoContainer.Database()

	// Usar colección única para este test
	collection := db.Collection("test_collection_" + t.Name())

	t.Run("InsertOne", func(t *testing.T) {
		doc := bson.M{"name": "test", "value": 123}
		result, err := collection.InsertOne(ctx, doc)
		if err != nil {
			t.Fatalf("Error insertando documento: %v", err)
		}
		if result.InsertedID == nil {
			t.Error("InsertedID es nil")
		}
	})

	t.Run("FindOne", func(t *testing.T) {
		var result bson.M
		err := collection.FindOne(ctx, bson.M{"name": "test"}).Decode(&result)
		if err != nil {
			t.Fatalf("Error buscando documento: %v", err)
		}
		if result["name"] != "test" {
			t.Errorf("Esperado name=test, obtenido %v", result["name"])
		}
	})

	t.Run("UpdateOne", func(t *testing.T) {
		update := bson.M{"$set": bson.M{"value": 456}}
		result, err := collection.UpdateOne(ctx, bson.M{"name": "test"}, update)
		if err != nil {
			t.Fatalf("Error actualizando documento: %v", err)
		}
		if result.ModifiedCount != 1 {
			t.Errorf("Esperado 1 documento modificado, obtenido %d", result.ModifiedCount)
		}
	})

	t.Run("DeleteOne", func(t *testing.T) {
		result, err := collection.DeleteOne(ctx, bson.M{"name": "test"})
		if err != nil {
			t.Fatalf("Error eliminando documento: %v", err)
		}
		if result.DeletedCount != 1 {
			t.Errorf("Esperado 1 documento eliminado, obtenido %d", result.DeletedCount)
		}
	})

	// Cleanup: drop collection
	if err := collection.Drop(ctx); err != nil {
		t.Errorf("Error eliminando colección: %v", err)
	}
}

func TestGetDatabase_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	config := containers.NewConfig().
		WithMongoDB(nil).
		Build()

	manager, err := containers.GetManager(t, config)
	if err != nil {
		t.Fatalf("Error creando manager: %v", err)
	}

	mongoContainer := manager.MongoDB()
	if mongoContainer == nil {
		t.Fatal("MongoDB container no disponible")
	}

	ctx := context.Background()
	connStr, err := mongoContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Error obteniendo connection string: %v", err)
	}

	cfg := mongodb.Config{
		URI:      connStr,
		Database: "test_db",
		Timeout:  10 * time.Second,
	}

	client, err := mongodb.Connect(cfg)
	if err != nil {
		t.Fatalf("Error conectando: %v", err)
	}
	defer client.Disconnect(ctx)

	t.Run("get database", func(t *testing.T) {
		db := mongodb.GetDatabase(client, "test_db")

		if db == nil {
			t.Error("Expected non-nil database")
		}

		if db.Name() != "test_db" {
			t.Errorf("Expected database name 'test_db', got '%s'", db.Name())
		}
	})

	t.Run("get different databases from same client", func(t *testing.T) {
		db1 := mongodb.GetDatabase(client, "db1")
		db2 := mongodb.GetDatabase(client, "db2")

		if db1 == nil || db2 == nil {
			t.Error("Expected non-nil databases")
		}

		if db1.Name() == db2.Name() {
			t.Error("Expected different database names")
		}
	})
}

func TestClose_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	config := containers.NewConfig().
		WithMongoDB(nil).
		Build()

	manager, err := containers.GetManager(t, config)
	if err != nil {
		t.Fatalf("Error creando manager: %v", err)
	}

	mongoContainer := manager.MongoDB()
	if mongoContainer == nil {
		t.Fatal("MongoDB container no disponible")
	}

	ctx := context.Background()
	connStr, err := mongoContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Error obteniendo connection string: %v", err)
	}

	cfg := mongodb.Config{
		URI:      connStr,
		Database: "test_db",
		Timeout:  10 * time.Second,
	}

	client, err := mongodb.Connect(cfg)
	if err != nil {
		t.Fatalf("Error conectando: %v", err)
	}

	t.Run("close successfully", func(t *testing.T) {
		err := mongodb.Close(client)
		if err != nil {
			t.Errorf("Close failed: %v", err)
		}
	})

	t.Run("close nil client", func(t *testing.T) {
		err := mongodb.Close(nil)
		// Should handle nil gracefully or return error
		_ = err
	})
}
