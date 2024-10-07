package forum

import (
	"log"
	"net/http"
	"text/template"
)

func RenderErrorPage(w http.ResponseWriter, r *http.Request) {
	renderErrorTemplate(w, r.Response.StatusCode)
}

// loades an error page from all aspects of the project but is an int rather than a http request
func renderErrorTemplate(w http.ResponseWriter, statusCode int) {
	errorTemplate, err := template.ParseFiles("templates/error.html")
	if err != nil {
		log.Printf("Error parsing error template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// we placed the following to automatically access the error page using one withoit having to manually right all of it
	errorPage := ErrorPage{
		Title:       http.StatusText(statusCode),
		Code:        statusCode,
		Description: StatusDescription(statusCode),
	}
	// executes the error pages
	w.WriteHeader(statusCode)
	if err := errorTemplate.Execute(w, errorPage); err != nil {
		log.Printf("Error executing error template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}