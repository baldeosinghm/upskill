package users

import (
	"encoding/json"
	"errors"
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

type LoginResponse struct {
	UserID string `json:"user_id"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// r *http.Request has a package, Context(), that contains timelines,
	// deadlines, etc. (same as the context.Context package) to give
	// context, for lack of a better word, for cancellations
	user, err := h.service.Create(
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	id, err := h.service.Login(
		r.Context(),
		req.Email,
		req.Password,
	)

	if errors.Is(err, ErrInvalidCredentials) {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Printf("login error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// id is "" when err is non-nil; build it after all error checks
	var response LoginResponse
	response.UserID = id

	// Return something back to the user that can be used to verify identity
	// when accessing user features once logged in
	log.Printf("user id: %v", id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
