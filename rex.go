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
	"flag"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/goanywhere/rex/internal"
	"github.com/goanywhere/rex/modules"
	"github.com/goanywhere/x/env"
	"github.com/goanywhere/x/fs"
)

var (
	mux = New()

	options = struct {
		debug    bool
		port     int
		maxprocs int
	}{
		debug:    true,
		port:     5000,
		maxprocs: runtime.NumCPU(),
	}
)

// Get is a shortcut for mux.HandleFunc(pattern, handler).Methods("GET"),
// it also fetch the full function name of the handler (with package) to name the route.
func Get(pattern string, handler interface{}) {
	mux.register("GET", pattern, handler)
}

// Head is a shortcut for mux.HandleFunc(pattern, handler).Methods("HEAD")
// it also fetch the full function name of the handler (with package) to name the route.
func Head(pattern string, handler interface{}) {
	mux.register("HEAD", pattern, handler)
}

// Options is a shortcut for mux.HandleFunc(pattern, handler).Methods("OPTIONS")
// it also fetch the full function name of the handler (with package) to name the route.
// NOTE method OPTIONS is **NOT** cachable, beware of what you are going to do.
func Options(pattern string, handler interface{}) {
	mux.register("OPTIONS", pattern, handler)
}

// Post is a shortcut for mux.HandleFunc(pattern, handler).Methods("POST")
// it also fetch the full function name of the handler (with package) to name the route.
func Post(pattern string, handler interface{}) {
	mux.register("POST", pattern, handler)
}

// Put is a shortcut for mux.HandleFunc(pattern, handler).Methods("PUT")
// it also fetch the full function name of the handler (with package) to name the route.
func Put(pattern string, handler interface{}) {
	mux.register("PUT", pattern, handler)
}

// Delete is a shortcut for mux.HandleFunc(pattern, handler).Methods("DELETE")
// it also fetch the full function name of the handler (with package) to name the route.
func Delete(pattern string, handler interface{}) {
	mux.register("Delete", pattern, handler)
}

// Group creates a new application group under the given path.
func Group(path string) *Router {
	return mux.Group(path)
}

// FileServer registers a handler to serve HTTP requests
// with the contents of the file system rooted at root.
func FileServer(prefix, dir string) {
	mux.FileServer(prefix, dir)
}

// Use appends middleware module into the serving list, modules will be served in FIFO order.
func Use(module func(http.Handler) http.Handler) {
	mux.Use(module)
}

func Run() {
	mux.Use(modules.Logger)
	mux.Run()
}

func init() {
	// ----------------------------------------
	// Project Root
	// ----------------------------------------
	var root = fs.Getcd(2)
	env.Set(internal.Root, root)
	env.Load(filepath.Join(root, ".env"))

	// cmd arguments
	flag.BoolVar(&options.debug, "debug", options.debug, "flag to toggle debug mode")
	flag.IntVar(&options.port, "port", options.port, "port to run the application server")
	flag.IntVar(&options.maxprocs, "maxprocs", options.maxprocs, "maximum cpu processes to run the server")

	flag.Parse()
}
