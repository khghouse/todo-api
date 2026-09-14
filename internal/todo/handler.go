package todo

import (
	"encoding/json"
	"errors"
	"net/http"
)

// Handler는 Todo HTTP 요청을 처리한다.
type Handler struct {
	// HTTP 처리에 필요한 비즈니스 로직을 Service에 위임한다.
	service *Service
}

// createTodoRequest는 Todo 등록 요청에서 클라이언트가 전달할 JSON 형식을 나타낸다.
// ID, 완료 여부, 생성 시각은 서버가 정하므로 요청에서는 title만 받는다.
type createTodoRequest struct {
	Title string `json:"title"`
}

// NewHandler는 Service를 주입받아 Handler를 생성한다.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// FindAll은 모든 Todo를 조회하여 JSON으로 응답한다.
func (h *Handler) FindAll(w http.ResponseWriter, r *http.Request) {
	// HTTP Handler는 저장소를 직접 호출하지 않고 Service를 통해 Todo를 조회한다.
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

// Create는 JSON 요청에서 제목을 읽어 새로운 Todo를 등록한다.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	// Decode한 요청 데이터를 담을 zero value 구조체를 준비한다.
	var request createTodoRequest

	// 요청 본문을 순차적으로 읽으면서 JSON을 Go 구조체로 변환한다.
	decoder := json.NewDecoder(r.Body)
	// DTO에 정의되지 않은 필드가 들어오면 오타나 잘못된 요청으로 판단한다.
	decoder.DisallowUnknownFields()

	// JSON 문법이 잘못됐거나 구조체로 변환할 수 없으면 400으로 응답한다.
	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 제목 검증과 Todo 생성 및 저장은 Service에 위임한다.
	todo, err := h.service.Create(request.Title)
	if err != nil {
		// errors.Is는 오류가 나중에 %w로 감싸져도 특정 원인인지 확인할 수 있다.
		if errors.Is(err, ErrTitleRequired) {
			http.Error(
				w, err.Error(), http.StatusBadRequest,
			)
			return
		}

		// 검증 오류가 아닌 저장소 오류 등은 서버 내부 오류로 처리한다.
		http.Error(
			w, "failed to create todo", http.StatusInternalServerError,
		)
		return
	}

	// Header는 상태 코드나 응답 본문을 쓰기 전에 설정해야 한다.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	// 생성된 Todo를 JSON으로 변환하여 201 응답 본문에 작성한다.
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		return
	}
}
