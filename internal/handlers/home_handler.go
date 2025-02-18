package handlers

import (
	"net/http"

	"github.com/taewony/go-fullstack-webapp/internal/models"
)

// GET /err?msg=
// shows the error message page
func ErrorHandler(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	_, err := session(writer, request)
	if err != nil {
		generateHTML(writer, vals.Get("msg"), "layout", "public.navbar", "error")
	} else {
		generateHTML(writer, vals.Get("msg"), "layout", "private.navbar", "error")
	}
}

func HomeHandler(writer http.ResponseWriter, request *http.Request) {
	threads, err := models.Threads()
	if err != nil {
		error_message(writer, request, "Cannot get threads")
	} else {
		_, err := session(writer, request)
		if err != nil {
			generateHTML(writer, threads, "layout", "public.navbar", "index")
		} else {
			generateHTML(writer, threads, "layout", "private.navbar", "index")
		}
	}
}

func IndexHandler(writer http.ResponseWriter, request *http.Request) {
	threads, err := models.Threads()
	if err != nil {
		error_message(writer, request, "Cannot get threads")
	} else {
		_, err := session(writer, request)
		if err != nil {
			generateHTML(writer, threads, "layout", "public.navbar", "index") // logged out state
			// templ.Layout(threads, false)
		} else {
			generateHTML(writer, threads, "layout", "private.navbar", "index") // logged in state
		}
	}
}
