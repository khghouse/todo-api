# Todo API

Go로 CRUD, REST API, 테스트, 데이터베이스 연동을 학습하기 위한 Todo 관리 API 프로젝트입니다.

## 개발 환경

- macOS
- Go
- VS Code
- Git

## Go 설치

Homebrew가 설치되어 있다면 터미널에서 아래 명령을 실행합니다.

```bash
brew install go
```

설치가 완료되면 버전을 확인합니다.

```bash
go version
```

다음과 비슷한 결과가 나오면 정상입니다.

```text
go version go1.27.1 darwin/arm64
```

Homebrew가 없다면 [Go 공식 다운로드 페이지](https://go.dev/dl/)에서 macOS 설치 파일을 내려받아 설치할 수 있습니다.

## VS Code 설정

1. VS Code를 실행합니다.
2. Extensions 화면을 엽니다. (`Cmd + Shift + X`)
3. `Go`를 검색합니다.
4. 게시자가 **Go Team at Google**인 공식 `Go` 확장을 설치합니다.

터미널에서 설치하려면 아래 명령을 사용할 수도 있습니다.

```bash
code --install-extension golang.go
```

### Go 개발 도구 설치

Go 확장 설치 후 VS Code에서 `Cmd + Shift + P`를 누르고 아래 명령을 실행합니다.

```text
Go: Install/Update Tools
```

다음 도구를 선택해 설치합니다.

- `gopls`: 코드 자동완성, 오류 분석, 정의 이동
- `goimports`: Go 코드 포맷 및 import 자동 정리
- `delve`: 디버깅 도구

## 프로젝트 실행

프로젝트 최상위 폴더에서 아래 명령을 실행합니다.

```bash
go run .
```
