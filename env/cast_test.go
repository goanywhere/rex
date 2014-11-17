/**
 *  ------------------------------------------------------------
 *  @project
 *  @file       case_test.go
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
package env

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestToBool(t *testing.T) {
	Convey("ToBool conversion test", t, func() {
		val, err := ToBool("true")
		So(val, ShouldBeTrue)
		So(err, ShouldBeNil)

		val, err = ToBool("hello")
		So(val, ShouldBeFalse)
		So(err, ShouldNotBeNil)

		val, err = ToBool("TrUe")
		So(val, ShouldBeTrue)
		So(err, ShouldBeNil)

		val, err = ToBool(1)
		So(val, ShouldBeTrue)
		So(err, ShouldBeNil)

		val, err = ToBool(0)
		So(val, ShouldBeFalse)
		So(err, ShouldBeNil)

		val, err = ToBool(true)
		So(val, ShouldBeTrue)
		So(err, ShouldBeNil)

		val, err = ToBool(false)
		So(val, ShouldBeFalse)
		So(err, ShouldBeNil)

		val, err = ToBool("")
		So(val, ShouldBeFalse)
		So(err, ShouldBeNil)
	})

}

func TestToFloat(t *testing.T) {
	Convey("ToFloat conversion test", t, func() {
		val, err := ToFloat("123.45")
		So(val, ShouldEqual, 123.45)
		So(err, ShouldBeNil)

		val, err = ToFloat(12345)
		So(val, ShouldEqual, 12345.0)
		So(err, ShouldBeNil)

		val, err = ToFloat(12345.000)
		So(val, ShouldEqual, 12345.0)
		So(err, ShouldBeNil)
	})
}

func TestToInt(t *testing.T) {
	Convey("ToInt conversion test", t, func() {
		val, err := ToInt("23")
		So(val, ShouldEqual, 23)
		So(err, ShouldBeNil)

		val, err = ToInt(23)
		So(val, ShouldEqual, 23)
		So(err, ShouldBeNil)

		val, err = ToInt(23.01)
		So(val, ShouldEqual, 23)
		So(err, ShouldBeNil)

		val, err = ToInt(nil)
		So(val, ShouldEqual, 0)
		So(err, ShouldBeNil)
	})
}

func TestToString(t *testing.T) {
	Convey("ToString conversion test", t, func() {
		val, err := ToString(123)
		So(val, ShouldEqual, "123")
		So(err, ShouldBeNil)

		val, err = ToString(123.34)
		So(val, ShouldEqual, "123.34")
		So(err, ShouldBeNil)

		val, err = ToString([]byte{96})
		So(val, ShouldEqual, "`")
		So(err, ShouldBeNil)

		val, err = ToString(nil)
		So(val, ShouldEqual, "")
		So(err, ShouldBeNil)
	})
}
