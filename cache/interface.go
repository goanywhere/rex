/**
 *  ------------------------------------------------------------
 *  @project	web.go
 *  @file       interface.go
 *  @date       2014-11-11
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
package cache

// Storage serve as a generic key value cache storage interface.
// The storage may be persistent (e.g., a database) or volatile (e.g., cache).
// All Storage implementations must be thread safe.
type Storage interface {
	// Get retrieves a single value from the storage.
	Get(key interface{}) (interface{}, error)

	// GetMulti retrieves multiple values from the storage.
	// The items are returned in the same order as the input keys.
	GetMulti(keys ...interface{}) ([]interface{}, error)

	// Set stores a single item into the storage.
	Set(item interface{}) error

	// SetMulti stores multiple items into the storage.
	SetMulti(items ...interface{}) error

	// Delete removes a single item from the storage.
	Delete(key interface{}) error

	// DeleteMulti removes multiple items from the storage.
	DeleteMulti(keys ...interface{}) error

	// Flush wipes all items from the storage.
	Flush() error
}
