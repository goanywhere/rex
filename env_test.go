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
package web

import (
	"os"
	"testing"
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
	Env.Set("app", "example")
	Env.Set("debug", "true")
	Env.Set("total", "100")
	Env.Set("version", "32.1")
	Env.Set("multiple_words_tag", "ALT")
}

func TestLoad(t *testing.T) {
	var spec Spec
	setup()
	if err := Env.Load(&spec); err != nil {
		t.Error(err.Error())
	}
}

func TestGetString(t *testing.T) {
	var spec Spec

	setup()

	if err := Env.Load(&spec); err != nil {
		t.Error(err.Error())
	}
	if spec.App != "example" {
		t.Errorf("Expect: 'example', Got: %s", spec.App)
	}
}

func TestGetBool(t *testing.T) {
	var spec Spec
	setup()
	if err := Env.Load(&spec); err != nil {
		t.Error(err.Error())
	}

	if !spec.Debug {
		t.Errorf("Expect: true, Got: %v", spec.Debug)
	}
}

func TestGetInt(t *testing.T) {
	var spec Spec
	setup()
	if err := Env.Load(&spec); err != nil {
		t.Error(err.Error())
	}

	if spec.Total != 100 {
		t.Errorf("Expect: 100, Got: %d", spec.Total)
	}
}

func TestGetFloat(t *testing.T) {
	var spec Spec
	setup()
	if err := Env.Load(&spec); err != nil {
		t.Error(err.Error())
	}

	if spec.Version != 32.1 {
		t.Errorf("Expect: 32.1, Got: %f", spec.Version)
	}
}

func TestTag(t *testing.T) {
	var spec Spec
	setup()
	if err := Env.Load(&spec); err != nil {
		t.Error(err.Error())
	}
	if spec.Tag != "ALT" {
		t.Errorf("Expect: 'MULTIPLE_WORDS_TAG', Got: %s", spec.Tag)
	}
}

func TestAccess(t *testing.T) {
	Env.Set("shell", "/bin/zsh")
	if Env.Get("shell") != "/bin/zsh" {
		t.Errorf("Expect: /bin/zsh, Got: %s", Env.Get("shell"))
	}

	Env.Set("Anything", "content")
	if Env.Get("anything") != "content" {
		t.Errorf("Expect: 'content', Got: %s", Env.Get("anything"))
	}
}

func TestValues(t *testing.T) {
	os.Clearenv()
	values := Env.Values()
	if len(values) != 0 {
		t.Errorf("Expect: 0, Got: %d", len(values))
	}
	Env.Set("app", "me")

	values = Env.Values()
	if len(values) != 1 {
		t.Errorf("Expect: 1, Got: %d", len(values))
	}
	if values[Prefix+"_APP"] != "me" {
		t.Errorf("Expect: 'me', Got: '%s'", values[Prefix+"_APP"])
	}
}
