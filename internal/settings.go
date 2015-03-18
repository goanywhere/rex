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

import "sync"

type session struct {
	Name     string `env:"SESSION_COOKIE_NAME"`
	Path     string `env:"SESSION_COOKIE_PATH"`
	Domain   string `env:"SESSION_COOKIE_DOMAIN"`
	Secure   bool   `env:"SESSION_COOKIE_SECURE"`
	HttpOnly bool   `env:"SESSION_COOKIE_HTTPONLY"`
	MaxAge   int    `env:"SESSION_COOKIE_MAXAGE"`
}

type settings struct {
	Session *session

	Root       string
	Debug      bool     `env:"DEBUG"`
	Port       int      `env:"PORT"`
	View       string   `env:"VIEW"`
	SecretKey  string   `env:"SECRET_KEY"`
	SecretKeys []string `env:"SECRET_KEYS"`
}

var (
	once   sync.Once
	config *settings
)

func Settings() *settings {
	once.Do(func() {
		config = new(settings)
		// session cookie defaults
		config.Debug = true
		config.Port = 5000
		config.View = "views"

		config.Session = new(session)
		config.Session.Name = "session"
		config.Session.Path = "/"
		config.Session.Domain = ""
		config.Session.Secure = false
		config.Session.HttpOnly = true
		config.Session.MaxAge = 3600 * 24 * 7
	})
	return config
}
