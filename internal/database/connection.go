package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"

	"portfolio-backend/internal/config"
)

type DB struct {
	*sql.DB
}

func NewConnection(cfg *config.DatabaseConfig) (*DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool as specified in CLAUDE.md
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info().
		Str("host", cfg.Host).
		Int("port", cfg.Port).
		Str("database", cfg.Database).
		Int("max_open_conns", cfg.MaxOpenConns).
		Int("max_idle_conns", cfg.MaxIdleConns).
		Dur("conn_max_lifetime", cfg.ConnMaxLifetime).
		Msg("Database connection established successfully")

	return &DB{DB: db}, nil
}

func (db *DB) Close() error {
	if db.DB == nil {
		return nil
	}

	if err := db.DB.Close(); err != nil {
		log.Error().Err(err).Msg("Failed to close database connection")
		return err
	}

	log.Info().Msg("Database connection closed successfully")
	return nil
}

func (db *DB) Health(ctx context.Context) error {
	if db.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.DB.PingContext(pingCtx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}