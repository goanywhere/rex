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
	"log"
	"net/url"
	"time"

	redigo "github.com/garyburd/redigo/redis"

	"github.com/goanywhere/rex/internal"
)

var Redis = struct {
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}{
	MaxIdle:     options.Int("cache.redis.maxidle", 10),
	MaxActive:   options.Int("cache.redis.maxactive", 100),
	IdleTimeout: time.Duration(options.Int("cache.redis.idletimeout", 180)) * time.Second,
}

type redis struct {
	*internal.Sharding
	pools []*redigo.Pool
}

func New() *redis {
	servers := options.Strings("cache.redis.servers")

	create := func(rawurl string) *redigo.Pool {
		pool := new(redigo.Pool)
		pool.MaxIdle = Redis.MaxIdle
		pool.MaxActive = Redis.MaxActive
		pool.IdleTimeout = Redis.IdleTimeout
		pool.Dial = func() (conn redigo.Conn, err error) {
			var URL *url.URL
			URL, err = url.Parse(rawurl)
			if err != nil {
				log.Fatalf("rex/cache: failed to setup cache server: %v", err)
			}
			// regular connection.
			if conn, err = redigo.Dial("tcp", URL.Host); err != nil {
				return
			}

			if URL.User == nil {
				if _, err = conn.Do("PING"); err != nil {
					conn.Close()
				}
			} else {
				// authenticate connection if password exists.
				if password, exists := URL.User.Password(); exists {
					if _, err = conn.Do("AUTH", password); err != nil {
						conn.Close()
					}
				}
			}
			return
		}
		// custom connection test method
		pool.TestOnBorrow = func(conn redigo.Conn, t time.Time) error {
			if _, err := conn.Do("PING"); err != nil {
				return err
			}
			return nil
		}
		return pool
	}

	var pools []*redigo.Pool
	for _, rawurl := range servers {
		pools = append(pools, create(rawurl))
	}
	if len(pools) == 0 {
		log.Fatalf("Failed to setup Redis: 'rex.cache.redis.servers' missing?")
	}
	redis := new(redis)
	redis.Sharding = internal.NewSharding(len(servers))
	redis.pools = pools
	return redis
}

// Serialization ------------------------------------------------------------

// do raw Redis commands with interface & error in return.
func (self *redis) do(cmd, key string, args ...interface{}) (interface{}, error) {
	var pool *redigo.Pool
	if len(self.pools) == 1 {
		pool = self.pools[0]
	} else {
		pool = self.pools[self.Shard(key)]
	}
	conn := pool.Get()
	defer conn.Close()
	return conn.Do(cmd, append([]interface{}{key}, args...)...)
}

// Get value associated with the given key & assign to the given pointer.
func (self *redis) Get(key string, ptr interface{}) error {
	if reply, err := self.do("GET", key); err == nil {
		return internal.Deserialize(reply.([]byte), ptr)
	}
	return nil
}

// Set the key/value into the cache with specified time duration.
func (self *redis) Set(key string, value interface{}, expires time.Duration) error {
	bits, err := internal.Serialize(value)
	if err != nil {
		return err
	}

	if expires > 0 {
		_, err = self.do("SETEX", key, int32(expires/time.Second), bits)
	} else {
		_, err = self.do("SET", key, bits)
	}
	return err
}

// Delete value associated with the given key from cache.
func (self *redis) Del(key string, others ...string) error {
	_, err := self.do("DEL", key)
	for _, item := range others {
		_, err = self.do("DEL", item)
	}
	return err
}

// Increments the number assoicated with the given key by one.
// If the key does not exist, it is set to 0 before performing the operation.
// An error will be returned if the key contains a value of the wrong type
// or contains a string that can not be represented as integer.
func (self *redis) Incr(key string) error {
	_, err := self.do("INCR", key)
	return err
}

// Decrements the number assoicated with the given key by one.
// If the key does not exist, it is set to 0 before performing the operation.
// An error will be returned if the key contains a value of the wrong type
// or contains a string that can not be represented as integer.
func (self *redis) Decr(key string) error {
	_, err := self.do("DECR", key)
	return err
}

// Determine if a key exists.
func (self *redis) Exists(key string) bool {
	exists, _ := redigo.Bool(self.do("EXISTS", key))
	return exists
}

// ---------- Primitive Types Supports ----------

func (self *redis) String(key string) (string, error) {
	return redigo.String(self.do("GET", key))
}

func (self *redis) Flush() error {
	return nil
}
