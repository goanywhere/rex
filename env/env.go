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
package env

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// env processes all key/value under the prefix.
const Prefix string = "GoAnywhere"

var (
	pattern *regexp.Regexp
	space   *regexp.Regexp
)

// ---------------------------------------------------------------------------
//  Internal Helpers
// ---------------------------------------------------------------------------
// getKey constructs the real key for storing the name/value pair under prefix.
func getKey(name string) string {
	return fmt.Sprintf("%s_%s", Prefix, strings.ToUpper(name))
}

// findKeyValue finds ':' or '=' separated key/value pair from the given string.
func findKeyValue(str string) (key, value string) {
	result := pattern.FindString(str)
	if result != "" {
		raw := space.ReplaceAllString(result, "")
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

// ---------------------------------------------------------------------------
//  Public APIs
// ---------------------------------------------------------------------------
// Load fetches the values from '.env' from project's CWD.
// *NOTE* value *MUST* not include ":" or "=".
func Load() error {
	if dotenv, err := os.Open(filepath.Join(Get("root"), ".env")); err == nil {
		defer dotenv.Close()
		reader := bufio.NewReader(dotenv)
		for {
			var line string
			line, err = reader.ReadString('\n')
			if err != nil || err == io.EOF {
				return err
			}
			k, v := findKeyValue(line)
			if k != "" && v != "" {
				Set(k, v)
			}
		}
	} else {
		return err
	}
	return nil
}

// LoadObject fetches the key/value pairs under the prefix into the given spec. struct.
func LoadInto(spec interface{}) error {
	s := reflect.ValueOf(spec).Elem()
	if s.Kind() != reflect.Struct {
		return fmt.Errorf("Configuration Spec. *MUST* be a struct.")
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

			key := getKey(name)
			value := os.Getenv(key)
			if value == "" {
				continue
			}
			// converts the environmental value from string to its real type.
			// Supports: String | Bool | Float | Integer
			// TODO Complex Object?
			switch field.Kind() {
			case reflect.String:
				field.SetString(value)
			case reflect.Bool:
				if val, err := ToBool(value); err == nil {
					field.SetBool(val)
				} else {
					return fmt.Errorf("env.LoadInto: <%s (%s): %s>", name, field.Type().String(), value)
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				// if val, err := ToInt(value); err == nil {
				if val, err := strconv.ParseInt(value, 0, field.Type().Bits()); err == nil {
					field.SetInt(val)
				} else {
					return fmt.Errorf("env.LoadInto: <%s (%s): %s>", name, field.Type().String(), value)
				}
			case reflect.Float32, reflect.Float64:
				if val, err := ToFloat(value); err == nil {
					field.SetFloat(val)
				} else {
					return fmt.Errorf("env.LoadInto: <%s (%s): %s>", name, field.Type().String(), value)
				}
			}
		}
	}
	return nil
}

// Get returns the value for the name under env. prefix.
func Get(name string) string {
	return os.Getenv(getKey(name))
}

// GetBool returns & parses the stored string value to bool.
func GetBool(name string) (bool, error) {
	return ToBool(Get(name))
}

// GetFloat returns & parsed the stored string value to int.
func GetFloat(name string) (float64, error) {
	return ToFloat(Get(name))
}

// GetInt returns & parsed the stored string value to int.
func GetInt(name string) (int, error) {
	return ToInt(Get(name))
}

// Set sets the value for the name under env. prefix.
func Set(name, value string) error {
	return os.Setenv(getKey(name), value)
}

// Values constructs [string]string map for key/value under env. prefix.
func Values() map[string]string {
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
	// catch almost any printable characters expect "=" & ":".
	pattern = regexp.MustCompile(`(\w+)\s*(:|=)\s*([[:graph:]]+)`)
	space = regexp.MustCompile(`\s`)

	// Use CWD as fallback.
	//if Get("root") == "" {
	//cwd, _ := os.Getwd()
	//Set("root", cwd)
	//}
}
