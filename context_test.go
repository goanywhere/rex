/* ----------------------------------------------------------------------
 *       ______      ___                         __
 *      / ____/___  /   |  ____  __  ___      __/ /_  ___  ________
 *     / / __/ __ \/ /| | / __ \/ / / / | /| / / __ \/ _ \/ ___/ _ \
 *    / /_/ / /_/ / ___ |/ / / / /_/ /| |/ |/ / / / /  __/ /  /  __/
 *    \____/\____/_/  |_/_/ /_/\__. / |__/|__/_/ /_/\___/_/   \___/
 *                            /____/
 *
 * (C) Copyright 2014 GoAnywhere (http://goanywhere.io).
 * ----------------------------------------------------------------------
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 * ----------------------------------------------------------------------*/

package web

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goanywhere/web/crypto"
	"github.com/goanywhere/web/env"
	. "github.com/smartystreets/goconvey/convey"
)

func setup(handler HandlerFunc) {
	request, _ := http.NewRequest("GET", "/", nil)
	writer := httptest.NewRecorder()
	app := New()
	app.Get("/", handler)
	app.ServeHTTP(writer, request)
}

func deleteCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{Name: name, Path: "/", MaxAge: -1})
}

func TestCookie(t *testing.T) {
	Convey("context#Cookie", t, func() {
		cookie := &http.Cookie{Name: "number", Value: "123", Path: "/"}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := NewContext(w, r)
			ctx.String(ctx.Cookie(cookie.Name))
		}))
		defer server.Close()

		client := new(http.Client)
		request, _ := http.NewRequest("GET", server.URL, nil)
		request.AddCookie(cookie)

		response, _ := client.Do(request)
		defer response.Body.Close()

		body, _ := ioutil.ReadAll(response.Body)
		So(string(body), ShouldEqual, cookie.Value)
	})
}

func TestSetCookie(t *testing.T) {
	Convey("context#SetCookie", t, func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := NewContext(w, r)
			ctx.SetCookie(&http.Cookie{Name: "number", Value: "123", Path: "/"})
			ctx.String("Hello Cookie")
			return
		}))
		defer server.Close()

		if response, err := http.Get(server.URL); err == nil {
			cookie := response.Cookies()[0]
			So(cookie.Name, ShouldEqual, "number")
			So(cookie.Value, ShouldEqual, "123")
		}
	})
}

func TestSecureCookie(t *testing.T) {
	Convey("[contex#SecureCookie]", t, func() {
		env.Set("secret", crypto.RandomString(32, nil))
		// Ensure we use the same signature as context does.
		signature = crypto.NewSignature(env.Get("secret"))

		name, value := "number", "1234567890"
		src, _ := signature.Encode(name, []byte(value))
		cookie := &http.Cookie{
			Name:  name,
			Value: src,
			Path:  "/",
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := NewContext(w, r)
			ctx.String(ctx.SecureCookie(name))
		}))
		defer server.Close()

		client := new(http.Client)
		request, _ := http.NewRequest("GET", server.URL, nil)
		request.AddCookie(cookie)

		response, _ := client.Do(request)
		defer response.Body.Close()

		body, _ := ioutil.ReadAll(response.Body)
		So(string(body), ShouldEqual, value)
	})
}

func TestSetSecureCookie(t *testing.T) {
	Convey("context#SetSecureCookie", t, func() {
		env.Set("secret", crypto.RandomString(32, nil))
		// Ensure we use the same signature as context does.
		signature = crypto.NewSignature(env.Get("secret"))

		name, value := "number", "123"

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := NewContext(w, r)
			ctx.SetSecureCookie(&http.Cookie{Name: name, Value: value, Path: "/"})
			ctx.String("Hello Cookie")
		}))
		defer server.Close()

		if response, err := http.Get(server.URL); err == nil {
			cookie := response.Cookies()[0]
			bits, _ := signature.Decode(cookie.Name, cookie.Value)

			So(cookie.Name, ShouldEqual, name)
			So(string(bits), ShouldEqual, value)
		}
	})
}
