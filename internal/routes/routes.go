package routes

import (
	"encoding/json"
	"net/http"

	"github.com/baldeosinghm/upskill/internal/courses"
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

	// Users service
	uRepo := users.NewRepository(db)
	uService := users.NewService(uRepo)
	uHandler := users.NewHandler(uService)

	// Courses service
	cRepo := courses.NewRepository(db)
	cService := courses.NewService(cRepo, uService)
	cHandler := courses.NewHandler(cService)

	// Route requests
	r.Post("/users", uHandler.Create)
	r.Post("/login", uHandler.Login)

	r.Route("/courses", func(r chi.Router) {
		r.Post("/", cHandler.Create)
		r.Get("/{id}", cHandler.GetByID)
		r.Get("/", cHandler.List)
	})

	return r
}
