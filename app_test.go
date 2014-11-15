/**
 *  ------------------------------------------------------------
 *  @project
 *  @file       app_test.go
 *  @date       2014-11-15
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
	"net/http"
	"net/http/httptest"
	"testing"
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
	if body != "hello" {
		t.Errorf("Respsonse was not correctly parsed. Should be 'hello', but got %s.", body)
	}
	if status != http.StatusOK {
		t.Errorf("Expected Status: 200, Got: %d", status)
	}

	app.Get("/404", func(ctx *Context) {
		ctx.WriteHeader(http.StatusNotFound)
		ctx.String("NotFound")
	})
	r, _ = http.NewRequest("GET", "/404", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, r)
	body = w.Body.String()
	status = w.Code
	if status != http.StatusNotFound {
		t.Errorf("Expected Status: 404, Got: %d", status)
	}
	if body != "NotFound" {
		t.Errorf("Respsonse was not correctly parsed. Should be 'NotFound', but got %s.", body)
	}
}

func TestPost(t *testing.T) {}

func TestPut(t *testing.T) {}

func TestDelete(t *testing.T) {}

func TestPatch(t *testing.T) {}

func TestHead(t *testing.T) {}

func TestOptions(t *testing.T) {}
