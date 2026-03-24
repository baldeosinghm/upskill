package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
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
		return nil, fmt.Errorf("Unable to create connection pool: %w", err)
	}
	// Reach out to DB and verify connection works
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("Unable to ping database: %w", err)
	}
	return pool, nil
}
