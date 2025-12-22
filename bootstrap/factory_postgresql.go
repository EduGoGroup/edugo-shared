package bootstrap

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// =============================================================================
// POSTGRESQL FACTORY IMPLEMENTATION
// =============================================================================

// DefaultPostgreSQLFactory implementa PostgreSQLFactory usando GORM
type DefaultPostgreSQLFactory struct {
	logger logger.Interface
}

// NewDefaultPostgreSQLFactory crea una nueva instancia de DefaultPostgreSQLFactory
func NewDefaultPostgreSQLFactory(gormLogger logger.Interface) *DefaultPostgreSQLFactory {
	if gormLogger == nil {
		gormLogger = logger.Default.LogMode(logger.Info)
	}
	return &DefaultPostgreSQLFactory{
		logger: gormLogger,
	}
}

// CreateConnection crea una conexión GORM a PostgreSQL
func (f *DefaultPostgreSQLFactory) CreateConnection(ctx context.Context, config PostgreSQLConfig) (*gorm.DB, error) {
	dsn := f.buildDSN(config)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: f.logger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Configurar connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from GORM: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	// Verificar conexión
	if err := f.Ping(ctx, db); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	return db, nil
}

// CreateRawConnection crea una conexión SQL nativa
func (f *DefaultPostgreSQLFactory) CreateRawConnection(ctx context.Context, config PostgreSQLConfig) (*sql.DB, error) {
	dsn := f.buildDSN(config)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open raw PostgreSQL connection: %w", err)
	}

	// Configurar connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	// Verificar conexión
	if err := db.PingContext(ctx); err != nil {
		if closeErr := db.Close(); closeErr != nil {
			return nil, errors.Join(
				fmt.Errorf("failed to ping PostgreSQL: %w", err),
				fmt.Errorf("failed to close connection: %w", closeErr),
			)
		}
		return nil, fmt.Errorf("failed to ping raw PostgreSQL: %w", err)
	}

	return db, nil
}

// Ping verifica la conectividad con PostgreSQL
func (f *DefaultPostgreSQLFactory) Ping(ctx context.Context, db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	return nil
}

// Close cierra la conexión
func (f *DefaultPostgreSQLFactory) Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}

	return nil
}

// buildDSN construye el Data Source Name para PostgreSQL
func (f *DefaultPostgreSQLFactory) buildDSN(config PostgreSQLConfig) string {
	sslMode := config.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Database,
		sslMode,
	)
}

// Verificar que DefaultPostgreSQLFactory implementa PostgreSQLFactory
var _ PostgreSQLFactory = (*DefaultPostgreSQLFactory)(nil)
