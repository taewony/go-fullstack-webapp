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
