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
package auth

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEncrypt(t *testing.T) {
	src := "I'mAPlainSecret"
	settings.Set("AUTH_SECRET_KEY", "secretkey@example.com")
	Convey("auth.Encrypt Test", t, func() {
		secret := Encrypt(src)
		So(len(secret), ShouldEqual, 60)
		So(secret, ShouldNotEqual, Encrypt(src))
		So(secret, ShouldNotEqual, Encrypt(src))
		So(secret, ShouldNotEqual, Encrypt(src))
		So(secret, ShouldNotEqual, Encrypt(src))
		So(secret, ShouldNotEqual, Encrypt(src))
	})
}

func TestVerify(t *testing.T) {
	src := "I'mAPlainSecret"
	//key := "secretkey@example.com"
	Convey("auth.Verify Test", t, func() {
		for index := 0; index < 10; index++ {
			secret := Encrypt(src)
			So(Verify(src, secret), ShouldBeTrue)
		}
	})
}

func BenchmarkHash(b *testing.B) {
	src := "I'mAPlainSecret"
	settings.Set("AUTH_SECRET_KEY", "secretkey@example.com")
	for index := 0; index < b.N; index++ {
		hash(src)
	}
}

func BenchmarkEncrypt(b *testing.B) {
	src := "I'mAPlainSecret"
	settings.Set("AUTH_SECRET_KEY", "secretkey@example.com")
	for index := 0; index < b.N; index++ {
		Encrypt(src)
	}
}

func BenchmarkVerify(b *testing.B) {
	src := "I'mAPlainSecret"
	settings.Set("AUTH_SECRET_KEY", "secretkey@example.com")
	secret := Encrypt(src)
	for index := 0; index < b.N; index++ {
		Verify(src, secret)
	}
}
