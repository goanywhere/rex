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
	"log"
	"net/http"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"github.com/goanywhere/rex/internal"
	"github.com/goanywhere/rex/template"
)

type (
	App struct {
		mux     *mux.Router
		pool    sync.Pool
		modules []Module

		store  *sessions.CookieStore
		loader *template.Loader
	}

	// Conventional method to implement custom modules.
	Module func(http.Handler) http.Handler

	HandlerFunc func(*Context)
)

// New creates a plain web server without any middleware modules.
func New() *App {
	self := new(App)
	self.mux = mux.NewRouter()
	self.configure()
	self.pool.New = func() interface{} {
		return &Context{app: self}
	}
	return self
}

// configure initialize all application related settings before running.
// App Secret Keys
func (self *App) configure() {
	options := internal.Options()
	options.Load(".env")
	// Fetch @ least a secret key to setup secret cookie.
	var keys [][]byte
	for _, key := range options.Strings("secret.keys") {
		keys = append(keys, []byte(key))
	}
	if len(keys) == 0 || keys[0] == nil {
		log.Fatalf("Failed to setup application: secret key(s) missing")
	}
	self.store = sessions.NewCookieStore(keys...)
	self.loader = template.NewLoader(options.String("dir.templates"))
	self.loader.Load()
}

// context creates/fetches a rex.Context instance for App server.
//	NOTE mux.pool.Put(ctx) must be called to put back the created context.
func (self *App) context(w http.ResponseWriter, r *http.Request) *Context {
	var err error
	ctx := self.pool.Get().(*Context)
	ctx.session, err = self.store.Get(r, Options.String("session.cookie.name"))
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	ctx.Writer = w
	ctx.Request = r
	return ctx
}

// ---------------------------------------------------------------------------
//  HTTP Requests Handlers
// ---------------------------------------------------------------------------
// Supported Handler Types
//	* http.Handler
//	* http.HandlerFunc	=> func(w http.ResponseWriter, r *http.Request)
//	* rex.HandlerFunc	=> func(ctx *Context)
func (self *App) register(method, pattern string, handler interface{}) {
	// finds the full function name (with package) as its mappings.
	var name = runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()

	switch H := handler.(type) {
	case http.Handler:
		self.mux.Handle(pattern, H).Methods(method).Name(name)

	case func(http.ResponseWriter, *http.Request):
		self.mux.HandleFunc(pattern, H).Methods(method).Name(name)

	case func(*Context):
		self.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			ctx := self.context(w, r)
			defer self.pool.Put(ctx)
			H(ctx)
		}).Methods(method).Name(name)

	default:
		log.Fatalf("Unknown handler type (%v) passed in.", H)
	}
}

// Get is a shortcut for mux.HandleFunc(pattern, handler).Methods("GET"),
// it also fetch the full function name of the handler (with package) to name the route.
func (self *App) Get(pattern string, handler interface{}) {
	self.register("GET", pattern, handler)
}

// Post is a shortcut for mux.HandleFunc(pattern, handler).Methods("POST")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *App) Post(pattern string, handler interface{}) {
	self.register("POST", pattern, handler)
}

// Put is a shortcut for mux.HandleFunc(pattern, handler).Methods("PUT")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *App) Put(pattern string, handler interface{}) {
	self.register("PUT", pattern, handler)
}

// Delete is a shortcut for mux.HandleFunc(pattern, handler).Methods("DELETE")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *App) Delete(pattern string, handler interface{}) {
	self.register("DELETE", pattern, handler)
}

// Patch is a shortcut for mux.HandleFunc(pattern, handler).Methods("PATCH")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *App) Patch(pattern string, handler interface{}) {
	self.register("PATCH", pattern, handler)
}

// Head is a shortcut for mux.HandleFunc(pattern, handler).Methods("HEAD")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *App) Head(pattern string, handler interface{}) {
	self.register("HEAD", pattern, handler)
}

// Options is a shortcut for mux.HandleFunc(pattern, handler).Methods("OPTIONS")
// it also fetch the full function name of the handler (with package) to name the route.
// NOTE method OPTIONS is **NOT** cachable, beware of what you are going to do.
func (self *App) Options(pattern string, handler interface{}) {
	self.register("OPTIONS", pattern, handler)
}

// Group creates a new application group under the given path.
func (self *App) Group(path string) *App {
	return &App{mux: self.mux.PathPrefix(path).Subrouter()}
}

// FileServer registers a handler to serve HTTP requests
// with the contents of the file system rooted at root.
func (self *App) FileServer(dir string) {
	if abs, err := filepath.Abs(dir); err == nil {
		prefix := Options.String("url.static")
		self.mux.PathPrefix(prefix).Handler(http.StripPrefix(prefix, http.FileServer(http.Dir(abs))))
	} else {
		log.Fatalf("Failed to setup file server: %v", err)
	}
}

// Use appends middleware module into the serving list, modules will be served in FIFO order.
// TODO simplify ME.
func (self *App) Use(modules ...interface{}) {
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
		self.modules = append(self.modules, mod)
	}
}

// ServeHTTP: Implementation of "http.Handler" interface.
func (self *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var mux http.Handler = self.mux
	// Activate modules in FIFO order.
	if len(self.modules) > 0 {
		for index := len(self.modules) - 1; index >= 0; index-- {
			mux = self.modules[index](mux)
		}
	}
	mux.ServeHTTP(w, r)
}

// Run starts the application server to serve incoming requests at the given address.
func (self *App) Run(address string) {
	go func() {
		time.Sleep(500 * time.Millisecond)
		log.Printf("Application server started [%s]", address)
	}()
	if err := http.ListenAndServe(address, self); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
