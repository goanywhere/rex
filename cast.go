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
