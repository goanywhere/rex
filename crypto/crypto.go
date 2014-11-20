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

package crypto

import (
	"crypto/hmac"
	"encoding/hex"
	"hash"
)

// Hash creates a hex string using the given crypto.
// Example:
//	Hash(sha1.New, "secret salt", "encrypt me please")
func Hash(crypto func() hash.Hash, salt string, src []byte) string {
	h := hmac.New(crypto, []byte(salt))
	h.Write(src)
	return hex.EncodeToString(h.Sum(nil))
}
