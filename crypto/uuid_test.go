/**
 *  ------------------------------------------------------------
 *  @project	web.go
 *  @file       uuid_test.go
 *  @date       2014-11-12
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
package crypto

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func BenchmarkNewV1(b *testing.B) {
	for index := 0; index < b.N; index++ {
		NewV1()
	}
}

func BenchmarkNewV3(b *testing.B) {
	for index := 0; index < b.N; index++ {
		NewV3(NamespaceOID, "jim.zhan@me.com")
	}
}

func BenchmarkNewV4(b *testing.B) {
	for index := 0; index < b.N; index++ {
		NewV4()
	}
}

func BenchmarkNewV5(b *testing.B) {
	for index := 0; index < b.N; index++ {
		NewV5(NamespaceOID, "jim.zhan@me.com")
	}
}

func TestNewV3(t *testing.T) {
	//""f2107fc9-aea6-3bf0-9ad8-3bef1b5f808b
	uuid := NewV3(NamespaceOID, "test@example.com")
	str := "70cd6896-ecb5-3388-85ca-384edc3f3e66"
	Convey("UUID Version 3 test", t, func() {
		So(uuid.Version(), ShouldEqual, 3)
		So(uuid.String(), ShouldEqual, str)
	})
}

func TestNewV5(t *testing.T) {
	//""f2107fc9-aea6-3bf0-9ad8-3bef1b5f808b
	uuid := NewV5(NamespaceOID, "test@example.com")
	str := "067f23a9-76a5-5585-b119-32402a120978"
	Convey("UUID Version 5 test", t, func() {
		So(uuid.Version(), ShouldEqual, 5)
		So(uuid.String(), ShouldEqual, str)
	})
}
