package forum

import "net/http"

func isAuthenticated(r *http.Request) bool {
    _, err := r.Cookie("session_id")
    return err == nil
}
