package main

import (
	"log"
	"net/http"
	"time"
	"todo-api/internal/database"
	"todo-api/internal/health"
	"todo-api/internal/todo"
)

func main() {
	db, err := database.OpenSQLiteInMemory() // SQLite 인메모리 데이터베이스를 연다.
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil { // 프로그램 종료 시 데이터베이스 연결을 닫는다.
			log.Printf("failed to close database: %v", err)
		}
	}()

	log.Println("SQLite in-memory database connected")

	now := time.Now().UTC() // 현재 시각을 UTC 기준으로 표현한다.

	// 초기 Todo slice를 복사하여 메모리 저장소를 생성한다.
	repository := todo.NewMemoryRepository([]todo.Todo{
		{ID: "1", Title: "Buy groceries", Completed: false, CreatedAt: now},
		{ID: "2", Title: "Clean the house", Completed: true, CreatedAt: now},
		{ID: "3", Title: "Finish the project", Completed: false, CreatedAt: now},
	})

	// 메모리 저장소를 Repository 인터페이스로 Todo 핸들러에 주입한다.
	todoHandler := todo.NewHandler(repository)

	// 표준 라이브러리 라우터를 생성하고 요청 경로에 핸들러를 등록한다.
	mux := http.NewServeMux()
	mux.HandleFunc("/health", health.Handler)     // 서버 상태 확인 핸들러를 등록한다.
	mux.HandleFunc("/todos", todoHandler.FindAll) // Todo 목록 조회 핸들러를 등록한다.

	log.Println("Server is running on http://localhost:8080")

	// 8080 포트에서 서버를 실행하고 요청 처리를 mux에 위임한다.
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Printf("server stopped: %v", err)
	}
}
