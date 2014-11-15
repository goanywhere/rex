/**
 *  ------------------------------------------------------------
 *  @project
 *  @file       case_test.go
 *  @date       2014-11-15
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
	"testing"
)

func TestToBool(t *testing.T) {
	val, err := ToBool("true")
	if !val || err != nil {
		t.Errorf("Expect: true, Got: %v", val)
	}

	val, err = ToBool("hello")
	if val || err == nil {
		t.Errorf("Expect: false, Got: %v Err: %v", val, err)
	}

	val, err = ToBool("TrUe")
	if !val || err != nil {
		t.Errorf("Expect: true, Got: %v Err: %v", val, err)
	}

	val, err = ToBool(1)
	if !val || err != nil {
		t.Errorf("Expect: true, Got: %v", val)
	}

	val, err = ToBool(false)
	if val || err != nil {
		t.Errorf("Expect: false, Got: %v", val)
	}

	val, err = ToBool(0)
	if val || err != nil {
		t.Errorf("Expect: false, Got: %v", val)
	}

	val, err = ToBool("TrUE")
	if !val || err != nil {
		t.Errorf("Expect: true, Got: %v Err: %v", val, err)
	}

	val, err = ToBool("")
	if val || err != nil {
		t.Errorf("Expect: false, Got: %v Err: %v", val, err)
	}
}

func TestToFloat(t *testing.T) {
	val, err := ToFloat("123.45")
	if val != 123.45 || err != nil {
		t.Errorf("Expect 123.45, Got: %v, Err: %v", val, err)
	}

	val, err = ToFloat(12345)
	if val != 12345.0 || err != nil {
		t.Errorf("Expect 12345.0, Got: %v, Err: %v", val, err)
	}

	val, err = ToFloat(12345.000)
	if val != 12345.0 || err != nil {
		t.Errorf("Expect 12345.0, Got: %v, Err: %v", val, err)
	}

}

func TestToInt(t *testing.T) {
	val, err := ToInt("23")
	if val != 23 || err != nil {
		t.Errorf("Expect: 23, Got: %v, Err: %v", val, err)
	}

	val, err = ToInt(23)
	if val != 23 || err != nil {
		t.Errorf("Expect: 23, Got: %v, Err: %v", val, err)
	}

	val, err = ToInt(23.01)
	if val != 23 || err != nil {
		t.Errorf("Expect: 23, Got: %v, Err: %v", val, err)
	}

	val, err = ToInt(nil)
	if val != 0 || err != nil {
		t.Errorf("Expect: 0, Got: %v, Err: %v", val, err)
	}
}

func TestToString(t *testing.T) {
	val, err := ToString(123)
	if val != "123" || err != nil {
		t.Errorf("Expect '123', Got: %v, Err: %v", val, err)
	}

	val, err = ToString(123.34)
	if val != "123.34" || err != nil {
		t.Errorf("Expect '123.34', Got: %v, Err: %v", val, err)
	}

	val, err = ToString([]byte{96})
	if val != "`" || err != nil {
		t.Errorf("Expect '`', Got: %v, Err: %v", val, err)
	}

	val, err = ToString(nil)
	if val != "" || err != nil {
		t.Errorf("Expect empty string, Got: %v, Err: %v", val, err)
	}
}
