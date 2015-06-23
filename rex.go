package rex

import (
	"net/http"
	"path/filepath"

	"github.com/goanywhere/fs"
	"github.com/goanywhere/rex/internal"
	. "github.com/goanywhere/rex/middleware"
)

var (
	Default = New()
	Env     *internal.Env
)

// Get is a shortcut for mux.HandleFunc(pattern, handler).Methods("GET"),
// it also fetch the full function name of the handler (with package) to name the route.
func Get(pattern string, handler interface{}) {
	Default.Get(pattern, handler)
}

// Head is a shortcut for mux.HandleFunc(pattern, handler).Methods("HEAD")
// it also fetch the full function name of the handler (with package) to name the route.
func Head(pattern string, handler interface{}) {
	Default.Head(pattern, handler)
}

// Options is a shortcut for mux.HandleFunc(pattern, handler).Methods("OPTIONS")
// it also fetch the full function name of the handler (with package) to name the route.
// NOTE method OPTIONS is **NOT** cachable, beware of what you are going to do.
func Options(pattern string, handler interface{}) {
	Default.Options(pattern, handler)
}

// Post is a shortcut for mux.HandleFunc(pattern, handler).Methods("POST")
// it also fetch the full function name of the handler (with package) to name the route.
func Post(pattern string, handler interface{}) {
	Default.Post(pattern, handler)
}

// Put is a shortcut for mux.HandleFunc(pattern, handler).Methods("PUT")
// it also fetch the full function name of the handler (with package) to name the route.
func Put(pattern string, handler interface{}) {
	Default.Put(pattern, handler)
}

// Delete is a shortcut for mux.HandleFunc(pattern, handler).Methods("DELETE")
// it also fetch the full function name of the handler (with package) to name the route.
func Delete(pattern string, handler interface{}) {
	Default.Delete(pattern, handler)
}

// Group creates a new application group under the given path.
func Group(path string) *Server {
	return Default.Group(path)
}

// FileServer registers a handler to serve HTTP (GET|HEAD) requests
// with the contents of file system under the given directory.
func FileServer(prefix, dir string) {
	Default.FileServer(prefix, dir)
}

// Use appends middleware module into the serving list, modules will be served in FIFO order.
func Use(module func(http.Handler) http.Handler) {
	Default.Use(module)
}

func Run() {
	Default.Use(Logger)
	Default.Run()
}

func init() {
	var basedir = fs.Getcd(2)
	Env = internal.New(filepath.Join(basedir, ".env"))
}
