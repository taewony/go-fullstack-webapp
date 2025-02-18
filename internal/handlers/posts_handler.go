package handlers

import (
	"fmt"
	"net/http"

	"github.com/taewony/go-fullstack-webapp/internal/models"
)

// POST /thread/post : Create the post based on form data {body, uuid, }
func CreatePostHandler(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
	} else {
		err = request.ParseForm()
		if err != nil {
			danger(err, "Cannot parse form")
		}
		user, err := sess.User()
		if err != nil {
			danger(err, "Cannot get user from session")
		}
		body := request.PostFormValue("body")
		uuid := request.PostFormValue("uuid")
		thread, err := models.ThreadByUUID(uuid)
		if err != nil {
			error_message(writer, request, "Cannot read thread")
		}
		if _, err := user.CreatePost(thread, body); err != nil {
			danger(err, "Cannot create post")
		}
		url := fmt.Sprintf("/thread/%s", uuid)
		http.Redirect(writer, request, url, 302)
	}
}
