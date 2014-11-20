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

package web

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSerialization(t *testing.T) {
	Convey("[utils#Serialization]", t, func() {
		var input int = 1234567890
		var output int

		v, _ := Serialize(input)
		Deserialize(v, &output)

		So(output, ShouldEqual, input)
	})
}
