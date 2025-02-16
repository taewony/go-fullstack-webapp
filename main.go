package main

import (
	"fmt"
	"log"
	"net/http"

	models "github.com/taewony/go-fullstack-webapp/internal/models"
	templates "github.com/taewony/go-fullstack-webapp/internal/templates" // templates 패키지 import (go.mod 파일 경로에 따라 수정)

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
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

func getPosts() ([]models.Post, error) {
	var posts []models.Post
	err := db.Select(&posts, "SELECT id, content, author FROM posts") // sqlx.DB.Select 사용하여 복수 row 조회
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func main() {
	var err error
	db, err = sqlx.Open("sqlite", ":memory:") // Change "sqlite3" to "sqlite"
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
