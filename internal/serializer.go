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
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
	"strconv"
)

// Serializes the given value into byte array, primitive types
// conversions use reflect while complex types via package gob.
func Serialize(value interface{}) (bits []byte, err error) {
	switch v := value.(type) {
	case bool:
		bits = []byte(strconv.FormatBool(v))

	case float32, float64:
		bits = []byte(strconv.FormatFloat(reflect.ValueOf(v).Float(), 'g', -1, 64))

	case int, int8, int16, int32, int64:
		bits = []byte(strconv.FormatInt(reflect.ValueOf(v).Int(), 10))

	case uint, uint8, uint16, uint32, uint64:
		bits = []byte(strconv.FormatUint(reflect.ValueOf(v).Uint(), 10))

	case string:
		bits = []byte(v)

	default:
		buffer := bytes.NewBuffer(bits)
		err = gob.NewEncoder(buffer).Encode(v)
		bits = buffer.Bytes()
	}
	return
}

// Deserializes the given bytes into appointed address.
// Primitive types deserialized using reflect directly while
// complex types utilizes package gob.
func Deserialize(bits []byte, ptr interface{}) error {
	if value := reflect.ValueOf(ptr); value.Kind() == reflect.Ptr {
		switch elem := value.Elem(); elem.Kind() {
		case reflect.Bool:
			v, e := strconv.ParseBool(string(bits))
			elem.SetBool(v)
			return e

		case reflect.Float32, reflect.Float64:
			v, e := strconv.ParseFloat(string(bits), 64)
			elem.SetFloat(v)
			return e

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v, e := strconv.ParseInt(string(bits), 10, 64)
			elem.SetInt(v)
			return e

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v, e := strconv.ParseUint(string(bits), 10, 64)
			elem.SetUint(v)
			return e

		case reflect.String:
			elem.SetString(string(bits))
			return nil

		default:
			buffer := bytes.NewBuffer(bits)
			return gob.NewDecoder(buffer).Decode(ptr)
		}
	}
	return fmt.Errorf("Failed to deserialize the value into non-pointer target")
}
