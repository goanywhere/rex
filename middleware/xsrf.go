/**
 *  ------------------------------------------------------------
 *  @project	webapp
 *  @file       xsrf.go
 *  @date       2014-10-21
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
package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/goanywhere/web"
)

const (
	xsrfCookieName = "xsrf"
	xsrfHeaderName = "X-XSRF-Token"
	xsrfFieldName  = "xsrftoken"

	xsrfMaxAge  = 3600 * 24 * 365
	xsrfTimeout = time.Hour * 24 * 365
)

var (
	errInvalidReferer = "Referer URL is missing from the request or the value was malformed."
	errInvalidToken   = "Invalid/Mismatch XSRF tokens."

	xsrfPattern   = regexp.MustCompile("[^0-9a-zA-Z-_]")
	unsafeMethods = regexp.MustCompile("^(DELETE|POST|PUT)$")
)

type xsrf struct {
	context *web.Context
}

// See http://en.wikipedia.org/wiki/Same-origin_policy
func (self *xsrf) checkOrigin() bool {
	if self.context.Request.URL.Scheme == "https" {
		// See [OWASP]; Checking the Referer Header.
		referer, err := url.Parse(self.context.Request.Header.Get("Referer"))

		if err != nil || referer.String() == "" ||
			referer.Scheme != self.context.Request.URL.Scheme ||
			referer.Host != self.context.Request.URL.Host {

			return false
		}
	}
	return true
}

func (self *xsrf) checkToken(token string) bool {
	// Header always takes precedance of form field since some popular
	// JavaScript frameworks allow global custom headers for all AJAX requests.
	query := self.context.Request.Header.Get(xsrfFieldName)
	if query == "" {
		query = self.context.Request.FormValue(xsrfFieldName)
	}

	// 1) basic length comparison.
	if query == "" || len(query) != len(token) {
		return false
	}
	// *sanitize* incoming masked token.
	query = xsrfPattern.ReplaceAllString(query, "")

	// 2) byte-based comparison.
	a, _ := base64.URLEncoding.DecodeString(token)
	b, _ := base64.URLEncoding.DecodeString(query)
	if subtle.ConstantTimeCompare(a, b) != 1 {
		return false
	}

	// 3) issued time checking.
	index := bytes.LastIndex(b, []byte{'|'})
	if index != 40 {
		return false
	}

	nanos, err := strconv.ParseInt(string(b[index+1:]), 10, 64)
	if err != nil {
		return false
	}
	now := time.Now()
	issueTime := time.Unix(0, nanos)

	if now.Sub(issueTime) >= xsrfTimeout {
		return false
	}

	// Ensure the token is not from the *future*, allow 1 minute grace period.
	if issueTime.After(now.Add(1 * time.Minute)) {
		return false
	}

	return true
}

func (self *xsrf) generate() string {
	nano := time.Now().UnixNano()
	hash := hmac.New(sha1.New, []byte(web.RandomString(32, nil)))
	fmt.Fprintf(hash, "%s|%d", web.RandomString(12, nil), nano)
	token := fmt.Sprintf("%s|%d", hex.EncodeToString(hash.Sum(nil)), nano)
	return base64.URLEncoding.EncodeToString([]byte(token))
}

func (self *xsrf) token() string {
	var secure bool = false
	if self.context.Request.URL.Scheme == "https" {
		secure = true
	}
	// Ensure we have XSRF token in the cookie first.
	token := self.context.Cookie(xsrfCookieName)
	if token == "" {
		token = self.generate()
		// The max-age directive takes priority over Expires.
		//	http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
		self.context.SetCookie(&http.Cookie{
			Name:     xsrfCookieName,
			Value:    token,
			MaxAge:   xsrfMaxAge,
			Path:     "/",
			HttpOnly: true,
			Secure:   secure,
		})
	}
	self.context.Header().Set("X-XSRF-Token", token)
	self.context.Set("xsrftoken", token)
	return token
}

// ---------------------------------------------------------------------------
//  XSRF Middleware Supports
// ---------------------------------------------------------------------------
func XSRF(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := web.NewContext(w, r)
		x := &xsrf{ctx}
		token := x.token()

		if unsafeMethods.MatchString(r.Method) {
			// Ensure the URL came for "Referer" under HTTPS.
			if !x.checkOrigin() {
				ctx.Forbidden(errInvalidReferer)
			}

			// length => bytes => issue time checkpoints.
			if !x.checkToken(token) {
				ctx.Forbidden(errInvalidToken)
			}
		}

		// ensure browser will invalidate the cached XSRF token.
		w.Header().Add("Vary", "Cookie")

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
