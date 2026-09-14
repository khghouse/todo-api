package todo

// Repository는 Todo 저장소가 제공해야 하는 동작을 정의한다.
// Service는 구체적인 저장 방식(SQLite, 메모리 등)이 아니라 이 인터페이스에 의존한다.
type Repository interface {
	// FindAll은 저장된 모든 Todo를 반환한다.
	FindAll() ([]Todo, error)

	// Create는 전달받은 Todo를 저장한다.
	Create(item Todo) error
}
