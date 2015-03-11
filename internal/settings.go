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
		settings = env.New()
		settings.Set(DEBUG, true)
		settings.Set("TEMPLATES", "views")
		// default environmental headers for modules.Env
		settings.Set(HTTP_HEADER_X_UA_Compatible, "deny")
		settings.Set(HTTP_HEADER_X_Frame_Settings, "nosniff")
		settings.Set(HTTP_HEADER_X_XSS_Protection, "1; mode=block")
		settings.Set(HTTP_HEADER_X_Content_Type_Options, "IE=Edge,chrome=1")
		settings.Set(HTTP_HEADER_Strict_Transport_Security, "max-age=31536000; includeSubdomains; preload")
		// session cookie defaults
		settings.Set(SESSION_COOKIE_NAME, "session")
		settings.Set(SESSION_COOKIE_PATH, "/")
		settings.Set(SESSION_COOKIE_DOMAIN, "")
		settings.Set(SESSION_COOKIE_SECURE, false)
		settings.Set(SESSION_COOKIE_HTTPONLY, true)
		settings.Set(SESSION_COOKIE_MAXAGE, 3600*24*7)
	})
	return settings
}
