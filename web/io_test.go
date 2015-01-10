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
package web

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAcceptEncodings(t *testing.T) {
	var r = &http.Request{Method: "GET"}
	var encodings []string
	Convey("Accept-Encoding", t, func() {
		r.Header = http.Header{"Accept-Encoding": {"gzip,deflate,sdch"}}
		encodings = AcceptEncodings(r)
		So(encodings[0], ShouldEqual, "gzip")
		So(encodings[1], ShouldEqual, "deflate")
		So(encodings[2], ShouldEqual, "sdch")

		r.Header = http.Header{"Accept-Encoding": {"gzip,deflate;q=0.5,sdch;q=0.2,*;q=0.0"}}
		encodings = AcceptEncodings(r)
		So(len(encodings), ShouldEqual, 3)
		So(encodings[0], ShouldEqual, "gzip")
		So(encodings[1], ShouldEqual, "deflate")
		So(encodings[2], ShouldEqual, "sdch")

		r.Header = http.Header{"Accept-Encoding": {"deflate;q=0.5,gzip;q=1.0"}}
		encodings = AcceptEncodings(r)
		So(encodings[0], ShouldEqual, "gzip")
		So(encodings[1], ShouldEqual, "deflate")

		r.Header = http.Header{"Accept-Encoding": {"deflate,gzip;q=0.5"}}
		encodings = AcceptEncodings(r)
		So(encodings[0], ShouldEqual, "deflate")
		So(encodings[1], ShouldEqual, "gzip")

		r.Header = http.Header{"Accept-Encoding": {"compress,identify,gzip;q=1.0"}}
		encodings = AcceptEncodings(r)
		So(len(encodings), ShouldEqual, 3)
		So(encodings[0], ShouldEqual, "gzip")
		So(encodings[1], ShouldEqual, "compress")
		So(encodings[2], ShouldEqual, "identify")
	})
}
