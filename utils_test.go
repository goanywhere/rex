/**
 *  ------------------------------------------------------------
 *  @project	web.go
 *  @file       utils_test.go
 *  @date       2014-11-18
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
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSerialization(t *testing.T) {
	Convey("[utils] serialize/deserialize", t, func() {
		var input int = 1234567890
		var output int

		v, _ := Serialize(input)
		Deserialize(v, &output)

		So(output, ShouldEqual, input)
	})
}
