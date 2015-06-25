package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/goanywhere/env"
	"github.com/goanywhere/rex"
	"github.com/goanywhere/rex/livereload"
)

type User struct {
	Username string
}

// -------------------- HTML Template --------------------
func Index(w http.ResponseWriter, r *http.Request) {
	if html, err := template.ParseFiles(filepath.Join("views", "index.html")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "text/html")
		var user = User{Username: env.String("USER", "guest")}
		html.Execute(w, user)
	}
}

// -------------------- JSON Template --------------------
type Response struct {
	Code int
	Text string
}

func JSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func fetch(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	var buffer = new(bytes.Buffer)
	var response = Response{Code: http.StatusOK, Text: http.StatusText(http.StatusOK)}
	if err := json.NewEncoder(buffer).Encode(response); err == nil {
		w.Write(buffer.Bytes())
	} else {
		http.Error(w, err.Error(), http.StatusBadGateway)
	}
}

func create(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	var buffer = new(bytes.Buffer)
	var response = Response{Code: http.StatusCreated, Text: http.StatusText(http.StatusCreated)}
	if err := json.NewEncoder(buffer).Encode(response); err == nil {
		w.Write(buffer.Bytes())
	} else {
		http.Error(w, err.Error(), http.StatusBadGateway)
	}
}

func update(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	var buffer = new(bytes.Buffer)
	var response = Response{Code: http.StatusAccepted, Text: http.StatusText(http.StatusAccepted)}
	if err := json.NewEncoder(buffer).Encode(response); err == nil {
		w.Write(buffer.Bytes())
	} else {
		http.Error(w, err.Error(), http.StatusBadGateway)
	}
}

func remove(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusGone)
	var buffer = new(bytes.Buffer)
	var response = Response{Code: http.StatusGone, Text: http.StatusText(http.StatusGone)}
	if err := json.NewEncoder(buffer).Encode(response); err == nil {
		w.Write(buffer.Bytes())
	} else {
		http.Error(w, err.Error(), http.StatusBadGateway)
	}
}

func main() {
	app := rex.New()
	app.Use(livereload.Middleware)
	app.Get("/", Index)

	api := app.Group("/v1/")
	api.Use(JSON)
	api.Get("/", fetch)
	api.Post("/", create)
	api.Put("/", update)
	api.Delete("/", remove)

	app.Run()
}
