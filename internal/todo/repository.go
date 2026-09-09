package todo

type Repository interface {
	FindAll() ([]Todo, error)
}
