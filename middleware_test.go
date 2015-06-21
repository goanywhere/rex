package rex

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMiddleware(t *testing.T) {
	env := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Powered-By", "rex")
			next.ServeHTTP(w, r)
		})
	}

	json := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			w.Header().Set("Content-Type", "application/json")
		})
	}

	mw := new(middleware)
	mw.stack = append(mw.stack, env)
	mw.stack = append(mw.stack, json)
	Convey("rex.middleware", t, func() {

		request, _ := http.NewRequest("GET", "/", nil)
		response := httptest.NewRecorder()

		mw.ServeHTTP(response, request)

		So(response.Header().Get("X-Powered-By"), ShouldEqual, "rex")
		So(response.Header().Get("Content-Type"), ShouldEqual, "application/json")
	})
}
