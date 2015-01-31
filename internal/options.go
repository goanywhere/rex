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
	once    sync.Once
	options *env.Env
)

func Options() *env.Env {
	once.Do(func() {
		options = env.New("rex")
		options.Set("debug", true)
		options.Set("mode", "debug")
		options.Set("dir.templates", "templates")
		// default environmental headers for modules.Env
		options.Set("header.X-UA-Compatible", "deny")
		options.Set("header.X-Frame-Options", "nosniff")
		options.Set("header.X-XSS-Protection", "1; mode=block")
		options.Set("header.X-Content-Type-options", "IE=Edge,chrome=1")
		// session cookie defaults
		options.Set("context.cookie.path", "/")
		options.Set("context.cookie.name", "gsid")
		options.Set("context.cookie.secure", false)
		options.Set("context.cookie.httponly", true)
		options.Set("context.cookie.maxage", 3600*24*7)

		options.Set("session.cookie.name", "gsid")
		options.Set("session.cookie.maxage", 3600*24*7)
		options.Set("session.cookie.httponly", true)
		options.Set("session.cookie.path", "/")
		options.Set("session.cookie.secure", false)
	})
	return options
}
