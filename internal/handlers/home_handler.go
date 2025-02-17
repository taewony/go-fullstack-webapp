package handlers

import (
	"net/http"

	"github.com/taewony/go-fullstack-webapp/internal/models"
	templ "github.com/taewony/go-fullstack-webapp/internal/templates"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := models.Posts()
	if err != nil {
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}

	templ.Index(posts).Render(r.Context(), w)
}
