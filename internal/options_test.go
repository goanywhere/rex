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
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDefine(t *testing.T) {
	var options = Options()
	Convey("Define primitive values via os environment", t, func() {
		Define("name", "rex")
		v, exists := options.Get("name")
		So(v, ShouldEqual, "rex")
		So(exists, ShouldBeTrue)

		v, exists = options.Get("NotFound")
		So(v, ShouldEqual, "")
		So(exists, ShouldBeFalse)
	})
}

func TestOption(t *testing.T) {
	Convey("Load defined option into appointed address", t, func() {
		var bv bool
		Option("bool", &bv)
		So(bv, ShouldBeFalse)
		Define("bool", true)
		Option("bool", &bv)
		So(bv, ShouldBeTrue)

		var fv32 float32
		Option("fv32", &fv32)
		So(fv32, ShouldEqual, 0.0)

		Define("fv32", 123.45)
		Option("fv32", &fv32)
		So(fv32, ShouldEqual, 123.45)

		var fv64 float64
		Option("fv64", &fv64)
		So(fv64, ShouldEqual, 0.0)

		Define("fv64", 123456789987654321.123456789)
		Option("fv64", &fv64)
		So(fv64, ShouldEqual, 123456789987654321.123456789)

		var number int
		Option("number", &number)
		So(number, ShouldEqual, 0)

		Define("number", 123)
		Option("number", &number)
		So(number, ShouldEqual, 123)

		var n64 int64
		Option("n64", &n64)
		So(n64, ShouldEqual, 0)

		Define("n64", 9223372036854775807)
		Option("n64", &n64)
		So(n64, ShouldEqual, 9223372036854775807)

	})
}
