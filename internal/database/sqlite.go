package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// OpenSQLiteInMemory는 SQLite 인메모리 데이터베이스를 열고 연결을 확인한다.
func OpenSQLiteInMemory() (*sql.DB, error) {
	db, err := sql.Open( // Go 표준 라이브러리의 범용 SQL 인터페이스
		"sqlite", // 드라이버 이름
		"file:todo-api?mode=memory&cache=shared",
		// file:todo-api: 인메모리 DB의 논리적인 이름
		// mode=memory: 디스크 파일이 아닌 메모리에 생성
		// cache=shared: 여러 SQLite 연결이 같은 인메모리 DB를 공유
	)

	// SQLite 연결을 열 때 오류가 발생하면 오류를 반환한다.
	if err != nil {
		return nil, fmt.Errorf("open SQLite: %w", err)
	}

	// 인메모리 DB가 여러 물리 연결로 분리되지 않도록 연결을 하나로 제한한다.
	db.SetMaxOpenConns(1)

	// SQLite 연결을 확인한다.
	if err := db.Ping(); err != nil {
		_ = db.Close() // DB 핸들을 닫음
		return nil, fmt.Errorf("ping SQLite: %w", err)
	}

	return db, nil
}
