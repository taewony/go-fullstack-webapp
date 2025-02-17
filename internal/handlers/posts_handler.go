package handlers

import (
	"fmt"
	"net/http"
)

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PostsHandler")
	// templ.Posts().Render(r.Context(), w)
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreatePostHandler")

	// Output the specific topic
	if topic := r.FormValue("topic"); topic != "" {
		fmt.Printf("Topic: %s\n", topic)
	}

	// templ.CreatePost().Render(r.Context(), w)
}
