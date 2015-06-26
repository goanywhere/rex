package rex

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/goanywhere/env"
	"github.com/gorilla/mux"
)

var (
	debug    bool
	port     int
	maxprocs int

	once sync.Once
)

type server struct {
	middleware *middleware
	mux        *mux.Router
	ready      bool
	subservers []*server
}

func New() *server {
	self := &server{
		middleware: new(middleware),
		mux:        mux.NewRouter().StrictSlash(true),
	}
	self.configure()
	return self
}

func (self *server) configure() {
	once.Do(func() {
		flag.BoolVar(&debug, "debug", env.Bool("DEBUG", true), "flag to toggle debug mode")
		flag.IntVar(&port, "port", env.Int("PORT", 5000), "port to run the application server")
		flag.IntVar(&maxprocs, "maxprocs", env.Int("MAXPROCS", runtime.NumCPU()), "maximum cpu processes to run the server")
		flag.Parse()
	})
}

// build constructs all server/subservers along with their middleware modules chain.
func (self *server) build() http.Handler {
	if !self.ready {
		// * add server mux into middlware stack to serve as final http.Handler.
		self.Use(func(http.Handler) http.Handler {
			return self.mux
		})
		// * add subservers into middlware stack to serve as final http.Handler.
		for index := 0; index < len(self.subservers); index++ {
			server := self.subservers[index]
			server.Use(func(http.Handler) http.Handler {
				return server.mux
			})
		}
		self.ready = true
	}
	return self.middleware
}

// register adds the http.Handler/http.HandleFunc into Gorilla mux.
func (self *server) register(pattern string, handler interface{}, methods ...string) {
	// finds the full function name (with package) as its mappings.
	var name = runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()

	switch H := handler.(type) {
	case http.Handler:
		self.mux.Handle(pattern, H).Methods(methods...).Name(name)

	case func(http.ResponseWriter, *http.Request):
		self.mux.HandleFunc(pattern, H).Methods(methods...).Name(name)

	default:
		Fatalf("Unsupported handler (%s) passed in.", name)
	}
}

// Any maps most common HTTP methods request to the given `http.Handler`.
// Supports: GET | POST | PUT | DELETE | OPTIONS | HEAD
func (self *server) Any(pattern string, handler interface{}) {
	self.register(pattern, handler, "GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD")
}

// Group creates a new application group under the given path prefix.
func (self *server) Group(prefix string) *server {
	var middleware = new(middleware)
	self.mux.PathPrefix(prefix).Handler(middleware)
	var mux = self.mux.PathPrefix(prefix).Subrouter()

	server := &server{middleware: middleware, mux: mux}
	self.subservers = append(self.subservers, server)
	return server
}

// Name returns route name for the given request, if any.
func (self *server) Name(r *http.Request) (name string) {
	var match mux.RouteMatch
	if self.mux.Match(r, &match) {
		name = match.Route.GetName()
	}
	return name
}

// FileServer registers a handler to serve HTTP (GET|HEAD) requests
// with the contents of file system under the given directory.
func (self *server) FileServer(prefix, dir string) {
	if abs, err := filepath.Abs(dir); err == nil {
		fs := http.StripPrefix(prefix, http.FileServer(http.Dir(abs)))
		self.mux.PathPrefix(prefix).Handler(fs)
	} else {
		log.Fatalf("Failed to setup file server: %v", err)
	}
}

// Use add the middleware module into the stack chain.
func (self *server) Use(module func(http.Handler) http.Handler) {
	self.middleware.stack = append(self.middleware.stack, module)
}

// GET is a shortcut for mux.HandleFunc(pattern, handler).Methods("GET"),
// it also fetch the full function name of the handler (with package) to name the route.
func (self *server) GET(pattern string, handler interface{}) {
	self.register(pattern, handler, "GET")
}

// HEAD is a shortcut for mux.HandleFunc(pattern, handler).Methods("HEAD")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *server) HEAD(pattern string, handler interface{}) {
	self.register(pattern, handler, "HEAD")
}

// OPTIONS is a shortcut for mux.HandleFunc(pattern, handler).Methods("OPTIONS")
// it also fetch the full function name of the handler (with package) to name the route.
// NOTE method OPTIONS is **NOT** cachable, beware of what you are going to do.
func (self *server) OPTIONS(pattern string, handler interface{}) {
	self.register(pattern, handler, "OPTIONS")
}

// POST is a shortcut for mux.HandleFunc(pattern, handler).Methods("POST")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *server) POST(pattern string, handler interface{}) {
	self.register(pattern, handler, "POST")
}

// PUT is a shortcut for mux.HandleFunc(pattern, handler).Methods("PUT")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *server) PUT(pattern string, handler interface{}) {
	self.register(pattern, handler, "PUT")
}

// DELETE is a shortcut for mux.HandleFunc(pattern, handler).Methods("DELETE")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *server) DELETE(pattern string, handler interface{}) {
	self.register(pattern, handler, "DELETE")
}

// TRACE is a shortcut for mux.HandleFunc(pattern, handler).Methods("TRACE")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *server) TRACE(pattern string, handler interface{}) {
	self.register(pattern, handler, "TRACE")
}

// CONNECT is a shortcut for mux.HandleFunc(pattern, handler).Methods("CONNECT")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *server) CONNECT(pattern string, handler interface{}) {
	self.register(pattern, handler, "CONNECT")
}

// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (self *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	self.build().ServeHTTP(w, r)
}

// Run starts the application server to serve incoming requests at the given address.
func (self *server) Run() {
	runtime.GOMAXPROCS(maxprocs)

	go func() {
		time.Sleep(500 * time.Millisecond)
		Infof("Application server is listening at %d", port)
	}()

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), self); err != nil {
		Fatalf("Failed to start the server: %v", err)
	}
}

// Vars returns the route variables for the current request, if any.
func (self *server) Vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}
