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
	"bufio"
	"errors"
	"net"
	"net/http"
)

/*
Extension to http.WriterWriter with compression supports.
*/

type Writer interface {
	http.Flusher
	http.Hijacker
	http.CloseNotifier
	http.ResponseWriter

	Size() int
	Status() int
	Written() bool
}

type writer struct {
	http.ResponseWriter

	size   int
	status int
}

/* ----------------------------------------------------------------------
 * Implementations of http.Hijacker
 * ----------------------------------------------------------------------*/
func (self *writer) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := self.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("ResponseWriter doesn't support the Hijacker interface")
	}
	return hijacker.Hijack()
}

/* ----------------------------------------------------------------------
 * Implementations of http.CloseNotifier
 * ----------------------------------------------------------------------*/
func (self *writer) CloseNotify() <-chan bool {
	return self.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

/* ----------------------------------------------------------------------
 * Implementations of http.Flusher
 * ----------------------------------------------------------------------*/
func (self *writer) Flush() {
	flusher, ok := self.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}

/* ----------------------------------------------------------------------
 * Implementations of http.ResponseWriter
 * ----------------------------------------------------------------------*/
func (self *writer) WriteHeader(status int) {
	if status >= 100 && status < 512 {
		self.status = status
		self.ResponseWriter.WriteHeader(status)
	}
}

// Write: Implementation of http.ResponseWriter#Write
func (self *writer) Write(data []byte) (size int, err error) {
	size, err = self.ResponseWriter.Write(data)
	self.size += size
	return
}

/* ----------------------------------------------------------------------
 * Implementations of rex.Writer interface.
 * ----------------------------------------------------------------------*/
func (self *writer) Size() int {
	return self.size
}

// Status returns current status code of the Context.
func (self *writer) Status() int {
	return self.status
}

func (self *writer) Written() bool {
	return self.status != 0 || self.size > 0
}
