/* ----------------------------------------------------------------------
 *       ______      ___                         __
 *      / ____/___  /   |  ____  __  ___      __/ /_  ___  ________
 *     / / __/ __ \/ /| | / __ \/ / / / | /| / / __ \/ _ \/ ___/ _ \
 *    / /_/ / /_/ / ___ |/ / / / /_/ /| |/ |/ / / / /  __/ /  /  __/
 *    \____/\____/_/  |_/_/ /_/\__. / |__/|__/_/ /_/\___/_/   \___/
 *                            /____/
 *
 * (C) Copyright 2015 GoAnywhere (http://goanywhere.io).
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
package rex

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	cookie "github.com/gorilla/securecookie"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCookie(t *testing.T) {
	Convey("Context.Cookie", t, func() {
		a := cookie.GenerateRandomKey(64)
		b := cookie.GenerateRandomKey(32)
		securecookie = cookie.New(a, b)

		Convey("Get", func() {
			raw, _ := securecookie.Encode("number", 123)
			cookie := &http.Cookie{Name: "number", Value: raw, Path: "/"}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var value int
				ctx := NewContext(w, r)
				ctx.Cookie("number", &value)
				ctx.Write([]byte(fmt.Sprintf("%d", value)))
			}))
			defer server.Close()

			client := new(http.Client)
			request, _ := http.NewRequest("GET", server.URL, nil)
			request.AddCookie(cookie)

			response, _ := client.Do(request)
			defer response.Body.Close()

			body, _ := ioutil.ReadAll(response.Body)
			So(string(body), ShouldEqual, "123")
		})

		Convey("Set", func() {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				values := make(map[string]interface{})
				values["uid"] = 123
				values["username"] = "test@example.com"
				values["anonymous"] = false

				ctx := NewContext(w, r)
				ctx.SetCookie("user", values)
				return
			}))
			defer server.Close()

			if response, err := http.Get(server.URL); err == nil {
				So(len(response.Cookies()), ShouldEqual, 1)

				if len(response.Cookies()) == 1 {
					cookie := response.Cookies()[0]
					var values map[string]interface{}
					securecookie.Decode(cookie.Name, cookie.Value, &values)

					So(values["uid"], ShouldEqual, 123)
					So(values["username"], ShouldEqual, "test@example.com")
					So(values["anonymous"], ShouldBeFalse)
				}
			}
		})
	})
}

func TestContextRender(t *testing.T) {
	Convey("Context.Render", t, func() {

		Convey("HTML", func() {
			tmp := os.TempDir()
			filename := path.Join(tmp, "index.html")
			ioutil.WriteFile(filename, []byte("<html><body>{{ user.Username }}</body></html>"), os.ModePerm)

			loadViews(tmp)

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				type User struct {
					Username string
				}
				user := &User{Username: "test@example.com"}

				ctx := NewContext(w, r)
				ctx.Set("user", user)
				ctx.Render("index.html")
			}))
			defer server.Close()

			if response, err := http.Get(server.URL); err == nil {
				So(response.StatusCode, ShouldEqual, http.StatusOK)
				So(response.Header.Get("Content-Type"), ShouldStartWith, "text/html")

				bytes, err := ioutil.ReadAll(response.Body)
				So(err, ShouldBeNil)
				So(bytes, ShouldNotBeNil)
				So(string(bytes), ShouldContainSubstring, "test@example.com")

				os.Remove(filename)
			}

		})

		Convey("XML", func() {
			type A struct {
				Id   int    `xml:"id,attr"`
				Name string `xml:"name"`
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var obj = &A{Id: 123, Name: "test"}

				ctx := NewContext(w, r)
				ctx.Render(obj)
			}))
			defer server.Close()

			if response, err := http.Get(server.URL); err == nil {
				So(response.Header.Get("Content-Type"), ShouldStartWith, "application/xml")
				var obj A
				xml.NewDecoder(response.Body).Decode(&obj)
				So(obj.Id, ShouldEqual, 123)
				So(obj.Name, ShouldEqual, "test")
			}
		})

		Convey("JSON", func() {
			type B struct {
				Id   int    `json:"id"`
				Name string `json:"name"`
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var obj = &B{Id: 123, Name: "test"}
				ctx := NewContext(w, r)
				ctx.Render(obj)
			}))
			defer server.Close()

			if response, err := http.Get(server.URL); err == nil {
				So(response.Header.Get("Content-Type"), ShouldStartWith, "application/json")
				var obj B
				json.NewDecoder(response.Body).Decode(&obj)
				So(obj.Id, ShouldEqual, 123)
				So(obj.Name, ShouldEqual, "test")
			}
		})
	})
}
