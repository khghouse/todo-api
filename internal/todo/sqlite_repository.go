package todo

import (
	"database/sql"
	"fmt"
	"time"
)

// SQLiteRepository는 SQLite 데이터베이스에 Todo를 저장하고 조회한다.
type SQLiteRepository struct {
	// sql.DB는 하나의 연결이 아니라 데이터베이스 연결 풀을 관리하는 객체다.
	db *sql.DB
}

// NewSQLiteRepository는 전달받은 DB를 사용하는 SQLite 저장소를 생성한다.
func NewSQLiteRepository(db *sql.DB) (*SQLiteRepository, error) {
	// 외부에서 생성한 DB를 저장소가 사용할 수 있도록 필드에 보관한다.
	repository := &SQLiteRepository{db: db}

	// 저장소를 사용하기 전에 todos 테이블이 존재하도록 초기화한다.
	if err := repository.createTable(); err != nil {
		// 생성에 실패하면 사용할 수 있는 저장소가 없으므로 nil과 오류를 반환한다.
		return nil, err
	}

	// 초기화에 성공한 저장소와 오류가 없음을 나타내는 nil을 반환한다.
	return repository, nil
}

// createTable은 todos 테이블이 없을 때만 새로 생성한다.
func (r *SQLiteRepository) createTable() error {
	// raw string literal(``)을 사용하면 여러 줄 SQL을 그대로 작성할 수 있다.
	const query = `
	CREATE TABLE IF NOT EXISTS todos (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		completed INTEGER NOT NULL,
		created_at TEXT NOT NULL
	)
	`

	// Exec은 조회 결과가 필요 없는 CREATE, INSERT, UPDATE, DELETE 문에 사용한다.
	// 첫 번째 반환값(sql.Result)은 여기서 필요하지 않으므로 _로 무시한다.
	if _, err := r.db.Exec(query); err != nil {
		// %w로 원래 오류를 감싸면 문맥을 추가하면서 errors.Is/As도 사용할 수 있다.
		return fmt.Errorf("create todos table: %w", err)
	}

	// 테이블 생성 또는 기존 테이블 확인에 성공했다.
	return nil
}

// FindAll은 todos 테이블의 모든 Todo를 생성 시각과 ID 순으로 조회한다.
func (r *SQLiteRepository) FindAll() ([]Todo, error) {
	const query = `
	SELECT id, title, completed, created_at
	FROM todos
	ORDER BY created_at, id
	`

	// Query는 여러 행을 조회하며, 결과를 순회할 수 있는 *sql.Rows를 반환한다.
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query todos: %w", err)
	}
	// 함수가 끝날 때 Rows가 점유한 데이터베이스 자원을 해제한다.
	defer rows.Close()

	// nil slice가 아닌 빈 slice로 시작하여 결과가 없을 때 JSON []로 응답하게 한다.
	todos := make([]Todo, 0)

	// Next는 다음 행으로 이동하며, 더 이상 행이 없으면 false를 반환한다.
	for rows.Next() {
		// 현재 행의 값을 담을 Todo와 DB 타입 변환용 임시 변수를 준비한다.
		var item Todo
		var completed int
		var createdAt string

		// Scan은 현재 행의 각 열 값을 전달한 변수의 주소에 순서대로 저장한다.
		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&completed,
			&createdAt,
		); err != nil {
			return nil, fmt.Errorf("scan todo: %w", err)
		}

		// SQLite에 TEXT로 저장된 RFC3339Nano 시각을 Go의 time.Time으로 변환한다.
		parsedCreatedAt, err := time.Parse(time.RFC3339Nano, createdAt)
		if err != nil {
			return nil, fmt.Errorf("parse created_at: %w", err)
		}

		// SQLite의 0/1 값을 Go의 bool로 변환하고, 변환한 시각도 Todo에 넣는다.
		item.Completed = completed != 0
		item.CreatedAt = parsedCreatedAt

		// 변환이 완료된 Todo를 결과 slice의 마지막에 추가한다.
		todos = append(todos, item)
	}

	// 순회 도중 발생한 DB 또는 드라이버 오류는 Next가 끝난 후 Err로 확인한다.
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate todos: %w", err)
	}

	// 모든 행을 정상적으로 읽었으므로 Todo 목록과 nil 오류를 반환한다.
	return todos, nil
}
