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
package internal

import (
	"log"
	"reflect"
	"strconv"
	"sync"

	"github.com/goanywhere/x/env"
)

var (
	once    sync.Once
	options *env.Env
)

func Options() *env.Env {
	once.Do(func() {
		options = env.New("rex")
		options.Set("port", 5000)
		options.Set("mode", "debug")
		options.Set("dir.static", "build")
		options.Set("dir.templates", "templates")
		options.Set("url.static", "/static/")
		// default environmental headers for modules.Env
		options.Set("header.x.ua.compatible", "deny")
		options.Set("header.x.frame.options", "nosniff")
		options.Set("header.x.xss.protection", "1; mode=block")
		options.Set("header.x.content.type.options", "IE=Edge,chrome=1")
		// session cookie defaults
		options.Set("session.cookie.name", "gsid")
		options.Set("session.cookie.maxage", 1209600)
		options.Set("session.cookie.httponly", true)
		options.Set("session.cookie.path", "/")
		options.Set("session.cookie.secure", false)
	})
	return options
}

// Option fetches the value from os's enviroment into the appointed address.
func Option(key string, ptr interface{}, fallback ...interface{}) {
	sv, exists := options.Get(key)
	if !exists && len(fallback) == 0 {
		return
	}

	if reflect.TypeOf(ptr).Kind() == reflect.Ptr {
		rv := reflect.ValueOf(ptr).Elem()
		switch rv.Kind() {
		case reflect.Bool:
			if v, e := strconv.ParseBool(sv); e == nil {
				rv.SetBool(v)
			} else if len(fallback) > 0 {
				rv.SetBool(fallback[0].(bool))
			}

		case reflect.Float32, reflect.Float64:
			if v, e := strconv.ParseFloat(sv, 64); e == nil {
				rv.SetFloat(v)
			} else if len(fallback) > 0 {
				rv.SetFloat(reflect.ValueOf(fallback[0]).Float())
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v, e := strconv.ParseInt(sv, 10, 64); e == nil {
				rv.SetInt(v)
			} else if len(fallback) > 0 {
				rv.SetInt(reflect.ValueOf(fallback[0]).Int())
			}

		case reflect.String:
			if sv == "" {
				rv.SetString(sv)
			} else if len(fallback) > 0 {
				rv.SetString(fallback[0].(string))
			}

		default:
			log.Fatalf("Failed to retrieve value for <%s>: unsupported type", key)
		}
	}
}
