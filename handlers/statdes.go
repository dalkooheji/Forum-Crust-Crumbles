package forum

import (
	
)
func StatusDescription(code int) string {
	switch code {
	case 400:
		return "Oops! BAD REQUEST."
	case 404:
		return "Oops! page does not exist."
	case 405:
		return "Oops! not alloweeddd ."
	case 418:
		return "Oops! I AM A TEAPOT."
	case 500:
		return "Oops! Something went wrong with the server"
}
return ""
}