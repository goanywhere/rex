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
package template

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestExtends(t *testing.T) {
	Convey("Regex. Tag:Extends", t, func() {
		So(extends.MatchString("{% extends \"layouts/base.html\" %}"), ShouldBeTrue)
		So(extends.MatchString("{% extends \"c://Users/someone/templates/layouts/base.html\" %}"), ShouldBeTrue)
		So(len(extends.FindStringSubmatch("{% extends \"layouts/base.html\" %}")), ShouldEqual, 2)
		So(len(extends.FindStringSubmatch("{% extends \"c://Users/someone/templates/layouts/base.html\" %}")), ShouldEqual, 2)
	})
}

func TestInclude(t *testing.T) {
	Convey("Regex. Tag:Include", t, func() {
		So(include.MatchString("{% include \"partials/header.html\" %}"), ShouldBeTrue)
		So(include.MatchString("{% include \"c://Users/someone/templates/partials/nav.html\" %}"), ShouldBeTrue)
		So(len(include.FindStringSubmatch("{% include \"partials/nav.html\" %}")), ShouldEqual, 2)
		So(len(include.FindStringSubmatch("{% include \"c://Users/someone/templates/partials/nav.html\" %}")), ShouldEqual, 2)

		matches := include.FindAllStringSubmatch(`{% include "partials/header.html" %}\t\n{% include "partials/footer.html" %}\n`, -1)
		So(len(matches), ShouldEqual, 2)
		So(len(matches[0]), ShouldEqual, 2)
		So(matches[0][1], ShouldEqual, "partials/header.html")
		So(matches[1][1], ShouldEqual, "partials/footer.html")
	})
}
