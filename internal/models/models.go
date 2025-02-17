package models

import (
	"log"

	"github.com/jmoiron/sqlx"
)

// Post represents a blog post.
type Post struct {
	ID      int    `db:"id"`
	Content string `db:"content"`
	Author  string `db:"author"`
}

var Db *sqlx.DB

func InitDB() {
	var err error
	Db, err = sqlx.Open("sqlite", ":memory:") // Change "sqlite3" to "sqlite"
	if err != nil {
		log.Fatal(err)
	}

	_, err = Db.Exec(`
        CREATE TABLE IF NOT EXISTS posts (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            content TEXT NOT NULL,
            author TEXT NOT NULL
        );
    `) // posts 테이블 생성
	if err != nil {
		log.Fatal(err)
	}

	_, err = Db.Exec(`
        INSERT INTO posts (content, author) VALUES (?, ?), (?, ?)
    `, "첫 번째 게시글 내용입니다.", "John Doe", "두 번째 게시글입니다.", "Jane Doe") // 샘플 데이터 삽입
	if err != nil {
		log.Fatal(err)
	}
}

func GetPosts() ([]Post, error) {
	var posts []Post
	err := Db.Select(&posts, "SELECT id, content, author FROM posts") // sqlx.DB.Select 사용하여 복수 row 조회
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return posts, nil
}
