package rex

import (
	"flag"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/goanywhere/fs"
	"github.com/goanywhere/rex/internal"
	mw "github.com/goanywhere/rex/middleware"
	"github.com/goanywhere/x/env"
)

var (
	DefaultMux = New()

	config = &options{
		debug:    true,
		port:     5000,
		maxprocs: runtime.NumCPU(),
	}
)

// Get is a shortcut for mux.HandleFunc(pattern, handler).Methods("GET"),
// it also fetch the full function name of the handler (with package) to name the route.
func Get(pattern string, handler interface{}) {
	DefaultMux.Get(pattern, handler)
}

// Head is a shortcut for mux.HandleFunc(pattern, handler).Methods("HEAD")
// it also fetch the full function name of the handler (with package) to name the route.
func Head(pattern string, handler interface{}) {
	DefaultMux.Head(pattern, handler)
}

// Options is a shortcut for mux.HandleFunc(pattern, handler).Methods("OPTIONS")
// it also fetch the full function name of the handler (with package) to name the route.
// NOTE method OPTIONS is **NOT** cachable, beware of what you are going to do.
func Options(pattern string, handler interface{}) {
	DefaultMux.Options(pattern, handler)
}

// Post is a shortcut for mux.HandleFunc(pattern, handler).Methods("POST")
// it also fetch the full function name of the handler (with package) to name the route.
func Post(pattern string, handler interface{}) {
	DefaultMux.Post(pattern, handler)
}

// Put is a shortcut for mux.HandleFunc(pattern, handler).Methods("PUT")
// it also fetch the full function name of the handler (with package) to name the route.
func Put(pattern string, handler interface{}) {
	DefaultMux.Put(pattern, handler)
}

// Delete is a shortcut for mux.HandleFunc(pattern, handler).Methods("DELETE")
// it also fetch the full function name of the handler (with package) to name the route.
func Delete(pattern string, handler interface{}) {
	DefaultMux.Delete(pattern, handler)
}

// Group creates a new application group under the given path.
func Group(path string) *Router {
	return DefaultMux.Group(path)
}

// Use appends middleware module into the serving list, modules will be served in FIFO order.
func Use(module func(http.Handler) http.Handler) {
	DefaultMux.Use(module)
}

func Run() {
	DefaultMux.Use(mw.Logger)
	DefaultMux.Run()
}

func init() {
	// setup project root
	var root = fs.Getcd(2)
	env.Set(internal.ROOT, root)
	env.Load(filepath.Join(root, ".env"))

	// cmd arguments
	flag.BoolVar(&config.debug, "debug", config.debug, "flag to toggle debug mode")
	flag.IntVar(&config.port, "port", config.port, "port to run the application server")
	flag.IntVar(&config.maxprocs, "maxprocs", config.maxprocs, "maximum cpu processes to run the server")

	flag.Parse()
}
