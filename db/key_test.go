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
