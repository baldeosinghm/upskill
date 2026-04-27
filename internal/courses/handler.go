package courses

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

type CourseRequest struct {
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CourseRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	course, err := h.service.Create(
		r.Context(),
		req.Name,
		req.OwnerID,
	)
	if errors.Is(err, ErrOwnerNotFound) {
		http.Error(w, "owner not found", http.StatusNotFound)
		return
	}
	if errors.Is(err, ErrOwnerNotTeacher) {
		http.Error(w, "owner is not a teacher", http.StatusForbidden)
		return
	}
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(course)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing course id", http.StatusBadRequest)
		return
	}
	course, err := h.service.GetByID(r.Context(), id)
	if errors.Is(err, ErrCourseNotFound) {
		http.Error(w, "course not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(course)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	courses, err := h.service.List(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(courses)
}
