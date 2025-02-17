package router

import (
	"net/http"

	handlers "github.com/taewony/go-fullstack-webapp/internal/handlers"
)

// Define an error handler function
func errorHandler(writer http.ResponseWriter, request *http.Request) {
	// Handle the error response here
	http.Error(writer, "An error occurred", http.StatusInternalServerError)
}

func NewRouter() *http.ServeMux {
	r := http.NewServeMux()

	r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("/public"))))

	r.HandleFunc("/", handlers.HomeHandler)
	r.HandleFunc("GET /err", errorHandler)

	// r.HandleFunc("GET /login", handlers.LoginHandler)
	// r.HandleFunc("GET /logout", handlers.LogoutHandler)
	// r.HandleFunc("GET /signup", handlers.SignupHandler)
	// r.HandleFunc("POST /signup_account", handlers.SignupAccountHandler)
	// r.HandleFunc("POST /authenticate", handlers.AuthenticateHandler)

	r.HandleFunc("GET /thread/", handlers.ThreadHandler)
	r.HandleFunc("GET /thread/read", handlers.ThreadHandler)
	// r.HandleFunc("GET /thread/new", handlers.NewThreadHandler)
	r.HandleFunc("POST /thread/create", handlers.CreateThreadHandler)

	r.HandleFunc("POST /post/create", handlers.CreatePostHandler)
	return r
}
