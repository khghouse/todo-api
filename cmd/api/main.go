package main

import (
	"log"
	"net/http"
	"time"
	"todo-api/internal/health"
	"todo-api/internal/todo"
)

func main() {
	now := time.Now().UTC()

	repository := todo.NewMemoryRepository([]todo.Todo{
		{ID: "1", Title: "Buy groceries", Completed: false, CreatedAt: now},
		{ID: "2", Title: "Clean the house", Completed: true, CreatedAt: now},
		{ID: "3", Title: "Finish the project", Completed: false, CreatedAt: now},
	})

	todoHandler := todo.NewHandler(repository)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", health.Handler)
	mux.HandleFunc("/todos", todoHandler.FindAll)

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
