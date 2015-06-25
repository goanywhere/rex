package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goanywhere/rex"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNoCache(t *testing.T) {
	app := rex.New()
	app.Use(NoCache)
	app.GET("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "app")
	})

	Convey("rex.middleware.NoCache", t, func() {
		request, _ := http.NewRequest("GET", "/", nil)
		response := httptest.NewRecorder()

		app.ServeHTTP(response, request)

		header := response.Header()
		So(header.Get("Cache-Control"), ShouldEqual, "no-cache, no-store, must-revalidate")
		So(header.Get("Pragma"), ShouldEqual, "no-cache")
		So(header.Get("Expires"), ShouldEqual, "0")
	})
}
