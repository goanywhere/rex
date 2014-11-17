/**
 *  ------------------------------------------------------------
 *  @project
 *  @file       system_test.go
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
package web

import (
	"bufio"
	"os"
	"testing"
)

func before(handler func(f string)) {
	filename := "/tmp/tmp.py"
	if file, err := os.Create(filename); err == nil {
		defer file.Close()
		defer os.Remove(filename)
		buffer := bufio.NewWriter(file)
		buffer.WriteString("#!/usr/bin/env python")
		buffer.Flush()

		handler(filename)
	}
}

func TestExists(t *testing.T) {
	exists := Exists("/tmp")
	if !exists {
		t.Errorf("Expected: true, Got: %v", exists)
	}

	exists = Exists("/NotExists")
	if exists {
		t.Errorf("Expected: false, Got: %v", exists)
	}
}

func TestIsDir(t *testing.T) {
	before(func(filename string) {
		flag := IsDir(filename)
		if flag {
			t.Errorf("Expected: false, Got: %v", flag)
		}
	})

	flag := IsDir("/tmp")
	if !flag {
		t.Errorf("Expected: true, Got: %v", flag)
	}
}

func TestIsFile(t *testing.T) {
	before(func(filename string) {
		flag := IsFile(filename)
		if !flag {
			t.Errorf("Expected: true, Got: %v", flag)
		}
	})

	flag := IsFile("/tmp")
	if flag {
		t.Errorf("Expected: false, Got: %v", flag)
	}
}
