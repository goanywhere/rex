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

	URL = struct {
		WebSocket  string
		JavaScript string
	}{
		WebSocket:  "/livereload",
		JavaScript: "/livereload.js",
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
			Path:    URL.WebSocket,
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
func ServeWebSocket(w http.ResponseWriter, r *http.Request) {
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

// ServeJavaScript serves livereload.js for browser.
func ServeJavaScript(w http.ResponseWriter, r *http.Request) {
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
