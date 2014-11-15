/**
 *  ------------------------------------------------------------
 *  @project
 *  @file       env_test.go
 *  @date       2014-11-14
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
package env

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/goanywhere/web/crypto"
)

type Spec struct {
	App     string
	Debug   bool
	Total   int
	Version float32
	Tag     string `web:"multiple_words_tag"`
}

func setup() {
	os.Clearenv()
	Set("app", "example")
	Set("debug", "true")
	Set("total", "100")
	Set("version", "32.1")
	Set("multiple_words_tag", "ALT")
}

func TestFindKeyValue(t *testing.T) {
	k, v := findKeyValue(" test: value")
	if k != "test" || v != "value" {
		t.Errorf("Expect: <test: value>', Got: <%s: %s>", k, v)
	}

	k, v = findKeyValue(" test: value")
	if k != "test" || v != "value" {
		t.Errorf("Expect: <test: value>', Got: <%s: %s>", k, v)
	}

	k, v = findKeyValue("\ttest:\tvalue\t\n")
	if k != "test" || v != "value" {
		t.Errorf("Expect: <test: value>', Got: <%s: %s>", k, v)
	}

}

func TestLoad(t *testing.T) {
	Set("root", "/tmp")

	var pool = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*(-_+)")
	if env, err := os.Create("/tmp/.env"); err == nil {
		defer env.Close()
		secret := crypto.RandomString(64, pool)
		buffer := bufio.NewWriter(env)
		buffer.WriteString(fmt.Sprintf("secret=%s\n", secret))
		buffer.WriteString("app=myapp\n")
		buffer.Flush()

		Load()

		value := Get("secret")
		if value != secret {
			t.Errorf("Expected: %s, Got: %s", secret, value)
		}
	}
	os.Remove(".env")
}

func TestLoadInto(t *testing.T) {
	var spec Spec
	setup()
	if err := LoadInto(&spec); err != nil {
		t.Error(err.Error())
	}
}

func TestGetString(t *testing.T) {
	var spec Spec

	setup()

	if err := LoadInto(&spec); err != nil {
		t.Error(err.Error())
	}
	if spec.App != "example" {
		t.Errorf("Expect: 'example', Got: %s", spec.App)
	}

	Set("secret", "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*(-_=+)")
	if Get("secret") != "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*(-_=+)" {
		t.Errorf("Expect: '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*(-_=+)', Got: %s", Get("secret"))
	}
}

func TestGetBool(t *testing.T) {
	var spec Spec
	setup()
	if err := LoadInto(&spec); err != nil {
		t.Error(err.Error())
	}

	if !spec.Debug {
		t.Errorf("Expect: true, Got: %v", spec.Debug)
	}

	value, err := GetBool("NotFound")
	if value || err != nil {
		t.Errorf("Expect: false, Got: %v", value)
	}
}

func TestGetInt(t *testing.T) {
	var spec Spec
	setup()
	if err := LoadInto(&spec); err != nil {
		t.Error(err.Error())
	}

	if spec.Total != 100 {
		t.Errorf("Expect: 100, Got: %d", spec.Total)
	}
}

func TestGetFloat(t *testing.T) {
	var spec Spec
	setup()
	if err := LoadInto(&spec); err != nil {
		t.Error(err.Error())
	}

	if spec.Version != 32.1 {
		t.Errorf("Expect: 32.1, Got: %f", spec.Version)
	}
}

func TestTag(t *testing.T) {
	var spec Spec
	setup()
	if err := LoadInto(&spec); err != nil {
		t.Error(err.Error())
	}
	if spec.Tag != "ALT" {
		t.Errorf("Expect: 'MULTIPLE_WORDS_TAG', Got: %s", spec.Tag)
	}
}

func TestAccess(t *testing.T) {
	Set("shell", "/bin/zsh")
	if Get("shell") != "/bin/zsh" {
		t.Errorf("Expect: /bin/zsh, Got: %s", Get("shell"))
	}

	Set("Anything", "content")
	if Get("anything") != "content" {
		t.Errorf("Expect: 'content', Got: %s", Get("anything"))
	}

	if Get("NotFound") != "" {
		t.Errorf("Expect empty string, Got: %s", Get("NotFound"))
	}
}

func TestValues(t *testing.T) {
	os.Clearenv()
	values := Values()
	if len(values) != 0 {
		t.Errorf("Expect: 0, Got: %d", len(values))
	}
	Set("app", "me")

	values = Values()
	if len(values) != 1 {
		t.Errorf("Expect: 1, Got: %d", len(values))
	}
	if values[Prefix+"_APP"] != "me" {
		t.Errorf("Expect: 'me', Got: '%s'", values[Prefix+"_APP"])
	}
}
