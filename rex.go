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
/*
Package rex provides an out-of-box web server with common middleware modules.
Example:
	package main

	import (
		"net/http"
		"github.com/goanywhere/rex"
		"github.com/goanywhere/rex/web"
	)

	func index(ctx *web.Context) {
		ctx.String("Hello World")
	}

	func main() {
		rex.Get("/", index)
		rex.Run()
	}
*/
package rex

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/goanywhere/rex/modules"
	"github.com/goanywhere/rex/web"
	"github.com/goanywhere/x/env"
)

// default rex server with reasonable middleware modules.
var (
	port int
	mux  = web.New()
)

type H map[string]interface{}
type Context web.Context

func Get(pattern string, handler interface{}) {
	mux.Get(pattern, handler)
}

func Post(pattern string, handler interface{}) {
	mux.Post(pattern, handler)
}

func Put(pattern string, handler interface{}) {
	mux.Put(pattern, handler)
}

func Delete(pattern string, handler interface{}) {
	mux.Delete(pattern, handler)
}

func Patch(pattern string, handler http.HandlerFunc) {
	mux.Patch(pattern, handler)
}

func Head(pattern string, handler http.HandlerFunc) {
	mux.Head(pattern, handler)
}

func Options(pattern string, handler http.HandlerFunc) {
	mux.Options(pattern, handler)
}

func Group(path string) *web.Mux {
	return mux.Group(path)
}

func Use(modules ...interface{}) {
	mux.Use(modules...)
}

// Serve starts serving the requests at the pre-defined address from settings.
func Run() {
	flag.Parse()
	mux.Run(fmt.Sprintf(":%d", port))
}

func init() {
	mux.Use(modules.Env)
	mux.Use(modules.XSRF)

	if cwd, err := os.Getwd(); err != nil {
		log.Fatalf("Failed to retrieve project root: %v", err)
	} else {
		root, _ := filepath.Abs(cwd)
		env.Set("root", root)
	}
	env.Set("port", string(port))
	env.Set("mode", "debug")
	env.Set("dir.static", "build")
	env.Set("dir.templates", "templates")
	env.Set("url.static", "/static/")
	// default environmental headers for modules.Env
	env.Set("header.x.ua.compatible", "deny")
	env.Set("header.x.frame.options", "nosniff")
	env.Set("header.x.xss.protection", "1; mode=block")
	env.Set("header.x.content.type.options", "IE=Edge,chrome=1")
	// custom settings
	env.Load(".env")
	// cmd parameters take the priority
	flag.IntVar(&port, "port", 5000, "port to run the application server")
}
