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
	"math/rand"
	"time"
)

var (
	alphanum = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	random   *rand.Rand
)

// RandomString creates a securely generated random string.
//
//	Args:
//		length: length of the generated random string.
func RandomString(length int, chars []rune) string {
	bytes := make([]rune, length)

	var pool []rune
	if chars == nil {
		pool = alphanum
	} else {
		pool = chars
	}

	for index := range bytes {
		bytes[index] = pool[random.Intn(len(pool))]
	}
	return string(bytes)
}

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}
