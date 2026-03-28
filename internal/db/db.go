package db

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // registers postgres driver for migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"       // registers file source
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

// Set up database connection
//
// Passing context.Context as the first argument is a Go convention for anything
// that touches I/O.  It should only be passed as a parameter to a function and
// never to a struct type.
func NewPool(ctx context.Context, connStr string) (*pgxpool.Pool, error) {
	// Initialize and confiure a pool to allow for conncurrency.
	// This means we make multiple requests to DB possible
	// pgx requires context to internally set a timeout on initial connection attempt
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}
	// Reach out to DB and verify connection works
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}
	return pool, nil
}

// Run migrations on start up
//
// Migrations manages tables; also enables us to create them
func RunMigrations(connStr string) error {
	m, err := migrate.New("file://migrations", connStr)
	if err != nil {
		return fmt.Errorf("migration init failed: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}
	return nil
}
