/* ----------------------------------------------------------------------
 *       ______      ___                         __
 *      / ____/___  /   |  ____  __  ___      __/ /_  ___  ________
 *     / / __/ __ \/ /| | / __ \/ / / / | /| / / __ \/ _ \/ ___/ _ \
 *    / /_/ / /_/ / ___ |/ / / / /_/ /| |/ |/ / / / /  __/ /  /  __/
 *    \____/\____/_/  |_/_/ /_/\__. / |__/|__/_/ /_/\___/_/   \___/
 *                            /____/
 *
 * (C) Copyright 2015 GoAnywhere (http://goanywhere.io).
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
	"log"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/gorilla/mux"
)

type (
	Mux struct {
		router  *mux.Router
		modules []Module
	}

	HandlerFunc func(*Context)

	// Conventional method to implement custom modules.
	Module func(http.Handler) http.Handler
)

// Custom handler func provides Context Supports.
func (self HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	self(NewContext(w, r))
}

// New creates a plain web server without any middleware modules.
func New() *Mux {
	self := new(Mux)
	self.router = mux.NewRouter()
	return self
}

// ---------------------------------------------------------------------------
//  HTTP Requests Handlers
// ---------------------------------------------------------------------------
// Supported Handler Types
//	* http.Handler
//	* http.HandlerFunc	=> func(w http.ResponseWriter, r *http.Request)
//	* rex.HandlerFunc	=> func(ctx *Context)
func (self *Mux) register(method, pattern string, h interface{}) {
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
		log.Fatalf("Unknown handler type (%v) passed in.", h)
	}
	// finds the full function name (with package)
	name := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	self.router.Handle(pattern, handler).Methods(method).Name(name)
}

// Get is a shortcut for mux.HandleFunc(pattern, handler).Methods("GET"),
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Mux) Get(pattern string, handler interface{}) {
	self.register("GET", pattern, handler)
}

// Post is a shortcut for mux.HandleFunc(pattern, handler).Methods("POST")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Mux) Post(pattern string, handler interface{}) {
	self.register("POST", pattern, handler)
}

// Put is a shortcut for mux.HandleFunc(pattern, handler).Methods("PUT")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Mux) Put(pattern string, handler interface{}) {
	self.register("PUT", pattern, handler)
}

// Delete is a shortcut for mux.HandleFunc(pattern, handler).Methods("DELETE")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Mux) Delete(pattern string, handler interface{}) {
	self.register("DELETE", pattern, handler)
}

// Patch is a shortcut for mux.HandleFunc(pattern, handler).Methods("PATCH")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Mux) Patch(pattern string, handler http.HandlerFunc) {
	self.register("PATCH", pattern, handler)
}

// Head is a shortcut for mux.HandleFunc(pattern, handler).Methods("HEAD")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Mux) Head(pattern string, handler http.HandlerFunc) {
	self.register("HEAD", pattern, handler)
}

// Options is a shortcut for mux.HandleFunc(pattern, handler).Methods("OPTIONS")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Mux) Options(pattern string, handler http.HandlerFunc) {
	self.register("OPTIONS", pattern, handler)
}

// Group creates a new application group under the given path.
func (self *Mux) Group(path string) *Mux {
	mux := new(Mux)
	mux.router = self.router.PathPrefix(path).Subrouter()
	return mux
}

// Add appends middleware module into the serving list, modules will be served in FIFO order.
func (self *Mux) Use(modules ...interface{}) {
	var mod Module
	for _, module := range modules {
		switch module.(type) {
		// Standard http.Handler module.
		case func(http.Handler) http.Handler:
			mod = module.(func(http.Handler) http.Handler)

		case func(map[string]interface{}) func(http.Handler) http.Handler:
			// http.Handler with module options (using default Options).
			var options = make(map[string]interface{})
			mod = module.(func(map[string]interface{}) func(http.Handler) http.Handler)(options)

		default:
			log.Fatalf("Unknown module type (%v) passed in.", module)
		}
		//name := runtime.FuncForPC(reflect.ValueOf(module).Pointer()).Name()
		self.modules = append(self.modules, mod)
	}
}

// ServeHTTP: Implementation of "http.Handler" interface.
func (self *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var mux http.Handler = self.router
	// Activate modules in FIFO order.
	if len(self.modules) > 0 {
		for index := len(self.modules) - 1; index >= 0; index-- {
			mux = self.modules[index](mux)
		}
	}
	mux.ServeHTTP(w, r)
}

// Run starts the application server to serve incoming requests at the given address.
func (self *Mux) Run(address string) {
	go func() {
		time.Sleep(100 * time.Millisecond)
		log.Printf("Application server started [%s]", address)
	}()
	if err := http.ListenAndServe(address, self); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
