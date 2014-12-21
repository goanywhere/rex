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

package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGet(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	app := New()
	app.Get("/", func(ctx *Context) {
		ctx.String("hello")
	})
	app.Get("/404", func(ctx *Context) {
		ctx.WriteHeader(http.StatusNotFound)
		ctx.String("NotFound")
	})
	app.ServeHTTP(w, r)

	body := w.Body.String()
	status := w.Code
	Convey("http 200@hello GET test", t, func() {
		So(body, ShouldEqual, "hello")
		So(status, ShouldEqual, http.StatusOK)
	})

	app.Get("/404", func(ctx *Context) {
		ctx.WriteHeader(http.StatusNotFound)
		ctx.String("NotFound")
	})
	r, _ = http.NewRequest("GET", "/404", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, r)
	body = w.Body.String()
	status = w.Code
	Convey("http 400 GET test", t, func() {
		So(body, ShouldEqual, "NotFound")
		So(status, ShouldEqual, http.StatusNotFound)
	})
}

func TestPost(t *testing.T) {}

func TestPut(t *testing.T) {}

func TestDelete(t *testing.T) {}

func TestPatch(t *testing.T) {}

func TestHead(t *testing.T) {}

func TestOptions(t *testing.T) {}
