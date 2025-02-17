package handlers

import (
	"fmt"
	"net/http"
)

func ThreadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ThreadHandler")
	// templ.Thread().Render(r.Context(), w)
}

func CreateThreadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateThreadHandler")
	// templ.CreateThread().Render(r.Context(), w)
}
