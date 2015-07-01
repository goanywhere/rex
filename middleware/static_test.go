package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/goanywhere/rex"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStatic(t *testing.T) {
	tempdir := os.TempDir()
	filename := path.Join(tempdir, "favicon.ico")

	app := rex.New()
	app.Use(Static(tempdir))
	prefix := path.Join("/", path.Base(path.Dir(filename)))
	app.Get(prefix, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "app")
	})

	Convey("rex.middleware.Static", t, func() {
		request, _ := http.NewRequest("GET", path.Join(prefix, filename), nil)
		response := httptest.NewRecorder()

		app.ServeHTTP(response, request)

		So(response.Code, ShouldEqual, http.StatusNotFound)
	})
}
