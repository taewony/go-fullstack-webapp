# Go Fullstack Web Programming with GOTH stack and PostgreSQL (based on "GO Web Programming Book" chichat service)

## "GO Web Programming" 책의 chitchat 서비스 개발 Step-by-Step Guide

**개발 환경:**

- **운영체제:** Windows  
- **개발 툴:** VS Code  
- **프로그래밍 언어:** Go 1.24

**핵심 목표:**

- **Go 언어 기본 다지기:** Go 1.24 표준 라이브러리를 최대한 활용하여 웹 개발 기본기를 탄탄하게 다집니다.  
- **최소 외부 라이브러리:** net/http 표준 라이브러리의 강력함을 경험하고 웹 개발 핵심 원리 이해에 집중합니다.  
- **Full-Stack 개발:** HTMX와 templ을 사용하여 생산성 높고 유지보수 용이한 Full-Stack 웹 애플리케이션 개발 방식을 익힙니다.  
- **데이터 중심 설계:** chitchat 서비스의 데이터를 기반으로 가장 간단한 웹 페이지부터 시작하여 점진적으로 기능을 확장해 나갑니다.  
- **DB 연동:** 처음에는 SQLite3 in-memory와 sqlx를 사용하여 간편하게 데이터베이스를 연동하고 CRUD를 구현하며, 추후 pgx와 PostgreSQL로 전환하여 실제 운영 환경을 경험합니다.

---

## 단계별 개발 가이드:

### 0단계: 개발 환경 설정 (Windows & VS Code)

