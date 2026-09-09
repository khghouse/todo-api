package todo

type MemoryRepository struct {
	todos []Todo
}

func NewMemoryRepository(initialTodos []Todo) *MemoryRepository {
	todos := make([]Todo, len(initialTodos))
	copy(todos, initialTodos)

	return &MemoryRepository{todos: todos}
}

func (r *MemoryRepository) FindAll() ([]Todo, error) {
	todos := make([]Todo, len(r.todos))
	copy(todos, r.todos)

	return todos, nil
}
