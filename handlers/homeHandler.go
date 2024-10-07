package forum

import (
	"html/template"
	"net/http"
	"log"
	"strings"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	username, err := FetchUsernameByID(session.UserID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, map[string]string{
		"Username": username,
	})
}

type HomePageData struct {
	Username string
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	// Trim the path
	path := strings.TrimSpace(r.URL.Path)
	if path != "/" {
		renderErrorTemplate(w, http.StatusNotFound)
		return
	}

	// Only allow GET requests
	if r.Method != http.MethodGet {
		renderErrorTemplate(w, http.StatusMethodNotAllowed)
		return
	}

	// Attempt to get the current session
	session, err := GetSession(r)
	var username string
	if err == nil {
		// If a session is found, fetch the username
		username, err = FetchUsernameByID(session.UserID)
		if err != nil {
			log.Printf("Error fetching username: %v", err)
			renderErrorTemplate(w, http.StatusInternalServerError)
			return
		}
	}

	// Load the index template
	mainTemplate, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Error parsing index template: %v", err)
		renderErrorTemplate(w, http.StatusInternalServerError)
		return
	}

	// Data to pass to the template
	data := HomePageData{
		Username: username,
	}

	// Render the template with data
	if err := mainTemplate.Execute(w, data); err != nil {
		log.Printf("Error executing index template: %v", err)
		renderErrorTemplate(w, http.StatusInternalServerError)
		return
	}
}