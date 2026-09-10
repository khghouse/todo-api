package todo

// MemoryRepository는 Todo를 프로세스 메모리의 slice에 저장한다.
type MemoryRepository struct {
	todos []Todo
}

// NewMemoryRepository는 초기 Todo 목록을 복사하여 메모리 저장소를 생성한다.
func NewMemoryRepository(initialTodos []Todo) *MemoryRepository {
	// 전달받은 slice와 backing array를 공유하지 않도록 새 slice를 만든다.
	todos := make([]Todo, len(initialTodos))
	copy(todos, initialTodos)

	return &MemoryRepository{todos: todos}
}

// FindAll은 저장된 모든 Todo의 복사본을 반환한다.
func (r *MemoryRepository) FindAll() ([]Todo, error) {
	// 호출자가 반환된 slice를 변경해도 저장소 내부 데이터가 바뀌지 않게 복사한다.
	todos := make([]Todo, len(r.todos))
	copy(todos, r.todos)

	// 메모리 조회에서는 오류가 발생하지 않으므로 error 자리에 nil을 반환한다.
	return todos, nil
}
