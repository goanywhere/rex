package rex

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	middleware *middleware
	mux        *mux.Router
	ready      bool
	subservers []*Server
}

func New() *Server {
	return &Server{
		middleware: new(middleware),
		mux:        mux.NewRouter().StrictSlash(true),
	}
}

// build constructs all server/subservers along with their middleware modules chain.
func (self *Server) build() http.Handler {
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
func (self *Server) register(pattern string, handler interface{}, methods ...string) {
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
func (self *Server) Any(pattern string, handler interface{}) {
	self.register(pattern, handler, "GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD")
}

// Get is a shortcut for mux.HandleFunc(pattern, handler).Methods("GET"),
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Server) Get(pattern string, handler interface{}) {
	self.register(pattern, handler, "GET")
}

// Head is a shortcut for mux.HandleFunc(pattern, handler).Methods("HEAD")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Server) Head(pattern string, handler interface{}) {
	self.register(pattern, handler, "HEAD")
}

// Options is a shortcut for mux.HandleFunc(pattern, handler).Methods("OPTIONS")
// it also fetch the full function name of the handler (with package) to name the route.
// NOTE method OPTIONS is **NOT** cachable, beware of what you are going to do.
func (self *Server) Options(pattern string, handler interface{}) {
	self.register(pattern, handler, "OPTIONS")
}

// Post is a shortcut for mux.HandleFunc(pattern, handler).Methods("POST")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Server) Post(pattern string, handler interface{}) {
	self.register(pattern, handler, "POST")
}

// Put is a shortcut for mux.HandleFunc(pattern, handler).Methods("PUT")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Server) Put(pattern string, handler interface{}) {
	self.register(pattern, handler, "PUT")
}

// Delete is a shortcut for mux.HandleFunc(pattern, handler).Methods("DELETE")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Server) Delete(pattern string, handler interface{}) {
	self.register(pattern, handler, "DELETE")
}

// Group creates a new application group under the given path prefix.
func (self *Server) Group(prefix string) *Server {
	var middleware = new(middleware)
	self.mux.PathPrefix(prefix).Handler(middleware)
	var mux = self.mux.PathPrefix(prefix).Subrouter()

	server := &Server{middleware: middleware, mux: mux}
	self.subservers = append(self.subservers, server)
	return server
}

// Name returns route name for the given request, if any.
func (self *Server) Name(r *http.Request) (name string) {
	var match mux.RouteMatch
	if self.mux.Match(r, &match) {
		name = match.Route.GetName()
	}
	return name
}

// Use add the middleware module into the stack chain.
func (self *Server) Use(module func(http.Handler) http.Handler) {
	self.middleware.stack = append(self.middleware.stack, module)
}

// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (self *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	self.build().ServeHTTP(w, r)
}

// Run starts the application server to serve incoming requests at the given address.
func (self *Server) Run() {
	configure()
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
func (self *Server) Vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}
