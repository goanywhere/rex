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
	"bytes"
	"log"

	"github.com/gorilla/websocket"
)

type tunnel struct {
	ws      *websocket.Conn
	message chan []byte
}

func (self *tunnel) connect() {
	go func() {
		for message := range self.message {
			if err := self.ws.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Failed to write message to peer: %v", err)
				break
			} else {
				log.Printf("[WebSocket][Write] %s", message)
				if bytes.Contains(message, []byte(`"command": "hello"`)) {
					log.Printf("[WebSocket] connection established")
				}
			}
		}
		self.ws.Close()
	}()

	for {
		_, message, err := self.ws.ReadMessage()
		if err != nil {
			break
		}
		switch true {
		case bytes.Contains(message, []byte(`"command":"hello"`)):
			//data, _ := json.Marshal(hello)
			//self.message <- data
			self.message <- []byte(`{
				"command": "hello",
				"protocols": [ "http://livereload.com/protocols/official-7" ],
				"serverName": "Rex#Livereload"
			}`)
		}
	}
	self.ws.Close()
}
