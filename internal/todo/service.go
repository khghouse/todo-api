package todo

import "fmt"

// Repository 인터페이스를 구현한 저장소를 주입받는다.
type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) FindAll() ([]Todo, error) {
	todos, err := s.repository.FindAll()
	if err != nil {
		return nil, fmt.Errorf("find all todos: %w", err)
	}

	return todos, nil
}
