package containers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq" // Driver PostgreSQL
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// PostgresContainer envuelve el container de PostgreSQL
// PostgresContainer envuelve el container de PostgreSQL de testcontainers.
// Proporciona acceso directo a la conexión de base de datos y métodos
// de utilidad para truncar tablas y ejecutar scripts SQL.
type PostgresContainer struct {
	container *postgres.PostgresContainer
	db        *sql.DB
	config    *PostgresConfig
}

// createPostgres crea y configura un container de PostgreSQL
func createPostgres(ctx context.Context, cfg *PostgresConfig) (*PostgresContainer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("PostgresConfig no puede ser nil")
	}

	// Crear container con configuración
	container, err := postgres.Run(ctx,
		cfg.Image,
		postgres.WithDatabase(cfg.Database),
		postgres.WithUsername(cfg.Username),
		postgres.WithPassword(cfg.Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error creando container PostgreSQL: %w", err)
	}

	// Obtener connection string
	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("error obteniendo connection string: %w", err)
	}

	// Conectar con retry
	db, err := connectWithRetry(connStr, 10, 2*time.Second)
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("error conectando a PostgreSQL: %w", err)
	}

	pc := &PostgresContainer{
		container: container,
		db:        db,
		config:    cfg,
	}

	// Ejecutar scripts de inicialización si existen
	if len(cfg.InitScripts) > 0 {
		for _, script := range cfg.InitScripts {
			if err := pc.ExecScript(ctx, script); err != nil {
				pc.Terminate(ctx)
				return nil, fmt.Errorf("error ejecutando script %s: %w", script, err)
			}
		}
	}

	return pc, nil
}

// ConnectionString retorna el connection string del container
func (pc *PostgresContainer) ConnectionString(ctx context.Context) (string, error) {
	return pc.container.ConnectionString(ctx, "sslmode=disable")
}

// DB retorna la conexión de base de datos
func (pc *PostgresContainer) DB() *sql.DB {
	return pc.db
}

// ExecScript ejecuta un archivo SQL en la base de datos
func (pc *PostgresContainer) ExecScript(ctx context.Context, scriptPath string) error {
	return ExecSQLFile(ctx, pc.db, scriptPath)
}

// Truncate trunca las tablas especificadas
func (pc *PostgresContainer) Truncate(ctx context.Context, tables ...string) error {
	if len(tables) == 0 {
		return nil
	}

	tx, err := pc.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error iniciando transacción: %w", err)
	}
	defer tx.Rollback()

	// Deshabilitar foreign keys temporalmente
	if _, err := tx.ExecContext(ctx, "SET session_replication_role = 'replica'"); err != nil {
		return fmt.Errorf("error deshabilitando foreign keys: %w", err)
	}

	// Truncar cada tabla
	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", pq.QuoteIdentifier(table))
		if _, err := tx.ExecContext(ctx, query); err != nil {
			return fmt.Errorf("error truncando tabla %s: %w", table, err)
		}
	}

	// Rehabilitar foreign keys
	if _, err := tx.ExecContext(ctx, "SET session_replication_role = 'origin'"); err != nil {
		return fmt.Errorf("error rehabilitando foreign keys: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error commiteando transacción: %w", err)
	}

	return nil
}

// Host retorna el hostname del container PostgreSQL.
//
// Parámetros:
//   - ctx: contexto para la operación
//
// Retorna el hostname del container o un error si no se puede obtener.
// El hostname es típicamente "localhost" cuando el container está expuesto
// en el host local.
func (pc *PostgresContainer) Host(ctx context.Context) (string, error) {
	return pc.container.Host(ctx)
}

// MappedPort retorna el puerto mapeado del container PostgreSQL.
//
// Parámetros:
//   - ctx: contexto para la operación
//
// Retorna el número de puerto del host que está mapeado al puerto 5432
// del container, o un error si no se puede obtener. Este puerto se asigna
// dinámicamente cuando el container inicia y es necesario para establecer
// conexiones desde el host.
func (pc *PostgresContainer) MappedPort(ctx context.Context) (int, error) {
	port, err := pc.container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return 0, err
	}
	return port.Int(), nil
}

// Username retorna el nombre de usuario configurado para PostgreSQL.
//
// Retorna el nombre de usuario que se utilizó al crear el container.
// Este valor proviene de la configuración PostgresConfig y se usa
// para autenticación en las conexiones a la base de datos.
func (pc *PostgresContainer) Username() string {
	return pc.config.Username
}

// Password retorna la contraseña configurada para PostgreSQL.
//
// Retorna la contraseña que se utilizó al crear el container.
// Este valor proviene de la configuración PostgresConfig y se usa
// para autenticación en las conexiones a la base de datos.
func (pc *PostgresContainer) Password() string {
	return pc.config.Password
}

// Database retorna el nombre de la base de datos configurada.
//
// Retorna el nombre de la base de datos que se creó automáticamente
// al iniciar el container. Este valor proviene de la configuración
// PostgresConfig y es la base de datos por defecto para las conexiones.
func (pc *PostgresContainer) Database() string {
	return pc.config.Database
}

// Terminate termina el container y cierra las conexiones
func (pc *PostgresContainer) Terminate(ctx context.Context) error {
	if pc.db != nil {
		pc.db.Close()
	}
	if pc.container != nil {
		return pc.container.Terminate(ctx)
	}
	return nil
}

// connectWithRetry intenta conectar a la base de datos con reintentos
func connectWithRetry(connStr string, maxRetries int, delay time.Duration) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			time.Sleep(delay)
			continue
		}

		if err = db.Ping(); err != nil {
			db.Close()
			time.Sleep(delay)
			continue
		}

		// Configurar pool de conexiones
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(time.Hour)

		return db, nil
	}

	return nil, fmt.Errorf("no se pudo conectar después de %d intentos: %w", maxRetries, err)
}
