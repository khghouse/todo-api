package todo

import "time"

// 현재는 데이터 모델링 객체이자 응답 객체
type Todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"createdAt"`
}
