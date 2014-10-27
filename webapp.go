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
	"os"

	"github.com/gorilla/mux"
)

const ContentType = "Content-Type"

var (
	root     string
	Logger   = GetLogger("webapp")
	Settings *config
)

type (
	Application struct {
		router      *mux.Router
		middlewares []Middleware
	}

	// Conventional method to implement custom middlewares.
	Middleware func(http.Handler) http.Handler

	// Shortcut to create map.
	H map[string]interface{}
)

// Initialize application settings & basic environmetal variables.
func init() {
	root, _ = os.Getwd()
	Settings = configure("app")
}

// New creates a new webapp instance.
func New() *Application {
	return &Application{mux.NewRouter(), nil}
}

// ---------------------------------------------------------------------------
//  HTTP Requests Handlers
// ---------------------------------------------------------------------------
// GET is a shortcut for app.HandleFunc(pattern, handler).Methods("GET"),
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Application) GET(pattern string, handler http.HandlerFunc) {
	self.router.HandleFunc(pattern, handler).Methods("GET").Name(getFuncName(handler))
}

// POST is a shortcut for app.HandleFunc(pattern, handler).Methods("POST")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Application) POST(pattern string, handler http.HandlerFunc) {
	self.router.HandleFunc(pattern, handler).Methods("POST").Name(getFuncName(handler))
}

// PUT is a shortcut for app.HandleFunc(pattern, handler).Methods("PUT")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Application) PUT(pattern string, handler http.HandlerFunc) {
	self.router.HandleFunc(pattern, handler).Methods("PUT").Name(getFuncName(handler))
}

// DELETE is a shortcut for app.HandleFunc(pattern, handler).Methods("DELETE")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Application) DELETE(pattern string, handler http.HandlerFunc) {
	self.router.HandleFunc(pattern, handler).Methods("DELETE").Name(getFuncName(handler))
}

// PATCH is a shortcut for app.HandleFunc(pattern, handler).Methods("PATCH")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Application) PATCH(pattern string, handler http.HandlerFunc) {
	self.router.HandleFunc(pattern, handler).Methods("PATCH").Name(getFuncName(handler))
}

// HEAD is a shortcut for app.HandleFunc(pattern, handler).Methods("HEAD")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Application) HEAD(pattern string, handler http.HandlerFunc) {
	self.router.HandleFunc(pattern, handler).Methods("HEAD").Name(getFuncName(handler))
}

// OPTIONS is a shortcut for app.HandleFunc(pattern, handler).Methods("OPTIONS")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Application) OPTIONS(pattern string, handler http.HandlerFunc) {
	self.router.HandleFunc(pattern, handler).Methods("OPTIONS").Name(getFuncName(handler))
}

// HandleFunc is a shourtcut to router's HandleFunc with multiple methods supports,
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Application) HandleFunc(pattern string, handler http.HandlerFunc, methods ...string) {
	self.router.HandleFunc(pattern, handler).Methods(methods...).Name(getFuncName(handler))
}

// Group creates a new application group under the given path.
func (self *Application) Group(path string) *Application {
	return &Application{self.router.PathPrefix(path).Subrouter(), nil}
}

// ---------------------------------------------------------------------------
//  HTTP Server with Middleware Supports
// ---------------------------------------------------------------------------
func (self *Application) Use(middlewares ...Middleware) {
	self.middlewares = append(self.middlewares, middlewares...)
}

// ServeHTTP turn Application into http.Handler by implementing the http.Handler interface.
func (self *Application) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var app http.Handler = self.router
	// Activate middlewares in FIFO order.
	if len(self.middlewares) > 0 {
		for index := len(self.middlewares) - 1; index >= 0; index-- {
			app = self.middlewares[index](app)
		}
	}
	app.ServeHTTP(writer, request)
}

// Serve starts serving the requests at the pre-defined address from application settings file.
// TODO command line arguments.
func (self *Application) Serve() {
	Logger.Info("Application server started [" + Settings.GetString("address") + "]")
	if err := http.ListenAndServe(Settings.GetString("address"), self); err != nil {
		panic(err)
	}
}
