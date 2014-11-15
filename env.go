/**
 *  ------------------------------------------------------------
 *  @project	web.go
 *  @file       env.go
 *  @date       2014-11-14
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

// Env loads & parses the exported system environmental values via pre-defined struct.
package web

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// env processes all key/value under the prefix.
const Prefix = "GoAnywhere"

var Env *env = new(env)

type (
	env struct {
		Key   string
		Value string
		Name  string
		Type  string
	}
)

// Implements the Error interface.
func (self *env) Error() string {
	return fmt.Sprintf("Env.Load: <%s/%s> <%s => %s>", self.Key, self.Name, self.Value, self.Type)
}

// Load fetches the key/value pairs under the prefix into the given spec. struct.
func (self *env) Load(spec interface{}) error {
	s := reflect.ValueOf(spec).Elem()
	if s.Kind() != reflect.Struct {
		return errors.New("Configuration Spec. *MUST* be a struct.")
	}

	var stype reflect.Type = s.Type()
	var field reflect.Value

	for index := 0; index < s.NumField(); index++ {
		field = s.Field(index)
		if field.CanSet() {
			name := stype.Field(index).Tag.Get("web")
			if name == "" {
				name = stype.Field(index).Name
			}

			key := self.key(name)
			value := os.Getenv(key)
			if value == "" {
				continue
			}
			// converts the environmental value from string to its real type.
			// Supports: String | Bool | Float | Integer
			// TODO more formats
			switch field.Kind() {
			case reflect.String:
				field.SetString(value)
			case reflect.Bool:
				if val, err := strconv.ParseBool(value); err == nil {
					field.SetBool(val)
				} else {
					self.Key = key
					self.Value = value
					self.Name = name
					self.Type = field.Type().String()
					return self
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if val, err := strconv.ParseInt(value, 0, field.Type().Bits()); err == nil {
					field.SetInt(val)
				} else {
					self.Key = key
					self.Value = value
					self.Name = name
					self.Type = field.Type().String()
					return self
				}
			case reflect.Float32, reflect.Float64:
				if val, err := strconv.ParseFloat(value, field.Type().Bits()); err == nil {
					field.SetFloat(val)
				} else {
					self.Key = key
					self.Value = value
					self.Name = name
					self.Type = field.Type().String()
					return self
				}
			}
		}
	}
	return nil
}

// key constructs the real key for storing the name/value pair under prefix.
func (self *env) key(name string) string {
	return fmt.Sprintf("%s_%s", Prefix, strings.ToUpper(name))
}

// Get returns the value for the name under env. prefix.
func (self *env) Get(name string) string {
	return os.Getenv(self.key(name))
}

// Set sets the value for the name under env. prefix.
func (self *env) Set(name, value string) error {
	return os.Setenv(self.key(name), value)
}

// Values constructs [string]string map for key/value under env. prefix.
func (self *env) Values() map[string]string {
	environ := os.Environ()
	values := make(map[string]string)
	for _, pair := range environ {
		if strings.HasPrefix(pair, Prefix) {
			kv := strings.Split(pair, "=")
			if kv != nil && len(kv) >= 2 {
				values[kv[0]] = kv[1]
			}
		}
	}
	return values
}
