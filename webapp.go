/**
 *  ------------------------------------------------------------
 *  @project	webapp
 *  @file       webapp.go
 *  @date       2014-10-16
 *  @author     Jim Zhan <jim.zhan@me.com>
 *
 *  Copyright © 2014 Jim Zhan.
 *  ------------------------------------------------------------
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *  ------------------------------------------------------------
 */
package webapp

import (
	"net/http"
	"path"
	"runtime"
	"strings"

	"github.com/gorilla/mux"
)

var (
	cwd    string
	here   string
	logger = GetLogger("webapp")

	Settings *Config
)

type (
	Application struct {
		*mux.Router
	}

	HTTPRequest interface {
		GET(string, Handler)
		POST(string, Handler)
		PUT(string, Handler)
		DELETE(string, Handler)
		PATCH(string, Handler)
		HEAD(string, Handler)
		OPTIONS(string, Handler)
	}
)

// Initialize application settings & basic environmetal variables.
func init() {
	_, filename, _, _ := runtime.Caller(1)
	here = path.Dir(filename)
	Settings = Configure("app")
}

// New creates a new webapp instance.
func New() *Application {
	return &Application{mux.NewRouter()}
}

// ---------------------------------------------------------------------------
//  HTTP Requests Handlers
// ---------------------------------------------------------------------------
// GET is a shortcut for app.HandleFunc(pattern, handler).Methods("GET")
func (app *Application) GET(pattern string, handler http.HandlerFunc) {
	app.HandleFunc(pattern, handler).Methods("GET")
}

// POST is a shortcut for app.HandleFunc(pattern, handler).Methods("POST")
func (app *Application) POST(pattern string, handler http.HandlerFunc) {
	app.HandleFunc(pattern, handler).Methods("POST")
}

// PUT is a shortcut for app.HandleFunc(pattern, handler).Methods("PUT")
func (app *Application) PUT(pattern string, handler http.HandlerFunc) {
	app.HandleFunc(pattern, handler).Methods("PUT")
}

// DELETE is a shortcut for app.HandleFunc(pattern, handler).Methods("DELETE")
func (app *Application) DELETE(pattern string, handler http.HandlerFunc) {
	app.HandleFunc(pattern, handler).Methods("DELETE")
}

// PATCH is a shortcut for app.HandleFunc(pattern, handler).Methods("PATCH")
func (app *Application) PATCH(pattern string, handler http.HandlerFunc) {
	app.HandleFunc(pattern, handler).Methods("PATCH")
}

// HEAD is a shortcut for app.HandleFunc(pattern, handler).Methods("HEAD")
func (app *Application) HEAD(pattern string, handler http.HandlerFunc) {
	app.HandleFunc(pattern, handler).Methods("HEAD")
}

// OPTIONS is a shortcut for app.HandleFunc(pattern, handler).Methods("OPTIONS")
func (app *Application) OPTIONS(pattern string, handler http.HandlerFunc) {
	app.HandleFunc(pattern, handler).Methods("OPTIONS")
}

// Group creates a new application group under the given prefix.
func (app *Application) Group(prefix string) *Application {
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	return &Application{app.PathPrefix(prefix).Subrouter()}
}

// ---------------------------------------------------------------------------
//  HTTP Server
// ---------------------------------------------------------------------------
// Serve starts serving the requests at the pre-defined address from application settings file.
func (app *Application) Serve() {
	address := Settings.GetString("address")
	logger.Info("Server started [" + address + "]")
	err := http.ListenAndServe(address, app)
	if nil != err {
		panic(err)
	}
}
