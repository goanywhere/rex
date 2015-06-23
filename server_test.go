package rex

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAny(t *testing.T) {
	app := New()
	app.Any("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Powered-By", "rex")
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)

		case "POST":
			w.WriteHeader(http.StatusCreated)

		case "PUT":
			w.WriteHeader(http.StatusAccepted)

		case "DELETE":
			w.WriteHeader(http.StatusGone)

		default:
			w.Header().Set("X-HTTP-Method", r.Method)
		}
	})

	Convey("rex.Any", t, func() {
		var (
			request  *http.Request
			response *httptest.ResponseRecorder
		)
		request, _ = http.NewRequest("GET", "/", nil)
		response = httptest.NewRecorder()
		app.ServeHTTP(response, request)
		So(response.Code, ShouldEqual, http.StatusOK)

		request, _ = http.NewRequest("POST", "/", nil)
		response = httptest.NewRecorder()
		app.ServeHTTP(response, request)
		So(response.Code, ShouldEqual, http.StatusCreated)

		request, _ = http.NewRequest("PUT", "/", nil)
		response = httptest.NewRecorder()
		app.ServeHTTP(response, request)
		So(response.Code, ShouldEqual, http.StatusAccepted)

		request, _ = http.NewRequest("DELETE", "/", nil)
		response = httptest.NewRecorder()
		app.ServeHTTP(response, request)
		So(response.Code, ShouldEqual, http.StatusGone)
	})
}

func TestGet(t *testing.T) {
	app := New()
	app.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Powered-By", "rex")
		w.Header().Set("Content-Type", "application/json")
	})

	Convey("rex.Get", t, func() {
		request, _ := http.NewRequest("GET", "/", nil)
		response := httptest.NewRecorder()

		app.ServeHTTP(response, request)

		So(response.Header().Get("X-Powered-By"), ShouldEqual, "rex")
		So(response.Header().Get("Content-Type"), ShouldEqual, "application/json")
	})
}

func TestPost(t *testing.T) {
	app := New()
	app.Post("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Powered-By", "rex")
		w.Header().Set("Content-Type", "application/json")
	})

	Convey("rex.Post", t, func() {
		request, _ := http.NewRequest("POST", "/", nil)
		response := httptest.NewRecorder()

		app.ServeHTTP(response, request)

		So(response.Header().Get("X-Powered-By"), ShouldEqual, "rex")
		So(response.Header().Get("Content-Type"), ShouldEqual, "application/json")
	})
}

func TestPut(t *testing.T) {
	app := New()
	app.Put("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Powered-By", "rex")
		w.Header().Set("Content-Type", "application/json")
	})

	Convey("rex.Put", t, func() {
		request, _ := http.NewRequest("PUT", "/", nil)
		response := httptest.NewRecorder()

		app.ServeHTTP(response, request)

		So(response.Header().Get("X-Powered-By"), ShouldEqual, "rex")
		So(response.Header().Get("Content-Type"), ShouldEqual, "application/json")
	})
}

func TestDelete(t *testing.T) {
	app := New()
	app.Delete("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Powered-By", "rex")
		w.Header().Set("Content-Type", "application/json")
	})

	Convey("rex.Delete", t, func() {
		request, _ := http.NewRequest("DELETE", "/", nil)
		response := httptest.NewRecorder()

		app.ServeHTTP(response, request)

		So(response.Header().Get("X-Powered-By"), ShouldEqual, "rex")
		So(response.Header().Get("Content-Type"), ShouldEqual, "application/json")
	})
}

func TestGroup(t *testing.T) {
	app := New()
	app.Get("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "index")
	})
	user := app.Group("/users")
	user.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Powered-By", "rex")
	})

	Convey("rex.Group", t, func() {
		request, _ := http.NewRequest("GET", "/users/", nil)
		response := httptest.NewRecorder()

		app.ServeHTTP(response, request)

		So(response.Header().Get("X-Powered-By"), ShouldEqual, "rex")
	})
}

func TestUse(t *testing.T) {
	app := New()
	app.Get("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "index")
	})
	app.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})
	Convey("rex.Use", t, func() {
		request, _ := http.NewRequest("GET", "/", nil)
		response := httptest.NewRecorder()
		app.ServeHTTP(response, request)
		So(response.Header().Get("Content-Type"), ShouldEqual, "application/json")
	})
}

func TestVars(t *testing.T) {
	Convey("rex.Vars", t, func() {
		app := New()
		app.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
			vars := app.Vars(r)
			So(vars["id"], ShouldEqual, "123")
		})

		request, _ := http.NewRequest("GET", "/users/123", nil)
		response := httptest.NewRecorder()
		app.ServeHTTP(response, request)
	})
}
