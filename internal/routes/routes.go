package routes

import (
	"encoding/json"
	"net/http"

	"github.com/baldeosinghm/upskill/internal/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(db *pgxpool.Pool) *chi.Mux {
	// Set up router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Route requests
	repo := users.NewRepository(db)
	service := users.NewService(repo)
	handler := users.NewHandler(service)
	r.Post("/users", handler.CreateUser)
	// r.Get("/users/{email}", handler.FindByEmail)
	// r.Get("/users?email={email}&password={password}", handler.Login)

	return r
}
