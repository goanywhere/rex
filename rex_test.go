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
package rex

/*
import (
	"os"
	"path"
	"testing"

	"github.com/flosch/pongo2"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLoadViews(t *testing.T) {
	Convey("rex.loadViews", t, func() {
		temp := os.TempDir()
		html := path.Join(temp, "index.html")
		xml := path.Join(temp, "index.xml")

		os.Create(html)
		os.Create(xml)

		include := path.Join(temp, "include")
		header := path.Join(include, "header.html")
		os.Mkdir(include, os.ModePerm)
		os.Create(header)

		layouts := path.Join(temp, "layouts")
		base := path.Join(layouts, "base.html")
		os.Mkdir(layouts, os.ModePerm)
		os.Create(base)

		loadViews(temp)

		So(len(views), ShouldEqual, 2)

		os.RemoveAll(include)
		os.RemoveAll(layouts)
		os.Remove(html)
		os.Remove(xml)
		views = make(map[string]*pongo2.Template)
	})
}
*/
