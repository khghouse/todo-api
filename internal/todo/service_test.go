package todo

import (
	"errors"
	"testing"
)

func TestService_Create(t *testing.T) {
	// t.Run은 하나의 테스트 함수 안에서 독립적으로 표시되는 하위 테스트를 실행한다.
	t.Run("정상적인 제목이면 Todo를 생성하고 저장된다.", func(t *testing.T) {
		// Arrange: 실제 DB 없이 비즈니스 로직을 검사하도록 메모리 저장소를 주입한다.
		repository := NewMemoryRepository(nil)
		service := NewService(repository)

		// Act: 앞뒤 공백이 포함된 제목으로 Todo 생성을 요청한다.
		created, err := service.Create("   Go 테스트 학습   ")

		// 예상하지 않은 오류가 있으면 이후 결과 검증이 의미 없으므로 즉시 중단한다.
		if err != nil {
			t.Fatalf("예상하지 못한 오류가 발생했습니다: %v", err)
		}

		// Assert: ID, 정리된 제목, 기본 상태, 생성 시각이 올바른지 검사한다.
		if created.ID == "" {
			t.Error("ID가 생성되어야 합니다.")
		}

		if created.Title != "Go 테스트 학습" {
			t.Errorf(
				"제목이 일치하지 않습니다: want=%q, got=%q",
				"Go 테스트 학습",
				created.Title,
			)
		}

		if created.Completed {
			t.Error("새 Todo의 completed는 false여야 합니다")
		}

		if created.CreatedAt.IsZero() {
			t.Error("생성 시각이 설정되어야 합니다")
		}

		// 반환값뿐 아니라 Repository.Create를 통해 실제로 저장되었는지도 확인한다.
		stored, err := repository.FindAll()
		if err != nil {
			t.Fatalf("저장된 Todo 조회에 실패했습니다: %v", err)
		}

		// 개수가 다르면 stored[0] 접근 시 panic이 날 수 있으므로 Fatalf로 중단한다.
		if len(stored) != 1 {
			t.Fatalf("저장된 Todo 수가 일치하지 않습니다: want=1, got=%d", len(stored))
		}

		if stored[0].ID != created.ID {
			t.Errorf("저장된 Todo의 ID가 일치하지 않습니다: want=%q, got=%q", created.ID, stored[0].ID)
		}
	})

	// 입력값만 다르고 기대 결과가 같은 경우를 한 번에 검증하는 테이블 테스트다.
	tests := []struct {
		name  string // 하위 테스트의 이름
		title string // Service.Create에 전달할 입력값
	}{
		{name: "빈 문자열", title: ""},
		{name: "공백만 있는 문자열", title: "   "},
		{name: "탭과 줄바꿈만 있는 문자열", title: "\t\n"},
	}

	// 각 테스트 데이터를 순회하며 동일한 검증 로직을 실행한다.
	for _, tt := range tests {
		t.Run(tt.name+"이면 등록에 실패한다", func(t *testing.T) {
			// 각 하위 테스트가 독립된 상태를 갖도록 새 저장소와 Service를 생성한다.
			repository := NewMemoryRepository(nil)
			service := NewService(repository)

			// 반환되는 Todo는 이 테스트의 관심 대상이 아니므로 _로 무시한다.
			_, err := service.Create(tt.title)

			// errors.Is는 오류가 %w로 감싸진 경우에도 같은 원인인지 확인할 수 있다.
			if !errors.Is(err, ErrTitleRequired) {
				t.Errorf("ErrTitleRequired를 기대했지만 다른 결과를 받았습니다: %v", err)
			}

			stored, err := repository.FindAll()
			if err != nil {
				t.Fatalf("저장된 Todo 조회에 실패했습니다: %v", err)
			}

			// 입력 검증에 실패했다면 Repository에는 아무 Todo도 저장되지 않아야 한다.
			if len(stored) != 0 {
				t.Errorf("검증 실패 시 저장되면 안 됩니다: got=%d", len(stored))
			}
		})
	}
}
