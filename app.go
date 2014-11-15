/**
 *  ------------------------------------------------------------
 *  @project	web.go
 *  @file       app.go
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
package web

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/gorilla/mux"
)

type (
	Application struct {
		router      *mux.Router
		middlewares []Middleware
	}

	HandlerFunc func(*Context)

	// Conventional method to implement custom middlewares.
	Middleware func(http.Handler) http.Handler

	// Shortcut to create map.
	H map[string]interface{}
)

// New creates an application instance & setup its default settings..
func New() *Application {
	app := &Application{mux.NewRouter(), nil}
	return app
}

// ---------------------------------------------------------------------------
//  Custom handler func with Context Supports
// ---------------------------------------------------------------------------
func (self HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	self(NewContext(w, r))
}

// ---------------------------------------------------------------------------
//  HTTP Requests Handlers
// ---------------------------------------------------------------------------
// Supported Handler Types
//	* http.Handler
//	* http.HandlerFunc	=> func(w http.ResponseWriter, r *http.Request)
//	* web.HandlerFunc	=> func(ctx *Context)
func (self *Application) handle(method, pattern string, h interface{}) {
	var handler http.Handler

	switch h.(type) {
	// Standard net/http.Handler/HandlerFunc
	case http.Handler:
		handler = h.(http.Handler)
	case func(w http.ResponseWriter, r *http.Request):
		handler = http.HandlerFunc(h.(func(w http.ResponseWriter, r *http.Request)))
	case func(ctx *Context):
		handler = HandlerFunc(h.(func(ctx *Context)))
	default:
		panic(fmt.Sprintf("Unknown handler type (%v) passed in.", handler))
	}
	// finds the full function name (with package)
	name := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	self.router.Handle(pattern, handler).Methods(method).Name(name)
}

// GET is a shortcut for app.HandleFunc(pattern, handler).Methods("GET"),
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Application) Get(pattern string, handler interface{}) {
	self.handle("GET", pattern, handler)
}

// POST is a shortcut for app.HandleFunc(pattern, handler).Methods("POST")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Application) Post(pattern string, handler interface{}) {
	self.handle("POST", pattern, handler)
}

// PUT is a shortcut for app.HandleFunc(pattern, handler).Methods("PUT")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Application) Put(pattern string, handler interface{}) {
	self.handle("PUT", pattern, handler)
}

// DELETE is a shortcut for app.HandleFunc(pattern, handler).Methods("DELETE")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Application) Delete(pattern string, handler interface{}) {
	self.handle("DELETE", pattern, handler)
}

// PATCH is a shortcut for app.HandleFunc(pattern, handler).Methods("PATCH")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Application) Patch(pattern string, handler http.HandlerFunc) {
	self.handle("PATCH", pattern, handler)
}

// HEAD is a shortcut for app.HandleFunc(pattern, handler).Methods("HEAD")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Application) Head(pattern string, handler http.HandlerFunc) {
	self.handle("HEAD", pattern, handler)
}

// OPTIONS is a shortcut for app.HandleFunc(pattern, handler).Methods("OPTIONS")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Application) Options(pattern string, handler http.HandlerFunc) {
	self.handle("OPTIONS", pattern, handler)
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
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// Serve starts serving the requests at the pre-defined address from settings.
// TODO command line arguments.
func (self *Application) Serve() {
	address := fmt.Sprintf("%s:%s", Env.Get("host"), Env.Get("port"))
	Info("Application server started [%s]", address)
	if err := http.ListenAndServe(address, self); err != nil {
		panic(err)
	}
}

func init() {
	// Application Defaults
	var root string
	if cwd, err := os.Getwd(); err == nil {
		root, _ = filepath.Abs(cwd)
	} else {
		panic(err)
	}
	Env.Set("root", root)
	Env.Set("host", "0.0.0.0")
	Env.Set("port", "5000")
	Env.Set("templates", "templates")
}
