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
package rex

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/gorilla/mux"
)

type middleware struct {
	cache http.Handler
	stack []func(http.Handler) http.Handler
}

// build sets up the whole middleware modules in a FIFO chain.
func (self *middleware) build() http.Handler {
	if self.cache == nil {
		var next http.Handler = http.DefaultServeMux
		// Activate modules in FIFO order.
		for index := len(self.stack) - 1; index >= 0; index-- {
			next = self.stack[index](next)
		}
		self.cache = next
	}
	return self.cache
}

// Use add the middleware module into the stack chain.
func (self *middleware) Use(modules ...func(http.Handler) http.Handler) {
	self.stack = append(self.stack, modules...)
}

// Implements the net/http Handler interface and calls the middleware stack.
func (self *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	self.build().ServeHTTP(w, r)
}

type Router struct {
	middleware *middleware
	mux        *mux.Router
	ready      bool
	subrouters []*Router
}

func New() *Router {
	return &Router{
		middleware: new(middleware),
		mux:        mux.NewRouter().StrictSlash(true),
	}
}

// build constructs all router/subrouters along with their middleware modules chain.
func (self *Router) build() http.Handler {
	if !self.ready {
		self.ready = true
		// * activate router's middleware modules.
		self.middleware.Use(func(http.Handler) http.Handler {
			return self.mux
		})
		// * activate subrouters's middleware modules.
		for index := 0; index < len(self.subrouters); index++ {
			sr := self.subrouters[index]
			sr.middleware.Use(func(http.Handler) http.Handler {
				return sr.mux
			})
		}
	}
	return self.middleware
}

// register adds the http.Handler/http.HandleFunc into Gorilla mux.
func (self *Router) register(method string, pattern string, handler interface{}) {
	// finds the full function name (with package) as its mappings.
	var name = runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()

	switch H := handler.(type) {
	case http.Handler:
		self.mux.Handle(pattern, H).Methods(method).Name(name)

	case func(http.ResponseWriter, *http.Request):
		self.mux.HandleFunc(pattern, H).Methods(method).Name(name)

	default:
		Fatalf("Unsupported handler (%s) passed in.", name)
	}
}

// Get is a shortcut for mux.HandleFunc(pattern, handler).Methods("GET"),
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Router) Get(pattern string, handler interface{}) {
	self.register("GET", pattern, handler)
}

// Head is a shortcut for mux.HandleFunc(pattern, handler).Methods("HEAD")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Router) Head(pattern string, handler interface{}) {
	self.register("HEAD", pattern, handler)
}

// Options is a shortcut for mux.HandleFunc(pattern, handler).Methods("OPTIONS")
// it also fetch the full function name of the handler (with package) to name the route.
// NOTE method OPTIONS is **NOT** cachable, beware of what you are going to do.
func (self *Router) Options(pattern string, handler interface{}) {
	self.register("OPTIONS", pattern, handler)
}

// Post is a shortcut for mux.HandleFunc(pattern, handler).Methods("POST")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Router) Post(pattern string, handler interface{}) {
	self.register("POST", pattern, handler)
}

// Put is a shortcut for mux.HandleFunc(pattern, handler).Methods("PUT")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Router) Put(pattern string, handler interface{}) {
	self.register("PUT", pattern, handler)
}

// Delete is a shortcut for mux.HandleFunc(pattern, handler).Methods("DELETE")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Router) Delete(pattern string, handler interface{}) {
	self.register("Delete", pattern, handler)
}

// Group creates a new application group under the given path prefix.
func (self *Router) Group(prefix string) *Router {
	var middleware = new(middleware)
	self.mux.PathPrefix(prefix).Handler(middleware)
	var mux = self.mux.PathPrefix(prefix).Subrouter()

	router := &Router{middleware: middleware, mux: mux}
	self.subrouters = append(self.subrouters, router)
	return router
}

// Name returns route name for the given request, if any.
func (self *Router) Name(r *http.Request) (name string) {
	var match mux.RouteMatch
	if self.mux.Match(r, &match) {
		name = match.Route.GetName()
	}
	return name
}

// Use add the middleware module into the stack chain.
func (self *Router) Use(module func(http.Handler) http.Handler) {
	self.middleware.Use(module)
}

// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (self *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	self.build().ServeHTTP(w, r)
}

// Run starts the application server to serve incoming requests at the given address.
func (self *Router) Run() {
	runtime.GOMAXPROCS(config.maxprocs)

	go func() {
		time.Sleep(500 * time.Millisecond)
		Infof("Application server is listening at %d", config.port)
	}()

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.port), self); err != nil {
		Fatalf("Failed to start the server: %v", err)
	}
}
