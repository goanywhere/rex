package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goanywhere/rex"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCompress(t *testing.T) {

	app := rex.New()
	app.Use(Compress)
	app.GET("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "app")
	})

	Convey("rex.middleware.Compress", t, func() {
		request, _ := http.NewRequest("GET", "/", nil)
		request.Header.Set("Accept-Encoding", "gzip")
		response := httptest.NewRecorder()

		app.ServeHTTP(response, request)

		So(response.Header().Get("Content-Encoding"), ShouldEqual, "gzip")
	})
}
