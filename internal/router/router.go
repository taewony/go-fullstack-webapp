package router

import (
	"net/http"

	handlers "github.com/taewony/go-fullstack-webapp/internal/handlers"
)

func NewRouter() *http.ServeMux {
	r := http.NewServeMux()

	r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))

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
