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
package internal

import (
	"sync"

	"github.com/goanywhere/x/env"
)

var (
	once     sync.Once
	settings *env.Env
)

func Settings() *env.Env {
	once.Do(func() {
		settings = env.New("rex")
		settings.Set("debug", true)
		settings.Set("mode", "debug")
		settings.Set("dir.templates", "templates")
		// default environmental headers for modules.Env
		settings.Set("header.X-UA-Compatible", "deny")
		settings.Set("header.X-Frame-settings", "nosniff")
		settings.Set("header.X-XSS-Protection", "1; mode=block")
		settings.Set("header.X-Content-Type-settings", "IE=Edge,chrome=1")
		// session cookie defaults
		settings.Set("session.cookie.name", "session")
		settings.Set("session.cookie.path", "/")
		settings.Set("session.cookie.domain", "")
		settings.Set("session.cookie.secure", false)
		settings.Set("session.cookie.httponly", true)
		settings.Set("session.cookie.maxage", 3600*24*7)
	})
	return settings
}
