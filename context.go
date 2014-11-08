/**
 *  ------------------------------------------------------------
 *  @project	web.go
 *  @file       context.go
 *  @date       2014-10-21
 *  @author     Jim Zhan <jim.zhan@me.com>
 *
 *  Copyright © 2014 Jim Zhan.
 *  ------------------------------------------------------------
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://wwself.apache.org/licenses/LICENSE-2.0
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
	"bufio"
	"errors"
	"net"
	"net/http"
)

type (
	ResponseWriter interface {
		http.ResponseWriter
		http.Hijacker
		http.Flusher
		http.CloseNotifier

		Status() int
		Size() int
		Written() bool
	}

	Context struct {
		http.ResponseWriter
		status  int
		size    int
		Request *http.Request
		data    map[string]interface{}
	}
)

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	var ctx *Context = new(Context)
	ctx.ResponseWriter = w
	ctx.Request = r
	ctx.status = http.StatusOK
	ctx.size = -1
	return ctx
}

// ---------------------------------------------------------------------------
//  Enhancements for native http.ResponseWriter
// ---------------------------------------------------------------------------
func (self *Context) Status() int {
	return self.status
}

func (self *Context) Size() int {
	return self.size
}

func (self *Context) Written() bool {
	return self.size != -1
}

// ---------------------------------------------------------------------------
//  Implementation of http.ResponseWriter#WriteHeader
// ---------------------------------------------------------------------------
func (self *Context) WriteHeader(status int) {
	if status > 0 && !self.Written() {
		self.status = status
		self.ResponseWriter.WriteHeader(status)
	}
}

// ---------------------------------------------------------------------------
//  Implementation of http.ResponseWriter#Write
// ---------------------------------------------------------------------------
func (self *Context) Write(data []byte) (n int, err error) {
	if !self.Written() {
		self.WriteHeader(http.StatusOK)
	}
	size, err := self.ResponseWriter.Write(data)
	self.size += size
	return size, err
}

// ---------------------------------------------------------------------------
//  Implementations of http.Hijackeri#Hijack
// ---------------------------------------------------------------------------
func (self *Context) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := self.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("the ResponseWriter doesn't support the Hijacker interface")
	}
	return hijacker.Hijack()
}

// ---------------------------------------------------------------------------
//  Implementations of http.CloseNotifier#CloseNotify
// ---------------------------------------------------------------------------
func (self *Context) CloseNotify() <-chan bool {
	return self.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

// ---------------------------------------------------------------------------
//  Implementations of http.Flusher#Flush
// ---------------------------------------------------------------------------
func (self *Context) Flush() {
	flusher, ok := self.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}
