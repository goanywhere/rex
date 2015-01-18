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

package crypto

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	errSignatureInvalid = errors.New("The signature is invalid")
	errSignatureExpired = errors.New("The signature is already expired")
	regexSignature      = regexp.MustCompile(`((\d{19})\|(\w+)\|(\w{40}))`)
)

type Signature struct {
	secret  string
	timeout time.Duration
}

// NewSignature creates a signature with secret salt for encode/decode consequent values.
func NewSignature(secret string) *Signature {
	if secret == "" {
		log.Fatal("Failed to create signature: secret key missing")
	}
	signature := new(Signature)
	signature.secret = secret
	signature.timeout = time.Hour * 24 * 365
	return signature
}

// expired checks if the given time (in nanoseconds) is valid.
func (self *Signature) expired(nano int64) bool {
	now := time.Now()
	issue := time.Unix(0, nano)
	if now.Sub(issue) >= self.timeout {
		return true
	}
	// Ensure the signature is not from the *future*, allow 1 minute grace period.
	if issue.After(now.Add(1 * time.Minute)) {
		return true
	}
	return false
}

// Encode creates a checksum value for the given source using key, timestamp, crc.
// Pattern:
//	1) key|nano|src(hex)	=> hmac.sha1		=> crc
//	2) nano|src(hex)|crc	=> base64 encode	=> signed value
// NOTE struct value must be registered using gob.Register() first.
func (self *Signature) Encode(key string, src []byte) (value string, err error) {
	nano := time.Now().UnixNano()
	// 1) setup hash with secret salt, construct CRC using "key|nano|src(hex)" via hash.
	hash := hmac.New(sha1.New, []byte(self.secret))
	hash.Write([]byte(fmt.Sprintf("%s|%d|%x", key, nano, src)))
	// 2) construct raw values for base64 encoding using "nano|src(hex)|crc".
	raw := fmt.Sprintf("%d|%x|%s", nano, src, hex.EncodeToString(hash.Sum(nil)))
	value = base64.URLEncoding.EncodeToString([]byte(raw))
	return
}

// Decode unpacks the source string to original values.
// Pattern:
//	1) signed value		=> base64 decode	=> nano|src(hex)|crc
//	2) verify crc
//	3) verify nano timestamp
func (self *Signature) Decode(key, value string) (src []byte, err error) {
	if bits, err := base64.URLEncoding.DecodeString(value); err == nil {
		if regexSignature.Match(bits) {
			// values: nano|src(hex)|crc
			values := strings.Split(string(bits), "|")
			if nano, err := strconv.ParseInt(values[0], 0, 64); err == nil {
				// 1) reconstruct CRC using "key|nano|src(hex)" via hash.
				hash := hmac.New(sha1.New, []byte(self.secret))
				hash.Write([]byte(fmt.Sprintf("%s|%s|%s", key, values[0], values[1])))
				// 2) verify the incoming CRC.
				if values[2] == hex.EncodeToString(hash.Sum(nil)) {
					// 3) verify nano timestamp
					if self.expired(nano) {
						err = errSignatureExpired
					} else {
						src, err = hex.DecodeString(values[1])
					}
				} else {
					err = errSignatureInvalid
				}
			}
		} else {
			err = errSignatureInvalid
		}
	}
	return
}
