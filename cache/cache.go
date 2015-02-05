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
package cache

import (
	"time"

	"github.com/goanywhere/rex/internal"
)

var settings = internal.Settings()

type Cache interface {
	// Get value associated with the given key & assign to the given pointer.
	Get(key string, ptr interface{}) error

	// Set the key/value into the cache with specified time duration.
	Set(key string, value interface{}, expires time.Duration) error

	// Delete value associated with the given key from cache.
	Del(key string) error

	// Increments the number assoicated with the given key by one.
	// If the key does not exist, it is set to 0 before performing the operation.
	// An error will be returned if the key contains a value of the wrong type
	// or contains a string that can not be represented as integer.
	Incr(key string) error

	// Decrements the number assoicated with the given key by one.
	// If the key does not exist, it is set to 0 before performing the operation.
	// An error will be returned if the key contains a value of the wrong type
	// or contains a string that can not be represented as integer.
	Decr(key string) error

	// Determine if a key exists.
	Exists(key string) bool

	Flush() error
}
