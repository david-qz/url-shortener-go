package main

import (
	"fmt"
	"log"
	"net/http"

	"example.com/go-api/urls"
	"github.com/akrylysov/algnhsa"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/l/{key}", linkHandler)
	algnhsa.ListenAndServe(r, nil)
}

func linkHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	url, err := urls.GetUrl(key, r.Context())

	if err != nil {
		log.Println(err)
		sendErrorResponse(w, err)
		return
	}

	if url == nil {
		sendNotFoundResponse(w)
		return
	}

	http.Redirect(w, r, *url, 302)
}

func sendNotFoundResponse(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte("<h1>Not found</h1>\n"))
}

func sendErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf("<h1>Server error</h1><br><pre>%s</pre>\n", err.Error())))
}
