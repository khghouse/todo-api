package todo

import "sync"

// MemoryRepository는 Todo를 프로세스 메모리의 slice에 저장한다.
type MemoryRepository struct {
	// RWMutex는 여러 조회는 동시에 허용하고, 데이터 변경 중에는 다른 접근을 막는다.
	mu    sync.RWMutex
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
	// 데이터를 읽는 동안 쓰기가 발생하지 않도록 읽기 잠금을 획득한다.
	r.mu.RLock()
	defer r.mu.RUnlock() // 함수가 끝날 때 읽기 잠금을 반드시 해제한다.

	// 호출자가 반환된 slice를 변경해도 저장소 내부 데이터가 바뀌지 않게 복사한다.
	todos := make([]Todo, len(r.todos))
	copy(todos, r.todos)

	// 메모리 조회에서는 오류가 발생하지 않으므로 error 자리에 nil을 반환한다.
	return todos, nil
}

// Create는 전달받은 Todo를 메모리의 slice에 추가한다.
func (r *MemoryRepository) Create(item Todo) error {
	// append는 slice를 변경하므로 쓰기 잠금으로 다른 읽기와 쓰기를 모두 막는다.
	r.mu.Lock()
	defer r.mu.Unlock()

	// 새로운 Todo를 저장소의 slice에 추가한다.
	r.todos = append(r.todos, item)

	// 메모리 저장에서는 오류가 발생하지 않으므로 error 자리에 nil을 반환한다.
	return nil
}
