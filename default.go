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

import "flag"

// ---------------------------------------------------------------------------
//  Default Server Mux
// ---------------------------------------------------------------------------
var server *Server

func Run() {
	// common server middleware modules.
	//server.Use(modules.XSRF)
	//server.Use(modules.Env)

	flag.Parse()
	server.Run()
}

func Get(pattern string, handler interface{}) {
	server.Get(pattern, handler)
}

func Post(pattern string, handler interface{}) {
	server.Post(pattern, handler)
}

func Put(pattern string, handler interface{}) {
	server.Put(pattern, handler)
}

func Delete(pattern string, handler interface{}) {
	server.Delete(pattern, handler)
}

func Patch(pattern string, handler interface{}) {
	server.Patch(pattern, handler)
}

func Head(pattern string, handler interface{}) {
	server.Head(pattern, handler)
}

func (self *Server) Options(pattern string, handler interface{}) {
	server.Options(pattern, handler)
}

// Group creates a new application group under the given path.
func Group(path string) *Server {
	return server.Group(path)
}

// FileServer registers a handler to serve HTTP requests
// with the contents of the file system rooted at root.
func FileServer(prefix, dir string) {
	server.FileServer(prefix, dir)
}

// Use appends middleware module into the serving list, modules will be served in FIFO order.
func Use(modules ...Module) {
	server.Use(modules...)
}
