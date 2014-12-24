/* ----------------------------------------------------------------------
 *       ______      ___                         __
 *      / ____/___  /   |  ____  __  ___      __/ /_  ___  ________
 *     / / __/ __ \/ /| | / __ \/ / / / | /| / / __ \/ _ \/ ___/ _ \
 *    / /_/ / /_/ / ___ |/ / / / /_/ /| |/ |/ / / / /  __/ /  /  __/
 *    \____/\____/_/  |_/_/ /_/\__. / |__/|__/_/ /_/\___/_/   \___/
 *                            /____/
 *
 * (C) Copyright 2014 GoAnywhere (http://goanywhere.io).
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
	"regexp"
	"sync"

	"github.com/gorilla/websocket"
)

/* ----------------------------------------------------------------------
 * WebSocket Server
 * ----------------------------------------------------------------------*/
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
			Path:    "/livereload",
			LiveCSS: true,
		})
		broadcast <- bytes
	}()
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
	tunnel.handshake = regexp.MustCompile(`"command"\s*:\s*"hello"`)

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
	broadcast = make(chan []byte)
	tunnels = make(map[*tunnel]bool)

	in = make(chan *tunnel)
	out = make(chan *tunnel)

	go func() {
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
	}()
}
