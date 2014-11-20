/**
 * ----------------------------------------------------------------------
 *       ______      ___                         __
 *      / ____/___  /   |  ____  __  ___      __/ /_  ___  ________
 *     / / __/ __ \/ /| | / __ \/ / / / | /| / / __ \/ _ \/ ___/ _ \
 *    / /_/ / /_/ / ___ |/ / / / /_/ /| |/ |/ / / / /  __/ /  /  __/
 *    \____/\____/_/  |_/_/ /_/\__. / |__/|__/_/ /_/\___/_/   \___/
 *                            /____/
 * ----------------------------------------------------------------------
 *  Copyright © 2014 GoAnywhere Ltd. All Rights Reserved.
 * ----------------------------------------------------------------------*/

package web

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
)

// Deserialize converts base64-encoded string back to its original object.
func Deserialize(value string, object interface{}) (err error) {
	if bits, err := base64.URLEncoding.DecodeString(value); err == nil {
		err = gob.NewDecoder(bytes.NewBuffer(bits)).Decode(object)
	}
	return
}

// Serialize converts any given object into base64-encoded string using `encoding/gob`.
// NOTE struct must be registered using gob.Register() first.
func Serialize(object interface{}) (value string, err error) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	if err = encoder.Encode(object); err == nil {
		value = base64.URLEncoding.EncodeToString(buffer.Bytes())
	}
	return
}
