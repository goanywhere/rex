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
		options.Set("port", "5000")
		options.Set("mode", "debug")
		options.Set("dir.static", "build")
		options.Set("dir.templates", "templates")
		options.Set("url.static", "/static/")
		// default environmental headers for modules.Env
		options.Set("header.x.ua.compatible", "deny")
		options.Set("header.x.frame.options", "nosniff")
		options.Set("header.x.xss.protection", "1; mode=block")
		options.Set("header.x.content.type.options", "IE=Edge,chrome=1")
	})
	return options
}

// TODO strings Array/Slice
func Option(key string, pointer interface{}) (e error) {
	T := reflect.TypeOf(pointer)
	if T.Kind() == reflect.Ptr {
		rv := reflect.ValueOf(pointer).Elem()
		switch rv.Kind() {
		case reflect.Bool:
			rv.SetBool(options.Bool(key))
		case reflect.Float32, reflect.Float64:
			rv.SetFloat(options.Float(key))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			rv.SetInt(options.Int(key))
		case reflect.String:
			rv.SetString(options.String(key))
		default:
			log.Printf("unknown type")
		}
	}
	return
}
