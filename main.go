package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"todo-api/internal/todo"
)

var todos = []todo.Todo{
	{ID: "1", Title: "Buy groceries", Completed: false, CreatedAt: time.Now()},
	{ID: "2", Title: "Walk the dog", Completed: true, CreatedAt: time.Now()},
	{ID: "3", Title: "Read a book", Completed: false, CreatedAt: time.Now()},
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Fprintln(w, `{"status":"ok"}`)
}

func todosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(todos); err != nil {
		http.Error(w, "failed to encode todos", http.StatusInternalServerError)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /todos", todosHandler)

	log.Println("server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
