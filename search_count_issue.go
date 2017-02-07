package sample

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// setup handlers
func init() {
	http.HandleFunc("/api/search", handleSearch)
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	log.Infof(c, "Hello!")
}
