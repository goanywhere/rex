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
	"reflect"
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
	AbstractRequest interface {
		GET(string, http.HandlerFunc)
		POST(string, http.HandlerFunc)
		PUT(string, http.HandlerFunc)
		DELETE(string, http.HandlerFunc)
		PATCH(string, http.HandlerFunc)
		HEAD(string, http.HandlerFunc)
		OPTIONS(string, http.HandlerFunc)
	}

	Application struct {
		*mux.Router
	}
)

// Initialize application settings & basic environmetal variables.
func init() {
	here = path.Dir(getCurrentFile())
	Settings = Configure("app")
}

// New creates a new webapp instance.
func New() *Application {
	return &Application{mux.NewRouter()}
}

// getCurrentFile finds current working file with full path.
func getCurrentFile() string {
	_, filename, _, _ := runtime.Caller(1)
	return filename
}

// getFuncName finds the full function name (with package).
func getFuncName(function interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name()
}

// ---------------------------------------------------------------------------
//  HTTP Requests Handlers
// ---------------------------------------------------------------------------
// GET is a shortcut for app.HandleFunc(pattern, handler).Methods("GET"),
// it also fetch the full function name of the handler (with package) to name the route.
func (app *Application) GET(pattern string, handler http.HandlerFunc) {
	app.HandleFunc(pattern, handler).Methods("GET").Name(getFuncName(handler))
}

// POST is a shortcut for app.HandleFunc(pattern, handler).Methods("POST")
// it also fetch the full function name of the handler (with package) to name the route.
func (app *Application) POST(pattern string, handler http.HandlerFunc) {
	app.HandleFunc(pattern, handler).Methods("POST").Name(getFuncName(handler))
}

// PUT is a shortcut for app.HandleFunc(pattern, handler).Methods("PUT")
// it also fetch the full function name of the handler (with package) to name the route.
func (app *Application) PUT(pattern string, handler http.HandlerFunc) {
	app.HandleFunc(pattern, handler).Methods("PUT").Name(getFuncName(handler))
}

// DELETE is a shortcut for app.HandleFunc(pattern, handler).Methods("DELETE")
// it also fetch the full function name of the handler (with package) to name the route.
func (app *Application) DELETE(pattern string, handler http.HandlerFunc) {
	app.HandleFunc(pattern, handler).Methods("DELETE").Name(getFuncName(handler))
}

// PATCH is a shortcut for app.HandleFunc(pattern, handler).Methods("PATCH")
// it also fetch the full function name of the handler (with package) to name the route.
func (app *Application) PATCH(pattern string, handler http.HandlerFunc) {
	app.HandleFunc(pattern, handler).Methods("PATCH").Name(getFuncName(handler))
}

// HEAD is a shortcut for app.HandleFunc(pattern, handler).Methods("HEAD")
// it also fetch the full function name of the handler (with package) to name the route.
func (app *Application) HEAD(pattern string, handler http.HandlerFunc) {
	app.HandleFunc(pattern, handler).Methods("HEAD").Name(getFuncName(handler))
}

// OPTIONS is a shortcut for app.HandleFunc(pattern, handler).Methods("OPTIONS")
// it also fetch the full function name of the handler (with package) to name the route.
func (app *Application) OPTIONS(pattern string, handler http.HandlerFunc) {
	app.HandleFunc(pattern, handler).Methods("OPTIONS").Name(getFuncName(handler))
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
	if err := http.ListenAndServe(address, app); err != nil {
		panic(err)
	}
}
