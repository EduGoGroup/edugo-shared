package containers

import (
	"context"
	"testing"
)

// TestPostgresContainer_Integration tests completos de PostgreSQL
func TestPostgresContainer_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Omitiendo test de integración en modo short")
	}

	ctx := context.Background()
	cfg := &PostgresConfig{
		Image:    "postgres:15-alpine",
		Database: "test_db",
		Username: "test_user",
		Password: "test_pass",
		Port:     "5432",
	}

	container, err := createPostgres(ctx, cfg)
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

	t.Run("DB_Connection", func(t *testing.T) {
		db := container.DB()
		if db == nil {
			t.Fatal("DB no debería ser nil")
		}

		if err := db.Ping(); err != nil {
			t.Errorf("Error haciendo ping: %v", err)
		}
	})

	t.Run("CreateTable_And_Truncate", func(t *testing.T) {
		db := container.DB()

		// Crear tabla de prueba
		_, err := db.ExecContext(ctx, `
			CREATE TABLE test_users (
				id SERIAL PRIMARY KEY,
				name VARCHAR(100)
			)
		`)
		if err != nil {
			t.Fatalf("Error creando tabla: %v", err)
		}

		// Insertar datos
		_, err = db.ExecContext(ctx, `INSERT INTO test_users (name) VALUES ('Alice'), ('Bob')`)
		if err != nil {
			t.Fatalf("Error insertando datos: %v", err)
		}

		// Verificar datos
		var count int
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM test_users`).Scan(&count)
		if err != nil {
			t.Fatalf("Error contando registros: %v", err)
		}
		if count != 2 {
			t.Errorf("Esperado 2 registros, obtenido %d", count)
		}

		// Truncar tabla
		err = container.Truncate(ctx, "test_users")
		if err != nil {
			t.Errorf("Error truncando tabla: %v", err)
		}

		// Verificar que está vacía
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM test_users`).Scan(&count)
		if err != nil {
			t.Fatalf("Error contando después de truncate: %v", err)
		}
		if count != 0 {
			t.Errorf("Esperado 0 registros después de truncate, obtenido %d", count)
		}
	})

	t.Run("Truncate_MultipleTablesWithForeignKeys", func(t *testing.T) {
		db := container.DB()

		// Crear tablas con foreign keys
		_, err := db.ExecContext(ctx, `
			CREATE TABLE test_authors (
				id SERIAL PRIMARY KEY,
				name VARCHAR(100)
			);
			
			CREATE TABLE test_books (
				id SERIAL PRIMARY KEY,
				title VARCHAR(100),
				author_id INTEGER REFERENCES test_authors(id)
			);
		`)
		if err != nil {
			t.Fatalf("Error creando tablas: %v", err)
		}

		// Insertar datos con relaciones
		_, err = db.ExecContext(ctx, `
			INSERT INTO test_authors (id, name) VALUES (1, 'Author1');
			INSERT INTO test_books (title, author_id) VALUES ('Book1', 1);
		`)
		if err != nil {
			t.Fatalf("Error insertando datos: %v", err)
		}

		// Truncar ambas tablas (debería manejar las foreign keys)
		err = container.Truncate(ctx, "test_books", "test_authors")
		if err != nil {
			t.Errorf("Error truncando tablas con FK: %v", err)
		}

		// Verificar que ambas están vacías
		var count int
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM test_authors`).Scan(&count)
		if err != nil {
			t.Fatalf("Error contando registros en test_authors: %v", err)
		}
		if count != 0 {
			t.Errorf("test_authors debería estar vacía, tiene %d registros", count)
		}

		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM test_books`).Scan(&count)
		if err != nil {
			t.Fatalf("Error contando registros en test_books: %v", err)
		}
		if count != 0 {
			t.Errorf("test_books debería estar vacía, tiene %d registros", count)
		}
	})

	t.Run("Truncate_EmptyList", func(t *testing.T) {
		// No debería dar error con lista vacía
		err := container.Truncate(ctx)
		if err != nil {
			t.Errorf("Truncate con lista vacía no debería dar error: %v", err)
		}
	})
}

func TestCreatePostgres_NilConfig(t *testing.T) {
	ctx := context.Background()
	_, err := createPostgres(ctx, nil)
	if err == nil {
		t.Error("createPostgres con config nil debería dar error")
	}
}
