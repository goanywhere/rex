package livereload

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/goanywhere/env"
)

type writer struct {
	http.ResponseWriter
	host string
}

func (self *writer) addJavaScript(data []byte) []byte {
	javascript := fmt.Sprintf(`<script src="//%s%s"></script>
</head>`, self.host, URL.JavaScript)
	return regexp.MustCompile(`</head>`).ReplaceAll(data, []byte(javascript))
}

func (self *writer) Write(data []byte) (size int, e error) {
	if strings.Contains(self.Header().Get("Content-Type"), "html") {
		var encoding = self.Header().Get("Content-Encoding")
		if encoding == "" {
			data = self.addJavaScript(data)
		} else {
			var reader io.ReadCloser
			var buffer *bytes.Buffer = new(bytes.Buffer)

			if encoding == "gzip" {
				// decode to add javascript reference.
				reader, _ = gzip.NewReader(bytes.NewReader(data))
				io.Copy(buffer, reader)
				output := self.addJavaScript(buffer.Bytes())
				reader.Close()
				buffer.Reset()
				// encode back to HTML with added javascript reference.
				writer := gzip.NewWriter(buffer)
				writer.Write(output)
				writer.Close()
				data = buffer.Bytes()

			} else if encoding == "deflate" {
				// decode to add javascript reference.
				reader, _ = zlib.NewReader(bytes.NewReader(data))
				io.Copy(buffer, reader)
				output := self.addJavaScript(buffer.Bytes())
				reader.Close()
				buffer.Reset()
				// encode back to HTML with added javascript reference.
				writer := zlib.NewWriter(buffer)
				writer.Write(output)
				writer.Close()
				data = buffer.Bytes()
			}
		}
	}
	return self.ResponseWriter.Write(data)
}

func Module(next http.Handler) http.Handler {
	// ONLY run this under debug mode.
	if !env.Bool("DEBUG", true) {
		return next
	}
	Start()
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == URL.WebSocket {
			ServeWebSocket(w, r)

		} else if r.URL.Path == URL.JavaScript {
			ServeJavaScript(w, r)

		} else {
			writer := &writer{w, r.Host}
			next.ServeHTTP(writer, r)
		}
	}
	return http.HandlerFunc(fn)
}
