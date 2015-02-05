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
	"github.com/gorilla/securecookie"

	"github.com/goanywhere/rex/internal"
	"github.com/goanywhere/rex/template"
	"github.com/goanywhere/x/fs"
)

var app struct {
	HTML   *template.Loader
	codecs []securecookie.Codec
}

type (
	Server struct {
		mux     *mux.Router
		pool    sync.Pool
		modules []Module
	}

	// Conventional method to implement custom modules.
	Module func(http.Handler) http.Handler

	Handler interface {
		ServeHTTP(*Context)
	}

	HandlerFunc func(*Context)
)

// Serve wraps standard ServeHTTP function with context.
func (self HandlerFunc) ServeHTTP(ctx *Context) {
	self(ctx)
}

// New creates a plain web server without any middleware modules.
func New() *Server {
	self := new(Server)
	self.mux = mux.NewRouter()
	self.configure()
	self.pool.New = func() interface{} {
		return NewContext(nil, nil)
	}
	return self
}

// configure initialize all application related settings before running.
// Server Secret Keys
func (self *Server) configure() {
	settings := internal.Settings()
	settings.Load(".env")
	// ------------------------------------------------
	// if secret keys exists, create codecs.
	// ------------------------------------------------
	if keys := settings.Strings("secret.keys"); len(keys) > 0 {
		var bytes [][]byte
		for _, key := range keys {
			bytes = append(bytes, []byte(key))
		}
		app.codecs = securecookie.CodecsFromPairs(bytes...)
	} else {
		log.Fatalf("Failed to setup application: secret key(s) missing")
	}
	// ------------------------------------------------
	// templates folder exists => load HTML templates.
	// ------------------------------------------------
	if dir := settings.String("dir.templates", "templates"); fs.Exists(dir) {
		app.HTML = template.NewLoader(dir)
		app.HTML.Load()
	}
}

// creates a reusable context for consequent requests.
func (self *Server) createContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := self.pool.Get().(*Context)
	ctx.Writer = w
	ctx.Request = r
	return ctx
}

// put the created context object into pool for consequent creation.
func (self *Server) recycleContext(ctx *Context) {
	defer self.pool.Put(ctx)
}

// ---------------------------------------------------------------------------
//  HTTP Requests Handlers
// ---------------------------------------------------------------------------
// Supported Handler Types
//	* http.Handler
//	* http.HandlerFunc	=> func(w http.ResponseWriter, r *http.Request)
//	* rex.HandlerFunc	=> func(ctx *Context)
func (self *Server) register(method, pattern string, handler interface{}) {
	// finds the full function name (with package) as its mappings.
	var name = runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()

	switch H := handler.(type) {
	case http.Handler:
		self.mux.Handle(pattern, H).Methods(method).Name(name)

	case func(http.ResponseWriter, *http.Request):
		self.mux.HandleFunc(pattern, H).Methods(method).Name(name)

	case Handler:
		self.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			ctx := self.createContext(w, r)
			defer self.recycleContext(ctx)
			H.ServeHTTP(ctx)
		}).Methods(method).Name(name)

	case func(*Context):
		self.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			ctx := self.createContext(w, r)
			defer self.recycleContext(ctx)
			H(ctx)
		}).Methods(method).Name(name)

	default:
		log.Fatalf("Unknown handler type (%v) passed in.", H)
	}
}

// Get is a shortcut for mux.HandleFunc(pattern, handler).Methods("GET"),
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Server) Get(pattern string, handler interface{}) {
	self.register("GET", pattern, handler)
}

// Post is a shortcut for mux.HandleFunc(pattern, handler).Methods("POST")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Server) Post(pattern string, handler interface{}) {
	self.register("POST", pattern, handler)
}

// Put is a shortcut for mux.HandleFunc(pattern, handler).Methods("PUT")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Server) Put(pattern string, handler interface{}) {
	self.register("PUT", pattern, handler)
}

// Delete is a shortcut for mux.HandleFunc(pattern, handler).Methods("DELETE")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Server) Delete(pattern string, handler interface{}) {
	self.register("DELETE", pattern, handler)
}

// Patch is a shortcut for mux.HandleFunc(pattern, handler).Methods("PATCH")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Server) Patch(pattern string, handler interface{}) {
	self.register("PATCH", pattern, handler)
}

// Head is a shortcut for mux.HandleFunc(pattern, handler).Methods("HEAD")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Server) Head(pattern string, handler interface{}) {
	self.register("HEAD", pattern, handler)
}

// Options is a shortcut for mux.HandleFunc(pattern, handler).Methods("OPTIONS")
// it also fetch the full function name of the handler (with package) to name the route.
// NOTE method OPTIONS is **NOT** cachable, beware of what you are going to do.
func (self *Server) Options(pattern string, handler interface{}) {
	self.register("OPTIONS", pattern, handler)
}

// Group creates a new application group under the given path.
func (self *Server) Group(path string) *Server {
	return &Server{mux: self.mux.PathPrefix(path).Subrouter()}
}

// FileServer registers a handler to serve HTTP requests
// with the contents of the file system rooted at root.
func (self *Server) FileServer(prefix, dir string) {
	settings := internal.Settings()
	if abs, err := filepath.Abs(dir); err == nil {
		settings.Set("url.static", prefix)
		server := http.StripPrefix(prefix, http.FileServer(http.Dir(abs)))
		self.mux.PathPrefix(prefix).Handler(server)
	} else {
		log.Fatalf("Failed to setup file server: %v", err)
	}
}

// Use appends middleware module into the serving list, modules will be served in FIFO order.
func (self *Server) Use(modules ...Module) {
	self.modules = append(self.modules, modules...)
}

// ServeHTTP: Implementation of "http.Handler" interface.
func (self *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
func (self *Server) Run(address string) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	go func() {
		time.Sleep(500 * time.Millisecond)
		log.Printf("Application server started [%s]", address)
	}()
	if err := http.ListenAndServe(address, self); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
