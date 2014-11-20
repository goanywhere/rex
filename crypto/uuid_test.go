/**
 * ----------------------------------------------------------------------
 *       ______      ___                         __
 *      / ____/___  /   |  ____  __  ___      __/ /_  ___  ________
 *     / / __/ __ \/ /| | / __ \/ / / / | /| / / __ \/ _ \/ ___/ _ \
 *    / /_/ / /_/ / ___ |/ / / / /_/ /| |/ |/ / / / /  __/ /  /  __/
 *    \____/\____/_/  |_/_/ /_/\__. / |__/|__/_/ /_/\___/_/   \___/
 *                            /____/
 * ----------------------------------------------------------------------
 *  Copyright © 2014 GoAnywhere Ltd. All Rights Reserved.
 * ----------------------------------------------------------------------*/

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
