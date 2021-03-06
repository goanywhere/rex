package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goanywhere/rex"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHeader(t *testing.T) {
	var values = make(http.Header)
	values.Set("X-Powered-By", "rex")

	app := rex.New()
	app.Use(Header(values))
	app.Get("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "app")
	})

	Convey("rex.middleware.Header", t, func() {
		request, _ := http.NewRequest("GET", "/", nil)
		response := httptest.NewRecorder()

		app.ServeHTTP(response, request)

		So(response.Header().Get("X-Powered-By"), ShouldEqual, "rex")
	})
}
