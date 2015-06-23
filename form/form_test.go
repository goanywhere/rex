package form

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type user struct {
	Username string `schema:"username"`
	Password string `schema:"password"`
}

func (self *user) Validate() error {
	if self.Username == "" || self.Password == "" {
		return errors.New("username/password can not be empty")
	}

	if len(self.Password) < 8 {
		return errors.New("password must be greater than 8bits")
	}
	return nil
}

func TestParse(t *testing.T) {
	app := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var form user
		if err := parse(r, &form); err == nil {
			w.WriteHeader(http.StatusAccepted)
			io.WriteString(w, "uid")
		} else {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, err.Error())
		}
	})

	Convey("rex.form.Parse", t, func() {
		values := url.Values{}
		values.Set("username", "username")
		values.Set("password", "password")

		request, _ := http.NewRequest("POST", "/", bytes.NewBufferString(values.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()

		app.ServeHTTP(response, request)
		So(response.Code, ShouldEqual, http.StatusAccepted)

		values.Set("password", "7")
		request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(values.Encode()))
		response = httptest.NewRecorder()

		app.ServeHTTP(response, request)
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})
}
