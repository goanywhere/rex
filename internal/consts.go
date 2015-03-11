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

// Settings Keys
const (
	DEBUG = "DEBUG"

	HTTP_HEADER_X_UA_Compatible           = "HTTP_HEADER_X_UA_Compatible"
	HTTP_HEADER_X_Frame_Settings          = "HTTP_HEADER_X_Frame_Settings"
	HTTP_HEADER_X_XSS_Protection          = "HTTP_HEADER_X_XSS_Protection"
	HTTP_HEADER_X_Content_Type_Options    = "HTTP_HEADER_X_Content_Type_Options"
	HTTP_HEADER_Strict_Transport_Security = "HTTP_HEADER_Strict_Transport_Security"

	SESSION_COOKIE_NAME     = "SESSION_COOKIE_NAME"
	SESSION_COOKIE_PATH     = "SESSION_COOKIE_PATH"
	SESSION_COOKIE_DOMAIN   = "SESSION_COOKIE_DOMAIN"
	SESSION_COOKIE_SECURE   = "SESSION_COOKIE_SECURE"
	SESSION_COOKIE_HTTPONLY = "SESSION_COOKIE_HTTPONLY"
	SESSION_COOKIE_MAXAGE   = "SESSION_COOKIE_MAXAGE"
)

var ContentType = struct {
	Name string
	HTML string
	JSON string
	XML  string
	Text string
}{
	Name: "Content-Type",
	HTML: "text/html; charset=UTF-8",
	JSON: "application/json; charset=UTF-8",
	XML:  "application/xml; charset=UTF-8",
	Text: "text/plain; charset=UTF-8",
}
