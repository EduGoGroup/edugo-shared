package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-shared/bootstrap"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// Factory implementa la creacion de conexiones PostgreSQL con soporte GORM.
type Factory struct{}

// NewFactory crea una nueva Factory de PostgreSQL.
func NewFactory() *Factory {
	return &Factory{}
}

// CreateRawConnection crea una conexion *sql.DB usando pgx con soporte para
// SimpleProtocol (PgBouncer/Neon) y SearchPath configurable.
func (f *Factory) CreateRawConnection(ctx context.Context, cfg bootstrap.PostgreSQLConfig) (*sql.DB, error) {
	pgxCfg, err := f.buildPgxConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("bootstrap/postgres: parse config: %w", err)
	}

	db := stdlib.OpenDB(*pgxCfg)
	f.applyPoolConfig(db, cfg)

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("bootstrap/postgres: ping: %w", err)
	}

	return db, nil
}

// CreateGORMConnection crea una conexion *gorm.DB configurada.
// Usa pgx con SimpleProtocol por defecto (compatible con PgBouncer/Neon).
func (f *Factory) CreateGORMConnection(
	ctx context.Context,
	cfg bootstrap.PostgreSQLConfig,
	opts ...bootstrap.GORMOption,
) (*gorm.DB, error) {
	options := bootstrap.ApplyGORMOptions(opts...)

	if options.SimpleProtocol && options.PrepareStmt {
		return nil, fmt.Errorf("bootstrap/postgres: SimpleProtocol and PrepareStmt are mutually exclusive — disable one of them")
	}

	pgxCfg, err := f.buildPgxConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("bootstrap/postgres: parse config: %w", err)
	}

	// Si SimpleProtocol esta habilitado, no usar prepared statements en pgx
	if options.SimpleProtocol {
		pgxCfg.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	}

	sqlDB := stdlib.OpenDB(*pgxCfg)
	f.applyPoolConfig(sqlDB, cfg)

	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("bootstrap/postgres: ping: %w", err)
	}

	gormCfg := &gorm.Config{
		PrepareStmt: options.PrepareStmt,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	// Configurar GORM logger si se proporciono
	if options.Logger != nil {
		if gl, ok := options.Logger.(gormlogger.Interface); ok {
			gormCfg.Logger = gl
		}
	}

	db, err := gorm.Open(gormpostgres.New(gormpostgres.Config{
		Conn: sqlDB,
	}), gormCfg)
	if err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("bootstrap/postgres: gorm open: %w", err)
	}

	return db, nil
}

// Ping verifica la conectividad con PostgreSQL.
func (f *Factory) Ping(ctx context.Context, db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("bootstrap/postgres: get sql.DB: %w", err)
	}
	return sqlDB.PingContext(ctx)
}

// Close cierra la conexion GORM.
func (f *Factory) Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("bootstrap/postgres: get sql.DB: %w", err)
	}
	return sqlDB.Close()
}

// buildPgxConfig construye la configuracion pgx a partir de PostgreSQLConfig.
func (f *Factory) buildPgxConfig(cfg bootstrap.PostgreSQLConfig) (*pgx.ConnConfig, error) {
	sslMode := cfg.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, sslMode,
	)

	pgxCfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// Configurar search_path si se especifico
	if cfg.SearchPath != "" {
		pgxCfg.RuntimeParams["search_path"] = cfg.SearchPath
	}

	return pgxCfg, nil
}

// applyPoolConfig configura el connection pool con valores del config o defaults.
func (f *Factory) applyPoolConfig(db *sql.DB, cfg bootstrap.PostgreSQLConfig) {
	maxOpen := cfg.MaxOpenConns
	if maxOpen == 0 {
		maxOpen = 25
	}
	maxIdle := cfg.MaxIdleConns
	if maxIdle == 0 {
		maxIdle = 5
	}
	maxLife := cfg.ConnMaxLifetime
	if maxLife == 0 {
		maxLife = time.Hour
	}
	maxIdleTime := cfg.ConnMaxIdleTime
	if maxIdleTime == 0 {
		maxIdleTime = 10 * time.Minute
	}

	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(maxLife)
	db.SetConnMaxIdleTime(maxIdleTime)
}
