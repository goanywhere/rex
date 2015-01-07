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
package livereload

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

/* ----------------------------------------------------------------------
 * WebSocket Server
 * ----------------------------------------------------------------------*/
const (
	WebSocket  string = "/livereload"
	JavaScript string = "/livereload.js"
)

var (
	once sync.Once

	broadcast chan []byte
	tunnels   map[*tunnel]bool

	in  chan *tunnel
	out chan *tunnel

	mutex sync.RWMutex

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// Alert sends a notice message to browser's livereload.js.
func Alert(message string) {
	go func() {
		var bytes, _ = json.Marshal(&alert{
			Command: "alert",
			Message: message,
		})
		broadcast <- bytes
	}()
}

// Reload sends a reload message to browser's livereload.js.
func Reload() {
	go func() {
		var bytes, _ = json.Marshal(&reload{
			Command: "reload",
			Path:    WebSocket,
			LiveCSS: true,
		})
		broadcast <- bytes
	}()
}

// run watches/dispatches all tunnel & tunnel messages.
func run() {
	for {
		select {
		case tunnel := <-in:
			mutex.Lock()
			defer mutex.Unlock()
			tunnels[tunnel] = true

		case tunnel := <-out:
			mutex.Lock()
			defer mutex.Unlock()
			delete(tunnels, tunnel)
			close(tunnel.message)

		case m := <-broadcast:
			for tunnel := range tunnels {
				select {
				case tunnel.message <- m:
				default:
					mutex.Lock()
					defer mutex.Unlock()
					delete(tunnels, tunnel)
					close(tunnel.message)
				}
			}
		}
	}
}

// Serve serves as a livereload server for accepting I/O tunnel messages.
func Serve(w http.ResponseWriter, r *http.Request) {
	var socket, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	tunnel := new(tunnel)
	tunnel.socket = socket
	tunnel.message = make(chan []byte, 256)

	in <- tunnel
	defer func() { out <- tunnel }()

	tunnel.connect()
}

// ServeJS serves livereload.js for browser.
func ServeJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	w.Write(javascript)
}

// Start activates livereload server for accepting tunnel messages.
func Start() {
	once.Do(func() {
		broadcast = make(chan []byte)
		tunnels = make(map[*tunnel]bool)

		in = make(chan *tunnel)
		out = make(chan *tunnel)

		go run()
	})
}
