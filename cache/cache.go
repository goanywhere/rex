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
