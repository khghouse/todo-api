package todo

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"uuid"
)

// ErrTitleRequired는 공백을 제외한 Todo 제목이 없을 때 반환하는 검증 오류다.
var ErrTitleRequired = errors.New("title is required")

// Service는 Todo의 비즈니스 규칙을 처리하고 데이터 저장은 Repository에 위임한다.
type Service struct {
	// Repository 인터페이스에 의존하므로 SQLite와 메모리 구현체를 교체할 수 있다.
	repository Repository
}

// NewService는 Repository를 주입받아 Todo 서비스를 생성한다.
func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

// FindAll은 저장소에 보관된 모든 Todo를 조회한다.
func (s *Service) FindAll() ([]Todo, error) {
	todos, err := s.repository.FindAll()
	if err != nil {
		// 원본 오류를 %w로 감싸 호출자에게 조회 작업의 문맥과 원인을 함께 전달한다.
		return nil, fmt.Errorf("find all todos: %w", err)
	}

	return todos, nil
}

// Create는 제목을 검증하고 기본값을 채운 Todo를 생성하여 저장한다.
func (s *Service) Create(title string) (Todo, error) {
	// 앞뒤 공백을 제거하여 공백만 있는 제목도 빈 제목으로 취급한다.
	title = strings.TrimSpace(title)
	if title == "" {
		// Todo 생성에 실패했으므로 Todo의 zero value와 검증 오류를 반환한다.
		return Todo{}, ErrTitleRequired
	}

	// ID와 기본 상태 및 생성 시각은 클라이언트가 아니라 서버에서 결정한다.
	item := Todo{
		ID:        uuid.New().String(),
		Title:     title,
		Completed: false,
		CreatedAt: time.Now().UTC(),
	}

	// Service가 완성한 Todo의 실제 저장은 주입받은 Repository에 위임한다.
	if err := s.repository.Create(item); err != nil {
		// 저장 오류에 작업 문맥을 추가하되 원본 오류는 %w로 보존한다.
		return Todo{}, fmt.Errorf("create todo: %w", err)
	}

	// 저장에 성공한 Todo를 Handler가 응답으로 사용할 수 있도록 반환한다.
	return item, nil
}
