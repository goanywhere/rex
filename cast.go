/**
 *  ------------------------------------------------------------
 *  @project	web.go
 *  @file       cast.go
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
	"fmt"
	"strconv"
	"strings"
)

// ToBool converts the given value to bool.
// Supported Types:
//	bool, int, float32, float64, string
func ToBool(raw interface{}) (bool, error) {
	switch t := raw.(type) {
	case bool:
		return t, nil
	case nil:
		return false, nil
	case int:
		if raw.(int) == 0 {
			return false, nil
		}
		return true, nil
	case float32, float64:
		if raw.(float32) == 0.0 || raw.(float64) == 0.0 {
			return false, nil
		}
		return true, nil
	case string:
		return strconv.ParseBool(strings.ToLower(raw.(string)))
	}
	return false, fmt.Errorf("Unable to convert %v to bool", raw)
}

// ToFloat converts string/int to their float64 value.
func ToFloat(raw interface{}) (float64, error) {
	switch val := raw.(type) {
	case float32:
		return float64(val), nil
	case float64:
		return val, nil
	case int, int8, int16, int32, int64:
		return float64(val.(int)), nil
	case string:
		v, err := strconv.ParseFloat(val, 64)
		if err == nil {
			return float64(v), nil
		}
		return 0.0, fmt.Errorf("Unable to convert %v to float", val)
	default:
		return 0.0, fmt.Errorf("Unable to convert %v to float", val)
	}
}

// ToInt converts int/float/bool/nil to int value.
func ToInt(raw interface{}) (int, error) {
	switch val := raw.(type) {
	case bool:
		if bool(val) {
			return 1, nil
		}
		return 0, nil
	case int, int8, int16, int32, int64:
		return raw.(int), nil
	case float64:
		return int(val), nil
	case string:
		v, err := strconv.ParseInt(val, 0, 0)
		if err == nil {
			return int(v), nil
		}
		return 0, fmt.Errorf("Unable to convert %v to int", raw)
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("Unable to convert %v to int", raw)
	}
}

// ToString converts float/int/byte/nil into string value.
func ToString(raw interface{}) (string, error) {
	switch val := raw.(type) {
	case string:
		return val, nil
	case []byte:
		return string(val), nil
	case int:
		return strconv.FormatInt(int64(raw.(int)), 10), nil
	case float64:
		return strconv.FormatFloat(raw.(float64), 'f', -1, 64), nil
	case nil:
		return "", nil
	default:
		return "", fmt.Errorf("Unable to convert %v to string", raw)
	}
}
