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

type redis struct {
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration

	*internal.Sharding
	pools []*redigo.Pool
}

func New() *redis {
	servers := settings.Strings("CACHE_REDIS_SERVERS")

	create := func(rawurl string) *redigo.Pool {
		pool := new(redigo.Pool)
		pool.MaxIdle = settings.Int("CACHE_REDIS_MAXIDLE", 10)
		pool.MaxActive = settings.Int("CACHE_REDIS_MAXACTIVE", 100)
		pool.IdleTimeout = time.Duration(settings.Int("CACHE_REDIS_IDLETIMEOUT", 180)) * time.Second

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
		log.Fatalf("Failed to setup Redis: 'CACHE_REDIS_SERVERS' missing?")
	}
	redis := new(redis)
	redis.Sharding = internal.NewSharding(len(servers))
	redis.pools = pools
	return redis
}

// conn gets connection from pool with the given key with sharding supports.
func (self *redis) conn(key string) redigo.Conn {
	var pool *redigo.Pool
	if len(self.pools) == 1 {
		pool = self.pools[0]
	} else {
		pool = self.pools[self.Shard(key)]
	}
	return pool.Get()
}

// exec raw Redis commands with interface & error in return.
// commands are automatically sharded using key when there are more than a server.
func (self *redis) exec(cmd, key string, args ...interface{}) (interface{}, error) {
	conn := self.conn(key)
	defer conn.Close()
	return conn.Do(cmd, append([]interface{}{key}, args...)...)
}

// Get value associated with the given key & assign to the given pointer.
func (self *redis) Get(key string, ptr interface{}) error {
	if reply, err := self.exec("GET", key); err == nil {
		if reply == nil {
			return redigo.ErrNil
		}
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
		_, err = self.exec("SETEX", key, int32(expires/time.Second), bits)
	} else {
		_, err = self.exec("SET", key, bits)
	}
	return err
}

// Delete value associated with the given key from cache.
func (self *redis) Del(key string, others ...string) error {
	_, err := self.exec("DEL", key)
	for _, item := range others {
		_, err = self.exec("DEL", item)
	}
	return err
}

// Increments the number assoicated with the given key by one.
// If the key does not exist, it is set to 0 before performing the operation.
// An error will be returned if the key contains a value of the wrong type
// or contains a string that can not be represented as integer.
func (self *redis) Incr(key string) error {
	_, err := self.exec("INCR", key)
	return err
}

// Decrements the number assoicated with the given key by one.
// If the key does not exist, it is set to 0 before performing the operation.
// An error will be returned if the key contains a value of the wrong type
// or contains a string that can not be represented as integer.
func (self *redis) Decr(key string) error {
	_, err := self.exec("DECR", key)
	return err
}

// Determine if a key exists.
func (self *redis) Exists(key string) bool {
	exists, _ := redigo.Bool(self.exec("EXISTS", key))
	return exists
}

// Flush all redis cache servers.
func (self *redis) Flush() error {
	var (
		err  error
		conn redigo.Conn
	)
	for _, pool := range self.pools {
		conn = pool.Get()
		_, err = conn.Do("FLUSHDB")
		conn.Close()
	}
	return err
}
