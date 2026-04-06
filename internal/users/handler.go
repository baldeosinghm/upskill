package users

import (
	"encoding/json"
	"log"
	"net/http"
)

type Handler struct {
	service *Service
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// r *http.Request has a package, Context(), that contains timelines,
	// deadlines, etc. (same as the context.Context package) to give
	// context, for lack of a better word, for cancellations
	user, err := h.service.CreateUser(
		r.Context(),
		req.Username,
		req.Email,
		req.Password,
		req.Role,
	)
	if err != nil {
		log.Printf("internal error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
