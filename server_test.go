package rex

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/goanywhere/env"
	mw "github.com/goanywhere/rex/middleware"
	. "github.com/smartystreets/goconvey/convey"
)

func TestConfigure(t *testing.T) {
	Convey("rex.configure", t, func() {
		app := New()
		app.configure()

		So(debug, ShouldBeTrue)
		So(port, ShouldEqual, 5000)

		env.Set("PORT", 9394)
		app.configure()
		So(port, ShouldEqual, 5000)
	})
}

func TestBuild(t *testing.T) {
	Convey("rex.build", t, func() {
		app := New()
		app.build()

		So(len(app.middleware.stack), ShouldEqual, 1)

		app.Use(mw.NoCache)
		So(len(app.middleware.stack), ShouldEqual, 2)
	})
}

func TestRegister(t *testing.T) {
	Convey("rex.register", t, func() {
		app := New()
		app.register("/login", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusAccepted)
			w.Header().Set("X-Auth-Server", "rex")
		}, "POST")

		request, _ := http.NewRequest("POST", "/login", nil)
		response := httptest.NewRecorder()
		app.ServeHTTP(response, request)
		So(response.Code, ShouldEqual, http.StatusAccepted)
		So(response.Header().Get("X-Auth-Server"), ShouldEqual, "rex")
	})
}

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

func TestGET(t *testing.T) {
	app := New()
	app.GET("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Powered-By", "rex")
		w.Header().Set("Content-Type", "application/json")
	})

	Convey("rex.GET", t, func() {
		request, _ := http.NewRequest("GET", "/", nil)
		response := httptest.NewRecorder()

		app.ServeHTTP(response, request)

		So(response.Header().Get("X-Powered-By"), ShouldEqual, "rex")
		So(response.Header().Get("Content-Type"), ShouldEqual, "application/json")
	})
}

func TestPOST(t *testing.T) {
	app := New()
	app.POST("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Powered-By", "rex")
		w.Header().Set("Content-Type", "application/json")
	})

	Convey("rex.POST", t, func() {
		request, _ := http.NewRequest("POST", "/", nil)
		response := httptest.NewRecorder()

		app.ServeHTTP(response, request)

		So(response.Header().Get("X-Powered-By"), ShouldEqual, "rex")
		So(response.Header().Get("Content-Type"), ShouldEqual, "application/json")
	})
}

func TestPUT(t *testing.T) {
	app := New()
	app.PUT("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Powered-By", "rex")
		w.Header().Set("Content-Type", "application/json")
	})

	Convey("rex.PUT", t, func() {
		request, _ := http.NewRequest("PUT", "/", nil)
		response := httptest.NewRecorder()

		app.ServeHTTP(response, request)

		So(response.Header().Get("X-Powered-By"), ShouldEqual, "rex")
		So(response.Header().Get("Content-Type"), ShouldEqual, "application/json")
	})
}

func TestDELETE(t *testing.T) {
	app := New()
	app.DELETE("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Powered-By", "rex")
		w.Header().Set("Content-Type", "application/json")
	})

	Convey("rex.DELETE", t, func() {
		request, _ := http.NewRequest("DELETE", "/", nil)
		response := httptest.NewRecorder()

		app.ServeHTTP(response, request)

		So(response.Header().Get("X-Powered-By"), ShouldEqual, "rex")
		So(response.Header().Get("Content-Type"), ShouldEqual, "application/json")
	})
}

func TestCONNECT(t *testing.T) {
	app := New()
	app.CONNECT("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Powered-By", "rex")
	})

	Convey("rex.CONNECT", t, func() {
		request, _ := http.NewRequest("CONNECT", "/", nil)
		response := httptest.NewRecorder()

		app.ServeHTTP(response, request)

		So(response.Header().Get("X-Powered-By"), ShouldEqual, "rex")
	})
}

func TestTRACE(t *testing.T) {
	app := New()
	app.TRACE("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Powered-By", "rex")
	})

	Convey("rex.TRACE", t, func() {
		request, _ := http.NewRequest("TRACE", "/", nil)
		response := httptest.NewRecorder()

		app.ServeHTTP(response, request)

		So(response.Header().Get("X-Powered-By"), ShouldEqual, "rex")
	})
}

func TestGroup(t *testing.T) {
	app := New()
	app.GET("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "index")
	})
	user := app.Group("/users")
	user.GET("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Powered-By", "rex")
	})

	Convey("rex.Group", t, func() {
		request, _ := http.NewRequest("GET", "/users/", nil)
		response := httptest.NewRecorder()

		app.ServeHTTP(response, request)

		So(response.Header().Get("X-Powered-By"), ShouldEqual, "rex")
	})
}

func TestFileServer(t *testing.T) {
	Convey("rex.FileServer", t, func() {
		var (
			prefix   = "/assets/"
			filename = "logo.png"
		)
		tempdir := os.TempDir()
		filepath := path.Join(tempdir, filename)
		os.Create(filepath)
		defer os.Remove(filepath)

		app := New()
		app.FileServer(prefix, tempdir)

		request, _ := http.NewRequest("GET", path.Join(prefix, filename), nil)
		response := httptest.NewRecorder()
		app.ServeHTTP(response, request)
		So(response.Code, ShouldEqual, http.StatusOK)

		filename = "index.html"
		request, _ = http.NewRequest("HEAD", prefix, nil)
		response = httptest.NewRecorder()
		app.ServeHTTP(response, request)
		So(response.Code, ShouldEqual, http.StatusOK)
	})
}

func TestUse(t *testing.T) {
	app := New()
	app.GET("/", func(w http.ResponseWriter, r *http.Request) {
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
		app.GET("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
			vars := app.Vars(r)
			So(vars["id"], ShouldEqual, "123")
		})

		request, _ := http.NewRequest("GET", "/users/123", nil)
		response := httptest.NewRecorder()
		app.ServeHTTP(response, request)
	})
}
