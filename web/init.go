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
package web

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/goanywhere/rex/config"
	"github.com/goanywhere/rex/crypto"
	"github.com/goanywhere/rex/template"
)

const (
	xSessionKey = "sessionid"
)

var (
	settings = config.Settings()

	process   string
	loader    = template.NewLoader(settings.Templates)
	signature = crypto.NewSignature(settings.Secret)
)

func init() {
	// prepare a md5-based fixed length (32-bits) string for Go process.
	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}
	// system pid combined with timestamp to identity current go process.
	pid := fmt.Sprintf("%d:%d", os.Getpid(), time.Now().UnixNano())
	hash := hmac.New(md5.New, []byte(fmt.Sprintf("%s-%s", hostname, pid)))
	process = hex.EncodeToString(hash.Sum(nil))
}
