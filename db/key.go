/* ----------------------------------------------------------------------
 *       ______      ___                         __
 *      / ____/___  /   |  ____  __  ___      __/ /_  ___  ________
 *     / / __/ __ \/ /| | / __ \/ / / / | /| / / __ \/ _ \/ ___/ _ \
 *    / /_/ / /_/ / ___ |/ / / / /_/ /| |/ |/ / / / /  __/ /  /  __/
 *    \____/\____/_/  |_/_/ /_/\__. / |__/|__/_/ /_/\___/_/   \___/
 *                            /____/
 *
 * (C) Copyright 2014 GoAnywhere (http://goanywhere.io).
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

package db

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

// Ref: http://www.mongodb.org/display/DOCS/Object+IDs
type Key []byte

// counter is atomically incremented when generating a new Key
// using NewKey() function. It's used as a counter part of a key.
var counter uint32 = 0

// FIXME use more reliable & fixed identity to replace this hostname based mid.
// machineId stores machine id generated once and used in subsequent calls
// to NewKey function.
var machineId []byte

func NewKey() Key {
	var bytes [12]byte
	// Timestamp, 4 bytes, big endian
	binary.BigEndian.PutUint32(bytes[:], uint32(time.Now().Unix()))
	// Machine, first 3 bytes of md5(hostname)
	bytes[4] = machineId[0]
	bytes[5] = machineId[1]
	bytes[6] = machineId[2]
	// Pid, 2 bytes, specs don't specify endianness, but we use big endian.
	pid := os.Getpid()
	bytes[7] = byte(pid >> 8)
	bytes[8] = byte(pid)
	// Increment, 3 bytes, big endian
	index := atomic.AddUint32(&counter, 1)
	bytes[9] = byte(index >> 16)
	bytes[10] = byte(index >> 8)
	bytes[11] = byte(index)
	return Key(bytes[:])
}

func (self Key) Hex() string {
	return hex.EncodeToString(self)
}

func (self Key) Time() time.Time {
	// bytes[0:4] of Key is 32-bit big-endian seconds from epoch.
	secs := int64(binary.BigEndian.Uint32(self[0:4]))
	return time.Unix(secs, 0)
}

// Machine returns the 3-byte machine id part of the key.
func (self Key) Machine() string {
	return hex.EncodeToString(self[4:7])
}

// ProcessId returns the process id part of the key.
func (self Key) ProcessId() uint16 {
	return binary.BigEndian.Uint16(self[7:9])
}

// Counter returns the incrementing value part of the key.
func (self Key) Counter() int32 {
	bytes := self[9:]
	return int32(uint32(bytes[0])<<16 | uint32(bytes[1])<<8 | uint32(bytes[2]))
}

func (self Key) String() string {
	return fmt.Sprintf(`Key("%x")`, string(self))
}

func init() {
	var sum [3]byte
	id := sum[:]
	if hostname, err := os.Hostname(); err == nil {
		copy(id, hmac.New(md5.New, []byte(hostname)).Sum(nil))
		machineId = id
	} else {
		panic(fmt.Errorf("Can not fetch hostname: %v", err))
	}
}
