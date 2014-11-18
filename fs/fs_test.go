/**
 *  ------------------------------------------------------------
 *  @project	web.go
 *  @file       fs_test.go
 *  @date       2014-11-17
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
package fs

import (
	"bufio"
	"os"
	"os/exec"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func setup(handler func(f string)) {
	filename := "/tmp/tmpfile"
	if file, err := os.Create(filename); err == nil {
		defer file.Close()
		defer os.Remove(filename)
		buffer := bufio.NewWriter(file)
		buffer.WriteString("I'm just a temp. file")
		buffer.Flush()

		handler(filename)
	}
}

func TestAbs(t *testing.T) {
	Convey("Absolute path check", t, func() {
		So(Abs("/tmp"), ShouldEqual, "/tmp")
	})
}

func TestCopy(t *testing.T) {
	Convey("Copy files/directories recursively", t, func() {
		filename := "GoAnywhereFake"
		exec.Command("touch", Abs("~/"+filename)).Run()
		defer os.Remove("/tmp/" + filename)
		err := Copy("~/"+filename, "/tmp")
		So(Exists("/tmp/"+filename), ShouldBeTrue)
		So(err, ShouldBeNil)

		exec.Command("mkdir", Abs("~/GoAnywhere")).Run()
		exec.Command("touch", Abs("~/GoAnywhere/Fake")).Run()
		defer os.RemoveAll("~/GoAnywhere")
		err = Copy("~/GoAnywhere", "/tmp")
		So(Exists("/tmp/GoAnywhere"), ShouldBeTrue)
		So(err, ShouldBeNil)
	})
}

func TestExists(t *testing.T) {
	Convey("Checks if the given path exists", t, func() {
		exists := Exists("/tmp")
		So(exists, ShouldBeTrue)

		exists = Exists("/NotExists")
		So(exists, ShouldBeFalse)
	})
}

func TestIsDir(t *testing.T) {
	setup(func(filename string) {
		flag := IsDir(filename)
		Convey("Checks if the given path is a directory", t, func() {
			So(flag, ShouldBeFalse)
		})
	})

	flag := IsDir("/tmp")
	Convey("Checks if the given path is a directory", t, func() {
		So(flag, ShouldBeTrue)
	})
}

func TestIsFile(t *testing.T) {
	setup(func(filename string) {
		flag := IsFile(filename)
		Convey("Checks if the given path is a file", t, func() {
			So(flag, ShouldBeTrue)
		})
	})

	flag := IsFile("/tmp")
	Convey("Checks if the given path is a file", t, func() {
		So(flag, ShouldBeFalse)
	})
}
