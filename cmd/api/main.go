package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/baldeosinghm/upskill/internal/db"
	"github.com/baldeosinghm/upskill/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Create a context for startup.  You always need to create one.
	// This is the exact way of creating a context.
	ctx := context.Background()

	// 2. Connect to the database
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to read database url from env file: %v", err)
	}
	connStr := os.Getenv("DATABASE_URL")
	pool, err := db.NewPool(ctx, connStr)

	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	} else {
		defer pool.Close()
	}
	log.Println("database connection established")

	// Run migrations on startup
	if err := db.RunMigrations(connStr); err != nil {
		log.Println(err)
	}

	// 3. Set up router + route requests
	r := routes.RegisterRoutes(pool)

	// 4. Start up server (this happen last)
	log.Println("upskill-api starting on :8080...")
	// Pass chi as the router for http requests
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
