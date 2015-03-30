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
package rex

import "net/http"

// Shortcuts to Context.Send with standard
// HTTP status codes, defined in RFC 2616.

// HTTP 1XX
// ----------------------------------------
func (self *Context) Continue(v interface{}) {
	self.WriteHeader(http.StatusContinue)
	self.Send(v)
}

func (self *Context) SwitchingProtocols(v interface{}) {
	self.WriteHeader(http.StatusSwitchingProtocols)
	self.Send(v)
}

// HTTP 2XX
// ----------------------------------------
func (self *Context) OK(v interface{}) {
	self.WriteHeader(http.StatusOK)
	self.Send(v)
}

func (self *Context) Created(v interface{}) {
	self.WriteHeader(http.StatusCreated)
	self.Send(v)
}

func (self *Context) Accepted(v interface{}) {
	self.WriteHeader(http.StatusAccepted)
	self.Send(v)
}

func (self *Context) NonAuthoritativeInfo(v interface{}) {
	self.WriteHeader(http.StatusNonAuthoritativeInfo)
	self.Send(v)
}

func (self *Context) NoContent(v interface{}) {
	self.WriteHeader(http.StatusNoContent)
	self.Send(v)
}

func (self *Context) ResetContent(v interface{}) {
	self.WriteHeader(http.StatusResetContent)
	self.Send(v)
}

func (self *Context) PartialContent(v interface{}) {
	self.WriteHeader(http.StatusPartialContent)
	self.Send(v)
}

// HTTP 3XX
// ----------------------------------------
func (self *Context) MultipleChoices(v interface{}) {
	self.WriteHeader(http.StatusMultipleChoices)
	self.Send(v)
}

func (self *Context) MovedPermanently(v interface{}) {
	self.WriteHeader(http.StatusMovedPermanently)
	self.Send(v)
}

func (self *Context) Found(v interface{}) {
	self.WriteHeader(http.StatusFound)
	self.Send(v)
}

func (self *Context) SeeOther(v interface{}) {
	self.WriteHeader(http.StatusSeeOther)
	self.Send(v)
}

func (self *Context) NotModified(v interface{}) {
	self.WriteHeader(http.StatusNotModified)
	self.Send(v)
}

func (self *Context) UseProxy(v interface{}) {
	self.WriteHeader(http.StatusUseProxy)
	self.Send(v)
}

func (self *Context) TemporaryRedirect(v interface{}) {
	self.WriteHeader(http.StatusTemporaryRedirect)
	self.Send(v)
}

// HTTP 4XX
// ----------------------------------------
func (self *Context) BadRequest(v interface{}) {
	self.WriteHeader(http.StatusBadRequest)
	self.Send(v)
}

func (self *Context) Unauthorized(v interface{}) {
	self.WriteHeader(http.StatusUnauthorized)
	self.Send(v)
}

func (self *Context) PaymentRequired(v interface{}) {
	self.WriteHeader(http.StatusPaymentRequired)
	self.Send(v)
}

func (self *Context) Forbidden(v interface{}) {
	self.WriteHeader(http.StatusForbidden)
	self.Send(v)
}

func (self *Context) NotFound(v interface{}) {
	self.WriteHeader(http.StatusNotFound)
	self.Send(v)
}

func (self *Context) MethodNotAllowed(v interface{}) {
	self.WriteHeader(http.StatusMethodNotAllowed)
	self.Send(v)
}

func (self *Context) NotAcceptable(v interface{}) {
	self.WriteHeader(http.StatusNotAcceptable)
	self.Send(v)
}

func (self *Context) ProxyAuthRequired(v interface{}) {
	self.WriteHeader(http.StatusProxyAuthRequired)
	self.Send(v)
}

func (self *Context) RequestTimeout(v interface{}) {
	self.WriteHeader(http.StatusRequestTimeout)
	self.Send(v)
}

func (self *Context) Conflict(v interface{}) {
	self.WriteHeader(http.StatusConflict)
	self.Send(v)
}

func (self *Context) Gone(v interface{}) {
	self.WriteHeader(http.StatusGone)
	self.Send(v)
}

func (self *Context) LengthRequired(v interface{}) {
	self.WriteHeader(http.StatusLengthRequired)
	self.Send(v)
}

func (self *Context) PreconditionFailed(v interface{}) {
	self.WriteHeader(http.StatusPreconditionFailed)
	self.Send(v)
}

func (self *Context) RequestEntityTooLarge(v interface{}) {
	self.WriteHeader(http.StatusRequestEntityTooLarge)
	self.Send(v)
}

func (self *Context) RequestURITooLong(v interface{}) {
	self.WriteHeader(http.StatusRequestURITooLong)
	self.Send(v)
}

func (self *Context) UnsupportedMediaType(v interface{}) {
	self.WriteHeader(http.StatusUnsupportedMediaType)
	self.Send(v)
}

func (self *Context) RequestedRangeNotSatisfiable(v interface{}) {
	self.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
	self.Send(v)
}

func (self *Context) ExpectationFailed(v interface{}) {
	self.WriteHeader(http.StatusExpectationFailed)
	self.Send(v)
}

func (self *Context) Teapot(v interface{}) {
	self.WriteHeader(http.StatusTeapot)
	self.Send(v)
}

// HTTP 5XX
// ----------------------------------------
func (self *Context) InternalServerError(v interface{}) {
	self.WriteHeader(http.StatusInternalServerError)
	self.Send(v)
}

func (self *Context) NotImplemented(v interface{}) {
	self.WriteHeader(http.StatusNotImplemented)
	self.Send(v)
}

func (self *Context) BadGateway(v interface{}) {
	self.WriteHeader(http.StatusBadGateway)
	self.Send(v)
}

func (self *Context) ServiceUnavailable(v interface{}) {
	self.WriteHeader(http.StatusServiceUnavailable)
	self.Send(v)
}

func (self *Context) GatewayTimeout(v interface{}) {
	self.WriteHeader(http.StatusGatewayTimeout)
	self.Send(v)
}

func (self *Context) HTTPVersionNotSupported(v interface{}) {
	self.WriteHeader(http.StatusHTTPVersionNotSupported)
	self.Send(v)
}

/*
// New HTTP status codes from RFC 6585. Not exported yet in Go 1.1.
// See discussion at https://codereview.appspot.com/7678043/
PreconditionRequired          = 428
TooManyRequests               = 429
RequestHeaderFieldsTooLarge   = 431
NetworkAuthenticationRequired = 511
*/
