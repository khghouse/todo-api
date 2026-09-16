package todo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_Create_정상적인요청이면201과생성된Todo를응답한다(t *testing.T) {
	// Arrange: 실제 DB 대신 메모리 저장소를 사용하여 Handler까지의 의존성을 조립한다.
	repository := NewMemoryRepository(nil)
	service := NewService(repository)
	handler := NewHandler(service)

	// httptest.NewRequest는 실제 네트워크 연결 없이 테스트용 HTTP 요청을 만든다.
	// 세 번째 인자는 요청 본문이며, strings.Reader가 JSON 문자열을 io.Reader로 제공한다.
	request := httptest.NewRequest(
		http.MethodPost,
		"/todos",
		strings.NewReader(`{"title": "Go 테스트 학습"}`),
	)
	// Content-Type은 요청 본문이 JSON 형식임을 서버에 알리는 HTTP 헤더다.
	request.Header.Set("Content-Type", "application/json")

	// ResponseRecorder는 Handler가 작성한 상태 코드, 헤더, 본문을 기록한다.
	recorder := httptest.NewRecorder()

	// Act: 라우터나 실제 서버를 거치지 않고 Handler 메서드를 직접 호출한다.
	handler.Create(recorder, request)

	// Assert: Handler가 응답 본문의 형식을 JSON으로 올바르게 명시했는지 검사한다.
	wantContentType := "application/json; charset=utf-8"
	gotContentType := recorder.Header().Get("Content-Type")

	if gotContentType != wantContentType {
		t.Errorf(
			"Content-Type이 다릅니다: want=%q, got=%q",
			wantContentType,
			gotContentType,
		)
	}

	// 기록된 응답 JSON을 Todo로 역직렬화한다.
	// 디코딩에 실패하면 이후 필드 검증이 의미 없으므로 Fatalf로 테스트를 중단한다.
	var response Todo
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("응답 JSON을 해석할 수 없습니다: %v", err)
	}

	// Service가 ID, 정리된 제목, 기본 완료 상태, 생성 시각을 설정했는지 검사한다.
	if response.ID == "" {
		t.Error("응답 Todo에 ID가 있어야 합니다")
	}

	if response.Title != "Go 테스트 학습" {
		t.Errorf(
			"응답 Todo의 제목이 다릅니다: want=%q, got=%q",
			"Go 테스트 학습",
			response.Title,
		)
	}

	if response.Completed {
		t.Error("응답 Todo의 completed는 false여야 합니다")
	}

	if response.CreatedAt.IsZero() {
		t.Error("응답 Todo의 CreatedAt이 설정되어야 합니다")
	}

	// HTTP 응답뿐 아니라 생성된 Todo가 저장소에도 실제로 들어갔는지 검사한다.
	stored, err := repository.FindAll()
	if err != nil {
		t.Fatalf("저장된 Todo 조회에 실패했습니다: %v", err)
	}

	// 개수가 다르면 stored[0]에 안전하게 접근할 수 없으므로 Fatalf로 중단한다.
	if len(stored) != 1 {
		t.Fatalf("저장된 Todo 수가 일치하지 않습니다: want=1, got=%d", len(stored))
	}

	// 저장된 객체와 응답 객체가 같은 Todo인지 ID로 확인한다.
	if stored[0].ID != response.ID {
		t.Errorf(
			"저장된 Todo와 응답 Todo의 ID가 다릅니다: stored=%q, response=%q",
			stored[0].ID,
			response.ID,
		)
	}
}

func TestHandler_Create_빈제목이면400을응답한다(t *testing.T) {
	// Arrange: 각 테스트마다 새 저장소를 사용해 다른 테스트와 상태를 공유하지 않는다.
	repository := NewMemoryRepository(nil)
	service := NewService(repository)
	handler := NewHandler(service)

	// JSON 문법은 유효하지만 title이 공백뿐인 요청을 만든다.
	request := httptest.NewRequest(
		http.MethodPost, "/todos",
		strings.NewReader(`{"title": "   "}`),
	)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	// Act: Service의 제목 검증 오류를 Handler가 어떤 HTTP 응답으로 변환하는지 확인한다.
	handler.Create(recorder, request)

	// Assert: ErrTitleRequired는 클라이언트 입력 오류이므로 400이어야 한다.
	if recorder.Code != http.StatusBadRequest {
		t.Errorf(
			"HTTP 상태 코드가 다릅니다: want=%d, got=%d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}
}

func TestHandler_Create_잘못된JSON이면400을응답한다(t *testing.T) {
	// Arrange: Handler의 JSON 파싱 동작을 확인하기 위한 의존성을 준비한다.
	repository := NewMemoryRepository(nil)
	service := NewService(repository)
	handler := NewHandler(service)

	// 닫히지 않은 JSON을 본문으로 전달하여 의도적으로 Decode 오류를 발생시킨다.
	request := httptest.NewRequest(
		http.MethodPost,
		"/todos",
		strings.NewReader(`{"title":`),
	)
	recorder := httptest.NewRecorder()

	// Act: JSON 파싱에 실패하므로 Service.Create까지 호출되지 않아야 한다.
	handler.Create(recorder, request)

	// Assert: 해석할 수 없는 요청 본문은 400 Bad Request로 응답해야 한다.
	if recorder.Code != http.StatusBadRequest {
		t.Errorf(
			"HTTP 상태 코드가 다릅니다: want=%d, got=%d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}
}
