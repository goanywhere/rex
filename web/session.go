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
package web

import "reflect"

type Session interface {
	Clear()

	Del(key string)

	Get(key string, ptr interface{}) error

	Set(key string, value interface{})

	Save() error
}

// internal securecookie based session.
type session struct {
	ctx    *Context
	values map[string]interface{}
}

func (self *session) Clear() {
	for key, _ := range self.values {
		delete(self.values, key)
	}
}

func (self *session) Del(key string) {
	delete(self.values, key)
}

func (self *session) Get(key string, ptr interface{}) error {
	if value := self.values[key]; value != nil {
		if reflect.TypeOf(ptr).Kind() == reflect.Ptr {
			elem := reflect.ValueOf(ptr).Elem()
			elem.Set(reflect.ValueOf(value))
		}
	}
	return nil
}

func (self *session) Set(key string, value interface{}) {
	self.values[key] = value
}

func (self *session) Save() error {
	return self.ctx.SetSecureCookie(settings.Cookie.Name, self.values)
}
