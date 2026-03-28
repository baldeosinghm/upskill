package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/baldeosinghm/upskill/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Create a context for startup.  You always need to create one.
	// This is the exact way of creating a context.
	ctx := context.Background()

	// 2. Connect to the database
	err := godotenv.Load()
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

	// 3. Set up router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// 4. Start up server
	log.Println("upskill-api starting on :8080...")
	// Pass chi as the router for http requests
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
