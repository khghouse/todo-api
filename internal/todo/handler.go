package todo

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	repository Repository
}

func NewHandler(repository Repository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) FindAll(w http.ResponseWriter, r *http.Request) {
	todos, err := h.repository.FindAll()
	if err != nil {
		http.Error(w, "Failed to retrieve todos", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(todos); err != nil {
		http.Error(w, "failed to encode todos", http.StatusInternalServerError)
	}
}
