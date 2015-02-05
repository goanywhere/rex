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
package cache

import (
	"testing"
	"time"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/goanywhere/rex/internal"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCacheGet(t *testing.T) {
	settings := internal.Settings()
	settings.Set("cache.redis.servers", []string{"redis://127.0.0.1:6379"})
	cache := New()
	Convey("Redis Cache", t, func() {
		var name string
		cache.Get("notfound", &name)
		So(name, ShouldEqual, "")
	})
}

func TestCacheSet(t *testing.T) {
	settings := internal.Settings()
	settings.Set("cache.redis.servers", []string{"redis://127.0.0.1:6379"})
	cache := New()
	Convey("Redis Cache Set", t, func() {
		var test string
		cache.Set("test", "example", time.Second*30)
		cache.Get("test", &test)
		So(test, ShouldEqual, "example")
	})
}

func TestCacheDel(t *testing.T) {
	settings := internal.Settings()
	settings.Set("cache.redis.servers", []string{"redis://127.0.0.1:6379"})
	cache := New()
	Convey("Redis Cache Del", t, func() {
		var test string
		cache.Set("test", "example", time.Second*30)
		cache.Get("test", &test)
		So(test, ShouldEqual, "example")

		cache.Del("test")
		var notfound string
		err := cache.Get("test", &notfound)
		So(err, ShouldEqual, redigo.ErrNil)
		So(notfound, ShouldEqual, "")
	})
}

func TestCacheIncr(t *testing.T) {
	settings := internal.Settings()
	settings.Set("cache.redis.servers", []string{"redis://127.0.0.1:6379"})
	cache := New()
	Convey("Redis Cache Incr", t, func() {
		cache.Set("age", 100, time.Second*30)
		cache.Incr("age")
		var age int
		cache.Get("age", &age)
		So(age, ShouldEqual, 101)
	})
}

func TestCacheDecr(t *testing.T) {
	settings := internal.Settings()
	settings.Set("cache.redis.servers", []string{"redis://127.0.0.1:6379"})
	cache := New()
	Convey("Redis Cache Decr", t, func() {
		cache.Set("age", 100, time.Second*30)
		cache.Decr("age")
		var age int
		cache.Get("age", &age)
		So(age, ShouldEqual, 99)
	})
}

func TestCacheExists(t *testing.T) {
	settings := internal.Settings()
	settings.Set("cache.redis.servers", []string{"redis://127.0.0.1:6379"})
	cache := New()
	Convey("Redis Cache Exists", t, func() {
		So(cache.Exists("NotFound"), ShouldBeFalse)
	})
}

func TestCacheFlush(t *testing.T) {
	settings := internal.Settings()
	settings.Set("cache.redis.servers", []string{"redis://127.0.0.1:6379"})
	cache := New()
	Convey("Redis Cache Flush", t, func() {
		cache.Set("age", 100, time.Second*30)
		cache.Set("test", "example", time.Second*30)
		cache.Flush()
		So(cache.Exists("age"), ShouldBeFalse)
		So(cache.Exists("test"), ShouldBeFalse)
	})
}
