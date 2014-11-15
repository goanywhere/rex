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
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	Env *env
	// env processes all key/value under the prefix.
	Prefix string = "GoAnywhere"
)

type (
	env struct {
		pattern *regexp.Regexp
		space   *regexp.Regexp

		Key   string // The actual key to store in os.Environ.
		Value string // String value of the storage.
		Name  string // User's specified field name.
		Type  string // The actual type of the value.
	}
)

// Implements the Error interface.
func (self *env) Error() string {
	return fmt.Sprintf("Env.Load: <%s/%s> <%s => %s>", self.Key, self.Name, self.Value, self.Type)
}

// ---------------------------------------------------------------------------
//  Internal Helpers
// ---------------------------------------------------------------------------
// key constructs the real key for storing the name/value pair under prefix.
func (self *env) key(name string) string {
	return fmt.Sprintf("%s_%s", Prefix, strings.ToUpper(name))
}

// findKeyValue finds ':' or '=' separated key/value pair from the given string.
func (self *env) findKeyValue(str string) (key, value string) {
	result := self.pattern.FindString(str)
	if result != "" {
		raw := self.space.ReplaceAllString(result, "")
		var kv []string
		if strings.Index(raw, ":") >= 0 {
			kv = strings.Split(raw, ":")
		} else {
			kv = strings.Split(raw, "=")
		}
		key = kv[0]
		value = kv[1]
	}
	return
}

// Load fetches the values from '.env' from project's CWD.
func (self *env) Load() error {
	if file, err := os.Open(filepath.Join(self.Get("root"), ".env")); err == nil {
		defer file.Close()
		reader := bufio.NewReader(file)
		for {
			line, e := reader.ReadString('\n')
			if e != nil || e == io.EOF {
				return e
			}
			k, v := self.findKeyValue(line)
			if k != "" && v != "" {
				self.Set(k, v)
			}
		}
	} else {
		return err
	}
	return nil
}

// LoadObject fetches the key/value pairs under the prefix into the given spec. struct.
func (self *env) LoadInto(spec interface{}) error {
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
				if val, err := ToBool(value); err == nil {
					field.SetBool(val)
				} else {
					self.Key = key
					self.Value = value
					self.Name = name
					self.Type = field.Type().String()
					return self
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				// if val, err := ToInt(value); err == nil {
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
				if val, err := ToFloat(value); err == nil {
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

// Get returns the value for the name under env. prefix.
func (self *env) Get(name string) string {
	return os.Getenv(self.key(name))
}

// GetBool returns & parses the stored string value to bool.
func (self *env) GetBool(name string) (bool, error) {
	return ToBool(self.Get(name))
}

// GetFloat returns & parsed the stored string value to int.
func (self *env) GetFloat(name string) (float64, error) {
	return ToFloat(self.Get(name))
}

// GetInt returns & parsed the stored string value to int.
func (self *env) GetInt(name string) (int, error) {
	return ToInt(self.Get(name))
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

func init() {
	Env = new(env)
	Env.pattern = regexp.MustCompile(`(\w+)\s*(:|=)\s*(\w+)`)
	Env.space = regexp.MustCompile(`\s`)
}
