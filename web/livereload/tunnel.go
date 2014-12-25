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
	"regexp"

	"github.com/gorilla/websocket"
)

var handshake = regexp.MustCompile(`"command"\s*:\s*"hello"`)

/* ----------------------------------------------------------------------
 * WebSocket Server Tunnel
 * ----------------------------------------------------------------------*/
type tunnel struct {
	socket  *websocket.Conn
	message chan []byte
}

// connect reads/writes message for livereload.js.
func (self *tunnel) connect() {
	// ***********************
	// WebSocket Tunnel#Write
	// ***********************
	go func() {
		for message := range self.message {
			if err := self.socket.WriteMessage(websocket.TextMessage, message); err != nil {
				break
			} else {
				if handshake.Find(message) != nil {
					// Keep the tunnel opened after handshake(hello command).
					Reload()
				}
			}
		}
		self.socket.Close()
	}()
	// ***********************
	// WebSocket Tunnel#Read
	// ***********************
	for {
		_, message, err := self.socket.ReadMessage()
		if err != nil {
			break
		}
		switch true {
		case handshake.Find(message) != nil:
			var bytes, _ = json.Marshal(&hello{
				Command:    "hello",
				Protocols:  []string{"http://livereload.com/protocols/official-7"},
				ServerName: "Rex#Livereload",
			})
			self.message <- bytes
		}
	}
	self.socket.Close()
}
