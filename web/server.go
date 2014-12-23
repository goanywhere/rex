/* ----------------------------------------------------------------------
 *       ______      ___                         __
 *      / ____/___  /   |  ____  __  ___      __/ /_  ___  ________
 *     / / __/ __ \/ /| | / __ \/ / / / | /| / / __ \/ _ \/ ___/ _ \
 *    / /_/ / /_/ / ___ |/ / / / /_/ /| |/ |/ / / / /  __/ /  /  __/
 *    \____/\____/_/  |_/_/ /_/\__. / |__/|__/_/ /_/\___/_/   \___/
 *                            /____/
 *
 * (C) Copyright 2014 GoAnywhere (http://goanywhere.io).
 * ----------------------------------------------------------------------
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
 * ----------------------------------------------------------------------*/

package web

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"

	"github.com/gorilla/mux"
)

type (
	Server struct {
		router      *mux.Router
		middlewares []Middleware
	}

	HandlerFunc func(*Context)

	// Conventional method to implement custom middlewares.
	Middleware func(http.Handler) http.Handler
)

// New creates an application instance & setup its default settings..
func NewServer() *Server {
	return &Server{mux.NewRouter(), nil}
}

// Custom handler func provides Context Supports.
func (self HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	self(NewContext(w, r))
}

// ---------------------------------------------------------------------------
//  HTTP Requests Handlers
// ---------------------------------------------------------------------------
// Supported Handler Types
//	* http.Handler
//	* http.HandlerFunc	=> func(w http.ResponseWriter, r *http.Request)
//	* rex.HandlerFunc	=> func(ctx *Context)
func (self *Server) handle(method, pattern string, h interface{}) {
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
		log.Fatalf("Unknown handler type (%v) passed in.", handler)
	}
	// finds the full function name (with package)
	name := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	self.router.Handle(pattern, handler).Methods(method).Name(name)
}

// Get is a shortcut for mux.HandleFunc(pattern, handler).Methods("GET"),
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Server) Get(pattern string, handler interface{}) {
	self.handle("GET", pattern, handler)
}

// Post is a shortcut for mux.HandleFunc(pattern, handler).Methods("POST")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Server) Post(pattern string, handler interface{}) {
	self.handle("POST", pattern, handler)
}

// Put is a shortcut for mux.HandleFunc(pattern, handler).Methods("PUT")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Server) Put(pattern string, handler interface{}) {
	self.handle("PUT", pattern, handler)
}

// Delete is a shortcut for mux.HandleFunc(pattern, handler).Methods("DELETE")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Server) Delete(pattern string, handler interface{}) {
	self.handle("DELETE", pattern, handler)
}

// Patch is a shortcut for mux.HandleFunc(pattern, handler).Methods("PATCH")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Server) Patch(pattern string, handler http.HandlerFunc) {
	self.handle("PATCH", pattern, handler)
}

// Head is a shortcut for mux.HandleFunc(pattern, handler).Methods("HEAD")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Server) Head(pattern string, handler http.HandlerFunc) {
	self.handle("HEAD", pattern, handler)
}

// Options is a shortcut for mux.HandleFunc(pattern, handler).Methods("OPTIONS")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Server) Options(pattern string, handler http.HandlerFunc) {
	self.handle("OPTIONS", pattern, handler)
}

// Group creates a new application group under the given path.
func (self *Server) Group(path string) *Server {
	return &Server{self.router.PathPrefix(path).Subrouter(), nil}
}

// Livereload provides websocket-based livereload supports for browser.
// SEE http://feedback.livereload.com/knowledgebase/articles/86174-livereload-protocol
func (self *Server) Livereload() {
	/*
		self.Get("/livereload.js", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/javascript")
			w.Write(livereloadjs)
		})
		self.Get("/livereload", Livereload)
	*/
}

// Use append middleware into the serving list, middleware will be served in FIFO order.
func (self *Server) Use(middlewares ...Middleware) {
	self.middlewares = append(self.middlewares, middlewares...)
}

// ServeHTTP: Implementation of "http.Handler" interface.
func (self *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if settings.Debug {
		Livereload.Start()
		self.Get("/livereload", Livereload.Serve)
		self.Get("/livereload.js", Livereload.ServeJS)
	}

	var mux http.Handler = self.router
	// Activate middlewares in FIFO order.
	if len(self.middlewares) > 0 {
		for index := len(self.middlewares) - 1; index >= 0; index-- {
			mux = self.middlewares[index](mux)
		}
	}
	mux.ServeHTTP(writer, request)
}

// Serve starts serving the requests at the pre-defined address from settings.
func (self *Server) Run() {
	address := fmt.Sprintf("%s:%d", settings.Host, settings.Port)
	log.Printf("Application server started [%s]", address)
	if err := http.ListenAndServe(address, self); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
