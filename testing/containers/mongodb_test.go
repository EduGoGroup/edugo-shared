package containers

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestMongoDBContainer_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	ctx := context.Background()
	cfg := &MongoConfig{
		Image:    "mongo:7.0",
		Database: "test_db",
	}

	container, err := createMongoDB(ctx, cfg)
	if err != nil {
		t.Fatalf("Error creando container: %v", err)
	}
	defer container.Terminate(ctx)

	t.Run("ConnectionString", func(t *testing.T) {
		connStr, err := container.ConnectionString(ctx)
		if err != nil {
			t.Errorf("Error obteniendo connection string: %v", err)
		}
		if connStr == "" {
			t.Error("Connection string está vacío")
		}
	})

	t.Run("Client_And_Ping", func(t *testing.T) {
		client := container.Client()
		if client == nil {
			t.Fatal("Client no debería ser nil")
		}

		if err := client.Ping(ctx, nil); err != nil {
			t.Errorf("Error haciendo ping: %v", err)
		}
	})

	t.Run("Database_Access", func(t *testing.T) {
		db := container.Database()
		if db == nil {
			t.Fatal("Database no debería ser nil")
		}

		if db.Name() != cfg.Database {
			t.Errorf("Nombre de database esperado %s, obtenido %s", cfg.Database, db.Name())
		}
	})

	t.Run("Insert_And_DropCollections", func(t *testing.T) {
		db := container.Database()
		coll := db.Collection("test_users")

		// Insertar documentos
		docs := []interface{}{
			bson.M{"name": "Alice", "age": 30},
			bson.M{"name": "Bob", "age": 25},
		}
		_, err := coll.InsertMany(ctx, docs)
		if err != nil {
			t.Fatalf("Error insertando documentos: %v", err)
		}

		// Verificar datos
		count, err := coll.CountDocuments(ctx, bson.M{})
		if err != nil {
			t.Fatalf("Error contando documentos: %v", err)
		}
		if count != 2 {
			t.Errorf("Esperado 2 documentos, obtenido %d", count)
		}

		// Eliminar colección específica
		err = container.DropCollections(ctx, "test_users")
		if err != nil {
			t.Errorf("Error eliminando colección: %v", err)
		}

		// Verificar que fue eliminada (crear nueva referencia)
		coll = db.Collection("test_users")
		count, err = coll.CountDocuments(ctx, bson.M{})
		if err != nil {
			t.Fatalf("Error contando después de drop: %v", err)
		}
		if count != 0 {
			t.Errorf("Esperado 0 documentos después de drop, obtenido %d", count)
		}
	})

	t.Run("DropAllCollections", func(t *testing.T) {
		db := container.Database()

		// Crear múltiples colecciones con datos
		collections := []string{"coll1", "coll2", "coll3"}
		for _, collName := range collections {
			coll := db.Collection(collName)
			_, err := coll.InsertOne(ctx, bson.M{"data": "test"})
			if err != nil {
				t.Fatalf("Error insertando en %s: %v", collName, err)
			}
		}

		// Eliminar todas las colecciones
		err := container.DropAllCollections(ctx)
		if err != nil {
			t.Errorf("Error eliminando todas las colecciones: %v", err)
		}

		// Verificar que no hay colecciones
		names, err := db.ListCollectionNames(ctx, bson.M{})
		if err != nil {
			t.Fatalf("Error listando colecciones: %v", err)
		}
		if len(names) != 0 {
			t.Errorf("Esperado 0 colecciones, obtenido %d: %v", len(names), names)
		}
	})
}

func TestCreateMongoDB_NilConfig(t *testing.T) {
	ctx := context.Background()
	_, err := createMongoDB(ctx, nil)
	if err == nil {
		t.Error("createMongoDB con config nil debería dar error")
	}
}
