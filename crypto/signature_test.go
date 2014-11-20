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
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSignature(t *testing.T) {
	s := NewSignature(RandomString(128, nil))

	k, v := "name", []byte("Hello Signature")
	Convey("[crypto#Signature]", t, func() {
		value, err := s.Encode(k, v)
		So(value, ShouldNotBeNil)
		So(err, ShouldBeNil)

		src, err := s.Decode(k, value)
		So(bytes.Compare(v, src), ShouldEqual, 0)
		So(err, ShouldBeNil)
	})
}
