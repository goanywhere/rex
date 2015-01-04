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
