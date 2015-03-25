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

import "reflect"

func DetectType(v interface{}) string {
	switch v.(type) {
	case string:
		return "text/plain"

	default:
		var (
			elem  = reflect.ValueOf(v).Elem()
			stype = elem.Type()
		)
		for index := 0; index < elem.NumField(); index++ {
			if tag := stype.Field(index).Tag.Get("json"); tag != "" {
				return "application/json"
			}

			if tag := stype.Field(index).Tag.Get("xml"); tag != "" {
				return "application/xml"
			}
		}
	}

	return "text/plain"
}
