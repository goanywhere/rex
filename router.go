package rex

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/gorilla/mux"
)

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
		// * add router into middlware stack to serve as final http.Handler.
		self.Use(func(http.Handler) http.Handler {
			return self.mux
		})
		// * add subrouters into middlware stack to serve as final http.Handler.
		for index := 0; index < len(self.subrouters); index++ {
			router := self.subrouters[index]
			router.Use(func(http.Handler) http.Handler {
				return router.mux
			})
		}
		self.ready = true
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
	self.middleware.stack = append(self.middleware.stack, module)
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
