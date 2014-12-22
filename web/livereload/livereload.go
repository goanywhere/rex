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
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	livereload = struct {
		tunnels   map[*tunnel]bool
		broadcast chan []byte

		register   chan *tunnel
		deregister chan *tunnel
	}{
		tunnels:   make(map[*tunnel]bool),
		broadcast: make(chan []byte),

		register:   make(chan *tunnel),
		deregister: make(chan *tunnel),
	}

	hello = struct {
		command    string   `json:"command"`
		protocols  []string `json:"protocols"`
		serverName string   `json:"serverName"`
	}{
		command: "hello",
		protocols: []string{
			"http://livereload.com/protocols/official-7",
			"http://livereload.com/protocols/official-8",
			"http://livereload.com/protocols/official-9",
			"http://livereload.com/protocols/2.x-origin-version-negotiation",
			"http://livereload.com/protocols/2.x-remote-control",
		},
		serverName: "Rex#livereload",
	}
)

type (
	alert struct {
		command string `json:"command"`
		message string `json:"message"`
	}

	reload struct {
		command string `json:"command"`
		path    string `json:"path"`
		liveCSS bool   `json:"liveCSS"`
	}
)

func Start() {
	go func() {
		for {
			select {
			case tunnel := <-livereload.register:
				livereload.tunnels[tunnel] = true

			case tunnel := <-livereload.deregister:
				delete(livereload.tunnels, tunnel)
				close(tunnel.message)

			case msg := <-livereload.broadcast:
				for tunnel := range livereload.tunnels {
					select {
					case tunnel.message <- msg:
					default:
						delete(livereload.tunnels, tunnel)
						close(tunnel.message)
					}
				}
			}
		}
	}()
}

func Alert(message string) {
	go func() {
		//msg := new(alert)
		//msg.command = "alert"
		//msg.message = message
		//data, _ := json.Marshal(msg)
		//livereload.broadcast <- data
		livereload.broadcast <- []byte(`{
			"command":"alert",
			"message":"` + message + "\"" + `
		}`)
	}()
}

func Reload(path string) {
	go func() {
		livereload.broadcast <- []byte(`{
			"command":"reload",
			"path":"` + path + "\"" + `,
			"originalPath":"",
			"liveCSS":true,
			"liveImg":true
		}`)
		//msg := new(reload)
		//msg.command = "reload"
		//msg.path = path
		//msg.liveCSS = true
		//data, _ := json.Marshal(msg)
		//livereload.broadcast <- data
	}()
}

func Serve(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	tn := new(tunnel)
	tn.ws = ws
	tn.message = make(chan []byte, 256)
	livereload.register <- tn
	defer func() { livereload.deregister <- tn }()
	tn.connect()
}

func ServeJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	w.Write(javascript)
}
