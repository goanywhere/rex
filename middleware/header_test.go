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
	var values = make(map[string]string)
	values["X-Powered-By"] = "rex server"

	app := rex.New()
	app.Use(Header(values))
	app.GET("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "app")
	})

	Convey("rex.middleware.Header", t, func() {
		request, _ := http.NewRequest("GET", "/", nil)
		response := httptest.NewRecorder()

		app.ServeHTTP(response, request)

		So(response.Header().Get("X-Powered-By"), ShouldEqual, "rex server")
	})
}
