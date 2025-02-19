package handlers

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/taewony/go-fullstack-webapp/internal/models"
)

var logger *log.Logger

func init() {
	file, err := os.OpenFile("chitchat.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}
	logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Convenience function to redirect to the error message page
func error_message(w http.ResponseWriter, r *http.Request, msg string) {
	encodedMsg := "/err?msg=" + url.QueryEscape(msg)
	// HTMX 요청인지 확인 (HX-Request 헤더 체크)
	if r.Header.Get("HX-Request") == "true" {
		// HTMX 요청인 경우: HX-Retarget 설정 후 에러 페이지로 리다이렉트
		w.Header().Set("HX-Retarget", "#container") // HX-Retarget 헤더 설정: 타겟 엘리먼트 지정

	}
	// 일반적인 HTTP 요청 (HTMX 요청이 아닌 경우): 에러 페이지로 리다이렉트 (기존 방식 유지)
	// http.StatusFound (302)는 임시 리다이렉트에 사용합니다.
	http.Redirect(w, r, encodedMsg, http.StatusFound)
}

// Checks if the user is logged in and has a session, if not err is not nil
func session(writer http.ResponseWriter, request *http.Request) (sess models.Session, err error) {
	cookie, err := request.Cookie("_cookie")
	if err == nil {
		sess = models.Session{Uuid: cookie.Value}
		if ok, _ := sess.Check(); !ok {
			err = errors.New("Invalid session")
		}
	}
	return
}

// parse HTML templates
// pass in a list of file names, and get a template
func parseTemplateFiles(filenames ...string) (t *template.Template) {
	var files []string
	t = template.New("layout")
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	t = template.Must(t.ParseFiles(files...))
	return
}

// generate HTML from a text/template
func generateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(writer, "layout", data)
}

// for logging
func info(args ...interface{}) {
	logger.SetPrefix("INFO ")
	logger.Println(args...)
}

func danger(args ...interface{}) {
	logger.SetPrefix("ERROR ")
	logger.Println(args...)
}

func warning(args ...interface{}) {
	logger.SetPrefix("WARNING ")
	logger.Println(args...)
}
