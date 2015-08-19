package rex

import (
	"bytes"
	"encoding/json"
	"net/http"
	"path"

	"github.com/goanywhere/env"
	"github.com/goanywhere/fs"
)

// Shortcut for string based map.
type M map[string]interface{}

// Sends the HTTP response in JSON.
func Send(w http.ResponseWriter, v interface{}) {
	var buffer = new(bytes.Buffer)
	defer func() {
		buffer.Reset()
	}()

	w.Header().Set("Content-type", "application/json; charset=utf-8")
	if err := json.NewEncoder(buffer).Encode(v); err == nil {
		w.Write(buffer.Bytes())
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Renders a view (Pongo) and sends the rendered HTML string to the client.
// Optional parameters: value, local variables for the view.
// func Render(filename string, v ...interface{}) {}

func init() {
	var basedir = fs.Getcd(2)
	env.Set("basedir", basedir)
	env.Load(path.Join(basedir, ".env"))
}
