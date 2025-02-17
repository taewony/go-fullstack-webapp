package handlers

import (
	"fmt"
	"net/http"

	"github.com/taewony/go-fullstack-webapp/internal/models"
	templ "github.com/taewony/go-fullstack-webapp/internal/templates"
)

func ThreadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ThreadHandler")

	threads, _ := models.Threads()
	templ.Content(threads).Render(r.Context(), w)
}

func CreateThreadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateThreadHandler")
	// templ.CreateThread().Render(r.Context(), w)
}
