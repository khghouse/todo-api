package todo

import (
	"encoding/json"
	"net/http"
)

// Handler는 Todo HTTP 요청을 처리한다.
type Handler struct {
	service *Service
}

// NewHandler는 Service를 주입받아 Handler를 생성한다.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// FindAll은 모든 Todo를 조회하여 JSON으로 응답한다.
func (h *Handler) FindAll(w http.ResponseWriter, r *http.Request) {
	// Service Todo 목록과 조회 중 발생한 오류를 함께 반환한다.
	todos, err := h.service.FindAll()
	if err != nil {
		http.Error(w, "Failed to retrieve todos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Todo slice를 JSON으로 변환해 HTTP 응답 본문에 작성한다.
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		http.Error(w, "failed to encode todos", http.StatusInternalServerError)
	}
}
