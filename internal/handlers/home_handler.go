package handlers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/taewony/go-fullstack-webapp/internal/components"
	"github.com/taewony/go-fullstack-webapp/internal/models"
)

// GET /err?msg=
// shows the error message page
func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	encodedMsg := url.QueryEscape(vals.Get("msg"))
	if r.Header.Get("HX-Request") == "true" {
		components.ErrorTempl(encodedMsg).Render(r.Context(), w)
	} else {
		_, err := session(w, r)
		if err != nil {
			generateHTML(w, encodedMsg, "layout", "public.navbar", "error")
		} else {
			generateHTML(w, encodedMsg, "layout", "private.navbar", "error")
		}
	}
}

func IndexHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Header.Get("HX-Request") == "true" {
		fmt.Println("HX-Request")
	} else {
		fmt.Println("HTML-Request")
	}
	threads, err := models.Threads()
	if err != nil {
		error_message(writer, request, "Cannot get threads")
		return
	}
	_, err = session(writer, request)
	if err != nil {
		components.LayoutTempl(components.PublicNavbarTempl(), threads).Render(request.Context(), writer)
	} else {
		components.LayoutTempl(components.PrivateNavbarTempl(), threads).Render(request.Context(), writer)
	}
}

func HomeHandler(writer http.ResponseWriter, request *http.Request) {
	threads, err := models.Threads()
	if err != nil {
		error_message(writer, request, "Cannot get threads")
		return
	}
	_, err = session(writer, request)
	if err != nil {
		generateHTML(writer, threads, "layout", "public.navbar", "index")
	} else {
		generateHTML(writer, threads, "layout", "private.navbar", "index")
	}
}
