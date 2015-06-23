package main

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/goanywhere/rex"
	"github.com/goanywhere/rex/livereload"
)

type User struct {
	Username string
}

func Index(w http.ResponseWriter, r *http.Request) {
	if html, err := template.ParseFiles(filepath.Join("views", "index.html")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "text/html")
		var user = User{Username: rex.Env.String("USER", "guest")}
		html.Execute(w, user)
	}
}

func main() {
	rex.Use(livereload.Middleware)
	rex.Get("/", Index)
	rex.Run()
}
