package handlers

import (
	"net/http"
)

// ServeStatic обслуговує статичні файли
func ServeStatic() http.Handler {
	return http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static")))
}
