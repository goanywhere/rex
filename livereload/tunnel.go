package livereload

import (
	"encoding/json"
	"regexp"

	"github.com/gorilla/websocket"
)

var regexHandshake = regexp.MustCompile(`"command"\s*:\s*"hello"`)

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
				if regexHandshake.Find(message) != nil {
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
		case regexHandshake.Find(message) != nil:
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
