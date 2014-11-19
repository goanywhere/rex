/**
 *  ------------------------------------------------------------
 *  @project	web.go
 *  @file       context_test.go
 *  @date       2014-11-18
 *  @author     Jim Zhan <jim.zhan@me.com>
 *
 *  Copyright © 2014 Jim Zhan.
 *  ------------------------------------------------------------
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
 *  ------------------------------------------------------------
 */
package web

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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
	name, value := "number", 123456789
	Convey("Context Cookie", t, func() {
		request, _ := http.NewRequest("GET", "/", nil)
		writer := httptest.NewRecorder()
		deleteCookie(writer, name)
		app := New()
		app.Get("/set", func(ctx *Context) {
			ctx.SetCookie(name, value, nil)
		})
		app.Get("/get", func(ctx *Context) {
			So(ctx.Cookie(name).(int), ShouldEqual, value)
		})
		app.ServeHTTP(writer, request)

	})
}

func TestSetCookie(t *testing.T) {
	Convey("Context Set Cookie", t, func() {
		setup(func(ctx *Context) {
			name, value := "number", 1234567890
			src, _ := Serialize(value)

			ctx.SetCookie(name, value, nil)

			So(ctx.Header().Get("Set-Cookie"), ShouldEqual, fmt.Sprintf("%s=%s", name, src))
		})
	})
}

func TestSecureCookie(t *testing.T) {
	name, value := "number", 1234567890
	Convey("[contex#SecureCookie]", t, func() {
		request, _ := http.NewRequest("GET", "/", nil)
		writer := httptest.NewRecorder()
		app := New()
		app.Get("/set", func(ctx *Context) {
			deleteCookie(writer, name)
			ctx.SetSecureCookie(name, value, nil)
		})
		app.Get("/get", func(ctx *Context) {
			So(32423, ShouldEqual, value)
			t.Logf("Name (%s): %d", name, ctx.SecureCookie(name).(int))
		})
		app.ServeHTTP(writer, request)
	})
}
