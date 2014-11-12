/**
 *  ------------------------------------------------------------
 *  @project	web.go
 *  @file       uuid.go
 *  @date       2014-11-12
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
package crypto

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"net"
	"regexp"
	"sync"
	"time"
)

const (
	// Difference in 100-nanosecond intervals between UUID epoch (October 15, 1582) and Unix epoch (January 1, 1970).
	epoch = 122192928000000000

	VariantNCS = iota
	VariantRFC4122
	VariantMicrosoft
	VariantFuture
)

var (
	mutex    sync.Mutex
	sequence uint16
	lastTime uint64
	hardware [6]byte
	pattern  = regexp.MustCompile(`^(urn\:uuid\:)?[\{(\[]?([A-Fa-f0-9]{8})-?([A-Fa-f0-9]{4})-?([1-5][A-Fa-f0-9]{3})-?([A-Fa-f0-9]{4})-?([A-Fa-f0-9]{12})[\]\})]?$`)

	NamespaceDNS, _  = ParseUUID("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	NamespaceURL, _  = ParseUUID("6ba7b811-9dad-11d1-80b4-00c04fd430c8")
	NamespaceOID, _  = ParseUUID("6ba7b812-9dad-11d1-80b4-00c04fd430c8")
	NamespaceX500, _ = ParseUUID("6ba7b814-9dad-11d1-80b4-00c04fd430c8")
)

type UUID [16]byte

// ---------------------------------------------------------------------------
//  Internal Helpers
// ---------------------------------------------------------------------------
// now returns epoch timestamp (for V1/V2).
func now() uint64 {
	mutex.Lock()
	defer mutex.Unlock()

	// difference in 100-nanosecond intervals between UUID epoch (October 15, 1582) and now.
	now := epoch + uint64(time.Now().UnixNano()/100)

	// Clock changed backwards since last UUID generation.
	// Should increase clock sequence.
	if now <= lastTime {
		sequence++
	}
	lastTime = now

	return now
}

// digest a namespace UUID and a name using the given algorithm, which then marshals to a new UUID
func digest(hash hash.Hash, namespace UUID, name string) UUID {
	hash.Write(namespace[:])
	hash.Write([]byte(name))

	uuid := UUID{}
	copy(uuid[:], hash.Sum(nil))
	return uuid
}

// use specifies the version & variant for UUID.
func (self *UUID) use(version byte) {
	self[6] = (self[6] & 0x0f) | (version << 4)
	self[8] = (self[8] & 0xbf) | 0x80
}

// ---------------------------------------------------------------------------
//  Public APIs
// ---------------------------------------------------------------------------
// Version returns the version used to generate UUID.
func (self *UUID) Version() uint {
	return uint(self[6] >> 4)
}

// String converts UUID into plain string.
func (self UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", self[:4], self[4:6], self[6:8], self[8:10], self[10:])
}

// Parse converts the given raw string to UUID.
// Supported formats:
//		6ba7b8149dad11d180b400c04fd430c8
//		6ba7b814-9dad-11d1-80b4-00c04fd430c8
//		{6ba7b814-9dad-11d1-80b4-00c04fd430c8}
//		urn:uuid:6ba7b814-9dad-11d1-80b4-00c04fd430c8
//		[6ba7b814-9dad-11d1-80b4-00c04fd430c8]
func ParseUUID(input string) (uuid UUID, err error) {
	raw := pattern.FindStringSubmatch(input)
	if raw == nil {
		err = errors.New("Invalid UUID string")
		return
	}
	hash := raw[2] + raw[3] + raw[4] + raw[5] + raw[6]
	bytes, err := hex.DecodeString(hash)
	if err != nil {
		return
	}
	copy(uuid[:], bytes)
	return
}

// NewV1 generates a UUID from hardware address, sequence number, and the current time.
func NewV1() UUID {
	uuid := UUID{}
	now := now()
	binary.BigEndian.PutUint32(uuid[0:], uint32(now))
	binary.BigEndian.PutUint16(uuid[4:], uint16(now>>32))
	binary.BigEndian.PutUint16(uuid[6:], uint16(now>>48))
	binary.BigEndian.PutUint16(uuid[8:], sequence)
	copy(uuid[10:], hardware[:])
	uuid.use(1)
	return uuid
}

// NewV3 generates a UUID from the MD5 hash of a namespace UUID and a unique key.
func NewV3(namespace UUID, key string) UUID {
	uuid := digest(md5.New(), namespace, key)
	uuid.use(3)
	return uuid
}

// NewV4 generates a random UUID.
func NewV4() UUID {
	uuid := UUID{}
	rand.Read(uuid[:])
	uuid.use(4)
	return uuid
}

// NewV5 generates a UUID from the SHA-1 hash of a namespace UUID and a unique key.
func NewV5(namespace UUID, key string) UUID {
	uuid := digest(sha1.New(), namespace, key)
	uuid.use(5)
	return uuid
}

func init() {
	buf := make([]byte, 2)
	rand.Read(buf)
	sequence = binary.BigEndian.Uint16(buf)

	// in case the real one's absence.
	rand.Read(hardware[:])
	hardware[0] |= 0x01

	if interfaces, err := net.Interfaces(); err == nil {
		for _, iface := range interfaces {
			if len(iface.HardwareAddr) >= 6 {
				copy(hardware[:], iface.HardwareAddr)
				break
			}
		}
	}
}
