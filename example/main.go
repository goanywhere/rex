package main

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/goanywhere/rex"
	"github.com/goanywhere/rex/livereload"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if html, err := template.ParseFiles(filepath.Join("views", "index.html")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		html.Execute(w, nil)
	}
}

func main() {
	rex.Use(livereload.Middleware)
	rex.Get("/", Index)
	rex.Run()
}
