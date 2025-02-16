package main

import (
	"fmt"
	"net/http"

	templates "github.com/taewony/go-fullstack-webapp/internal/templates" // templates 패키지 import (go.mod 파일 경로에 따라 수정)
)

func handler(w http.ResponseWriter, r *http.Request) {
	templates.Index().Render(r.Context(), w) // index.templ 렌더링
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("chitchat 서버 시작... http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