1. **Go 설치:** Go 1.24 버전을 [공식 Go 웹사이트](https://go.dev/dl/)에서 다운로드하여 Windows에 설치합니다. 설치 후, `go version` 명령어를 터미널 (VS Code 터미널 또는 Windows PowerShell) 에 입력하여 Go가 제대로 설치되었는지 확인합니다.  
2. **VS Code 설치 및 Go 확장 설치:** VS Code가 설치되어 있지 않다면 [VS Code 웹사이트](https://code.visualstudio.com/)에서 다운로드하여 설치합니다. VS Code를 실행하고 확장 탭에서 "Go" 확장 (by Go Team at Google) 을 검색하여 설치합니다.  
3. **Go 프로젝트 폴더 생성:** VS Code에서 "File" -> "Open Folder..." 를 선택하거나, 적절한 위치에 chitchat 프로젝트 폴더를 생성합니다. 예시: `C:\workspace\chitchat`  
4. **go.mod 초기화:** VS Code 터미널을 열고 (Ctrl + ` 또는 터미널 메뉴), 프로젝트 폴더로 이동한 후 다음 명령어를 실행하여 Go modules를 초기화합니다.  
   ```bash
   go mod init github.com/your-username/chitchat  # github.com/your-username 부분은 실제 github 사용자 이름 또는 프로젝트 경로로 변경
   ```
   이 명령어는 go.mod 파일을 생성하며, 프로젝트 의존성 관리를 시작합니다.

### 1단계: 가장 간단한 웹 페이지 만들기 (Hello, World!)

1. **main.go 파일 생성:** 프로젝트 폴더에 main.go 파일을 생성합니다.  
2. **기본 웹 서버 코드 작성:** main.go 파일에 다음 코드를 작성합니다.  
   ```go
   package main

   import (
       "fmt"
       "net/http"
   )

   func handler(w http.ResponseWriter, r *http.Request) {
       fmt.Fprintln(w, "Hello, World!")
   }

   func main() {
       http.HandleFunc("/", handler) // "/" 경로로 요청이 들어오면 handler 함수 실행
       fmt.Println("chitchat 서버 시작... http://localhost:8080")
       http.ListenAndServe(":8080", nil) // 8080 포트에서 웹 서버 시작
   }
   ```

3. **실행 및 확인:** VS Code 터미널에서 `go run main.go` 명령어를 실행합니다. 터미널에 "chitchat 서버 시작... http://localhost:8080" 메시지가 출력되면 웹 서버가 정상적으로 시작된 것입니다. 웹 브라우저를 열고 http://localhost:8080 주소로 접속하면 "Hello, World!" 메시지가 표시되는 것을 확인할 수 있습니다.

### 2단계: templ 사용하여 HTML 템플릿 적용

1. **templ 설치:** VS Code 터미널에서 다음 명령어를 실행하여 templ CLI를 설치합니다.  
   ```bash
   go install github.com/a-h/templ/cmd/templ@latest
   ```
   (만약 templ 명령어를 찾을 수 없다는 오류가 발생하면, Go bin 폴더 (예: `C:\Users\YourUsername\go\bin`) 가 시스템 환경 변수 Path 에 등록되어 있는지 확인하고, VS Code를 재시작해 보세요.)  
2. **templ 파일 생성:** 프로젝트 폴더에 index.templ 파일을 생성하고, 간단한 HTML 구조를 작성합니다.  
   ```html
   <!DOCTYPE html>
   <html>
   <head>
       <title>Chitchat</title>
   </head>
   <body>
       <h1>Welcome to Chitchat!</h1>
       <p>This is a simple chitchat service.</p>
   </body>
   </html>
   ```

3. **Go 코드 수정 (templ 렌더링):** main.go 파일을 다음과 같이 수정합니다.  
   ```go
   package main

   import (
       "fmt"
       "net/http"
       "chitchat/templates" // templates 패키지 import (go.mod 파일 경로에 따라 수정)
   )

   func handler(w http.ResponseWriter, r *http.Request) {
       templates.Index().Render(r.Context(), w) // index.templ 렌더링
   }

   func main() {
       http.HandleFunc("/", handler)
       fmt.Println("chitchat 서버 시작... http://localhost:8080")
       http.ListenAndServe(":8080", nil)
   }
   ```

4. **templates 패키지 생성 및 templ 파일 컴파일:**  
   - 프로젝트 폴더에 templates 라는 폴더를 생성합니다.  
   - index.templ 파일을 templates 폴더로 이동합니다.  
   - VS Code 터미널에서 프로젝트 루트 폴더로 이동 후, 다음 명령어를 실행하여 index.templ 파일을 Go 코드로 컴파일합니다.  
     ```bash
     templ generate
     ```
     이 명령어는 templates 폴더 안에 index_templ.go 파일을 생성합니다. go.mod 파일에 github.com/a-h/templ 의존성이 자동으로 추가됩니다.  
5. **실행 및 확인:** `go run main.go` 명령어를 실행하고, 웹 브라우저에서 http://localhost:8080 에 접속하면, index.templ 에 작성한 HTML 구조가 표시되는 것을 확인할 수 있습니다.

### 3단계: SQLite3 in-memory 데이터베이스 연동 및 데이터 표시

1. **sqlx 라이브러리 추가:** VS Code 터미널에서 다음 명령어를 실행하여 sqlx 라이브러리를 프로젝트에 추가합니다.  
   ```bash
   go get github.com/jmoiron/sqlx  
   go get modernc.org/sqlite // 기존 github.com/mattn/go-sqlite3 대체
   ```

2. **데이터 모델 정의 (Post):** main.go 파일 또는 별도의 파일 (예: models.go) 에 간단한 게시글 데이터 모델 (Post) 을 정의합니다.  
   ```go
   package main // 또는 package models

   type Post struct {
       ID      int    `db:"id"`
       Content string `db:"content"`
       Author  string `db:"author"`
   }
   ```

3. **데이터베이스 초기화 및 데이터 삽입:** main.go 파일의 main 함수 안에 데이터베이스 초기화 및 샘플 데이터 삽입 코드를 추가합니다.  
   ```go
   package main

   import (
       "database/sql"
       "fmt"
       "log"
       "net/http"
       "chitchat/templates"
       "github.com/jmoiron/sqlx"
       _ "github.com/mattn/go-sqlite3" // sqlite3 드라이버 import (init 함수 실행)
   )

   var db *sqlx.DB

   func handler(w http.ResponseWriter, r *http.Request) {
       posts, err := getPosts() // 게시글 목록 조회 함수 호출 (아래에 구현)
       if err != nil {
           http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
           return
       }
       templates.Index(posts).Render(r.Context(), w) // index.templ 에 posts 데이터 전달하며 렌더링
   }

   func getPosts() ([]Post, error) {
       var posts []Post
       err := db.Select(&posts, "SELECT id, content, author FROM posts") // sqlx.DB.Select 사용하여 복수 row 조회
       if err != nil {
           return nil, err
       }
       return posts, nil
   }

   func main() {
       var err error
       db, err = sqlx.Open("sqlite3", ":memory:") // in-memory SQLite3 데이터베이스 연결
       if err != nil {
           log.Fatal(err)
       }
       defer db.Close()

       _, err = db.Exec(`
           CREATE TABLE IF NOT EXISTS posts (
               id INTEGER PRIMARY KEY AUTOINCREMENT,
               content TEXT NOT NULL,
               author TEXT NOT NULL
           );
       `) // posts 테이블 생성
       if err != nil {
           log.Fatal(err)
       }

       _, err = db.Exec(`
           INSERT INTO posts (content, author) VALUES (?, ?), (?, ?)
       `, "첫 번째 게시글 내용입니다.", "John Doe", "두 번째 게시글입니다.", "Jane Doe") // 샘플 데이터 삽입
       if err != nil {
           log.Fatal(err)
       }

       http.HandleFunc("/", handler)
       fmt.Println("chitchat 서버 시작... http://localhost:8080")
       http.ListenAndServe(":8080", nil)
   }
   ```

4. **templ 파일 수정 (데이터 표시):** templates/index.templ 파일을 수정하여 게시글 목록을 표시하도록 합니다.  
   ```html
   <!DOCTYPE html>
   <html>
   <head>
       <title>Chitchat</title>
   </head>
   <body>
       <h1>Chitchat 게시판</h1>
       <ul>
           @for post := range posts {
               <li>
                   <strong>{ post.Author }:</strong> { post.Content }
               </li>
           }
       </ul>
   </body>
   </html>

   @code script(posts []Post) {
       // index.templ 함수는 posts 라는 Post 슬라이스를 인자로 받도록 정의
   }
   ```

5. **templ 파일 컴파일 및 실행:** `templ generate` 명령어를 다시 실행하여 변경된 index.templ 파일을 컴파일합니다. `go run main.go` 명령어를 실행하고, 웹 브라우저에서 http://localhost:8080 에 접속하면, 데이터베이스에 저장된 게시글 목록이 웹 페이지에 표시되는 것을 확인할 수 있습니다.

### 4단계: Router/Handler 추가 및 templ 파일 정비

1. **data 중심 web Server 설계**
   ![Web Server 구조](public/images/web_app.png)

2. **Router 및 Handler 추가:** main.go 내용을 간단히 하고, router.go 파일에 url에 handler 함수를 연결합니다.
   **main.go**  
   ```go
   package main

   import (
       "net/http"
       "github.com/taewony/go-fullstack-webapp/internal/router"
   )

   func main() {
       router := router.NewRouter()
       http.ListenAndServe(":8080", router)
   }
   ```

   **router.go**  
   ```go
   func NewRouter() *http.ServeMux {
    r := http.NewServeMux()
	// home handlers
	r.HandleFunc("/", handlers.HomeHandler)
	r.HandleFunc("GET /err", handlers.ErrorHandler)
	// user handlers
	r.HandleFunc("GET /login", handlers.LoginHandler)
	r.HandleFunc("GET /logout", handlers.LogoutHandler)
	r.HandleFunc("GET /signup", handlers.SignupHandler)
	r.HandleFunc("POST /signup", handlers.SignupAccountHandler)
	r.HandleFunc("POST /authenticate", handlers.AuthenticateHandler)
	// thread handlers
	r.HandleFunc("GET /thread/list", handlers.ThreadListHandler)
	r.HandleFunc("POST /thread/create", handlers.CreateThreadHandler)
	r.HandleFunc("GET /thread/{id}", handlers.ThreadHandler)
	r.HandleFunc("POST /thread/post", handlers.CreatePostHandler)
    return r
   }
   ```

3. **HTMX CDN 추가 및 새 게시글 작성 폼 템플릿 생성:** templates/index.templ 파일
   ```html
   templ Index(posts []models.Post) {
       <!DOCTYPE html>
       <html>
       <head>
           <title>Chitchat 게시판</title>
           <script src="https://unpkg.com/htmx.org@2.0.4/dist/htmx.js" integrity="sha384-oeUn82QNXPuVkGCkcrInrS1twIxKhkZiFfr2TdiuObZ3n3yIeMiqcRzkIcguaof1" crossorigin="anonymous"></script>
       </head>
       <body>
           <h1>Chitchat 게시판</h1>
           <ul>
               for _, post := range posts {
                   <li>
                       <strong>{ post.Author }:</strong> { post.Content }
                   </li>
               }
           </ul>

           <form role="form" id="post-form">
               <div class="lead">Start a new thread with the following topic</div>
               <div class="form-group">
                   <textarea class="form-control" name="topic" id="topic" placeholder="Thread topic here" rows="4"></textarea>
                   <br/>
                   <br/>
                   <button class="btn btn-lg btn-primary pull-right" type="button" hx-post="/post/create" hx-target="#topic" hx-trigger="click" hx-swap="innerHTML">Start this thread</button>
               </div>
           </form>
       </body>
       </html>
   }
   ```

4. **templ 파일 컴파일 및 실행:** `templ generate` 명령어를 다시 실행하여 변경된 템플릿 파일들을 컴파일합니다. `go run main.go` 명령어를 실행하고, 웹 브라우저에서 http://localhost:8080 에 접속합니다. "새 게시글 작성" 버튼을 클릭하면 폼이 동적으로 나타나고, 작성 후 "작성 완료" 버튼을 누르면 페이지 새로고침 없이 게시글 목록이 업데이트되는 것을 확인할 수 있습니다.

### 5단계: 기능 확장 및 UI 개선 (점진적 개발)

- **데이터 모델 확장:** Post 모델에 CreatedAt, UpdatedAt 등의 필드를 추가하고, 필요에 따라 Thread, User, Comment 등의 모델을 추가합니다.  
  **데이터 구조:**
  - User—Representing the forum user's information
  - Session—Representing a user's current login session
  - Thread—Representing a forum thread (a conversation among forum users)
  - Post—Representing a post (a message added by a forum user) within a thread

| HTTP 메서드| 패스         | 개요                (model data)
|----------|------------|----------------------------------------------------------------
| GET      | `/`        | home("/") 페이지
| GET      | `/err`     | Error message 출력
| GET      | `/signup`  | 신규 사용자 등록 페이지   
| POST     | `/signup`  | 신규 사용자 등록 User(name,email,pswd), login 페이지로 redirect
| GET      | `/login`   | 기존 사용자 로그인 페이지  
| POST     | `/authenticate`| 기존 사용자 인증 User(email,pswd), home("/") 페이지로 redirect
| GET      | `/logout`  | 로그아웃 요청, home("/") 페이지로 redirect
| GET      | `/thread/new` | new thread 생성 페이지 보여주기
| POST     | `/thread/create` | Thread(topic) 생성
| GET      | `/thread/{thread_id}` | Show the details of the thread & posts, the form to write a post
| POST     | `/thread/post` | post 생성 Post(body)

- **CRUD 기능 확장:** 게시글 수정, 삭제 기능 추가
- **UI 개선:** CSS 스타일 적용, 더 나은 폼 디자인, 목록 디자인 개선  
- **유효성 검사:** 폼 입력 값에 대한 유효성 검사 추가 (서버 & 클라이언트)  
- **페이지네이션:** 게시글 목록 페이지네이션 기능 추가  
- **검색 기능:** 게시글 검색 기능 추가
- **사용자 인증/인가:** (추후) JWT 토큰 이용한 사용자 계정, 로그인/로그아웃, 권한 관리 기능 추가

### 6단계: PostgreSQL 및 pgx 드라이버로 전환

- **PostgreSQL 설치 및 설정:** Windows 에 PostgreSQL 을 설치하고, 데이터베이스 및 사용자 설정을 완료합니다.  
- **pgx 라이브러리 추가:** `go get github.com/jackc/pgx/v5/pgxpool` 명령어를 실행하여 pgx 라이브러리를 추가합니다.  
- **데이터베이스 연결 정보 변경:** main.go 파일에서 sqlx.Open 함수 부분을 PostgreSQL 연결 정보로 변경합니다. (connection string 은 PostgreSQL 설정에 따라 다름)  
  ```go
  // ... sqlite3 연결 코드 주석 처리 또는 제거 ...
  // db, err = sqlx.Open("sqlite3", ":memory:")

  dsn := "postgresql://사용자:비밀번호@localhost:5432/데이터베이스이름?sslmode=disable" // PostgreSQL DSN (Data Source Name)
  db, err = sqlx.Connect("pgx", dsn)
  if err != nil {
      log.Fatal(err)
  }
  ```

- **SQL 쿼리 수정:** SQLite3 와 PostgreSQL 은 SQL 문법이 약간 다를 수 있으므로, 테이블 생성 및 데이터 조회/삽입 쿼리를 PostgreSQL 문법에 맞게 수정합니다. (필요한 경우)  
- **실행 및 테스트:** `go run main.go` 명령어를 실행하고, 웹 애플리케이션이 PostgreSQL 데이터베이스에 정상적으로 연결되고 기능하는지 확인합니다.

### 7단계: 배포 (선택 사항)

- **실행 파일 빌드:** `go build -o chitchat.exe main.go` 명령어를 실행하여 실행 파일 chitchat.exe 를 빌드합니다.  
- **실행 파일 실행:** 빌드된 chitchat.exe 파일을 실행하여 웹 서버를 시작합니다.  
- **배포 환경 구축:** (선택 사항) AWS, Google Cloud, Azure 등의 클라우드 플랫폼 또는 Docker 를 이용하여 배포 환경을 구축하고, 애플리케이션을 배포합니다.

---

**주의 사항:**

- **에러 처리:** 각 단계별 코드에서 에러 처리를 꼼꼼하게 구현하여 안정적인 애플리케이션을 만드는 것이 중요합니다.  
- **보안:** 학생 학습 목적에서는 보안을 깊이 고려하지 않았지만, 실제 서비스 개발 시에는 보안 취약점을 항상 염두에 두어야 합니다. (특히 사용자 입력 값 검증, SQL Injection 방지, Cross-Site Scripting (XSS) 방지 등)  
- **코드 구조:** 규모가 커짐에 따라 코드를 여러 파일로 분리하고, 패키지 구조를 체계적으로 관리하는 것이 중요합니다. (예: handler, model, repository, service 패키지 등)  
- **테스트:** 단위 테스트, 통합 테스트 등을 작성하여 코드의 안정성을 확보하는 것이 중요합니다. (학습 단계에서는 생략 가능)

---

**학습 팁:**

- **단계별 커밋:** 각 단계를 완료할 때마다 Git 에 커밋하여 변경 사항을 기록하고 관리하는 습관을 들이세요.  
- **디버깅:** VS Code 디버깅 기능을 적극적으로 활용하여 코드 오류를 찾고 수정하는 연습을 하세요.  
- **Go 문서 및 온라인 자료 활용:** Go 공식 문서 ([https://go.dev/](https://go.dev/)), net/http 패키지 문서, sqlx 문서, templ 문서, HTMX 문서 등을 참고하고, Stack Overflow, Go Forum 등의 온라인 커뮤니티에서 정보를 얻으세요.  
- **꾸준한 학습:** 웹 개발은 꾸준한 학습과 실습이 중요합니다. chitchat 서비스 개발 외에도 다양한 프로젝트를 만들어 보면서 실력을 향상시켜 보세요.

---

**Sources**  
1. [https://github.com/sausheong/gwp](https://github.com/sausheong/gwp)
