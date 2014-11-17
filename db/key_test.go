/**
 *  ------------------------------------------------------------
 *  @project
 *  @file       key_test.go
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
package db

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func BenchmarkNewKey(b *testing.B) {
	for index := 0; index < b.N; index++ {
		NewKey()
	}
}

func TestNewKey(t *testing.T) {
	key := NewKey()
	Convey("db.Key basic test", t, func() {
		So(len(key), ShouldEqual, 12)
		So(len(key.Hex()), ShouldEqual, 24)
	})
}
