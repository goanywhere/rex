/* ----------------------------------------------------------------------
 *       ______      ___                         __
 *      / ____/___  /   |  ____  __  ___      __/ /_  ___  ________
 *     / / __/ __ \/ /| | / __ \/ / / / | /| / / __ \/ _ \/ ___/ _ \
 *    / /_/ / /_/ / ___ |/ / / / /_/ /| |/ |/ / / / /  __/ /  /  __/
 *    \____/\____/_/  |_/_/ /_/\__. / |__/|__/_/ /_/\___/_/   \___/
 *                            /____/
 *
 * (C) Copyright 2014 GoAnywhere (http://goanywhere.io).
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

// A space efficient permutation-based consistent hashing function.  This
// implementation supports up to a maximum of (1 << 16 - 1), 65535, number
// of shards.
//
// Implementation details:
//
// Unlike the standard ring-based algorithm (e.g., as described in dynamo db),
// this algorithm relays on shard permutations to determine the key's shard
// mapping. The idea is as follow:
//  1. Assume there exist a set of shard ids, S, which contains every possible
//     shard ids in the universe (in this case 0 .. 65535).
//  2. Now suppose, A (a subset of S), is the set of available shard ids, and we
//     want to find the shard mapping for key, K
//  3. Use K as the pseudorandom generator's seed, and generate a random
//     permutation of S using variable-base permutation encoding (see
//     http://stackoverflow.com/questions/1506078/fast-permutation-number-permutation-mapping-algorithms
//     for additional details)
//  4. Ignore all shard ids in the permutation that are not in set A
//  5. Finally, use the first shard id as K's shard mapping.
//
// NOTE: Because each key generates a different permutation, the data
// distribution is generally more uniform than the standard algorithm (The
// standard algorithm works around this issue by adding more points to the
// ring, which unfortunately uses even more memory).
//
// Complexity: this algorithm is O(1) in theory (because the max shard id is
// known), but O(n) in practice.
//
// Example:
//  1. Assume S = {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, and A = {0, 1, 2, 3, 4}.
//  2. Now suppose K = 31415 and perm(S, K) = (3, 1, 9, 4, 7, 5, 8, 2, 0, 6).
//  3. After ignoring S - A, the remaining ids are (3, 1, 4, 2, 0)
//  4. Therefore, the key belongs to shard 3.

package sharding

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/binary"
	"strings"
)

const (
	maxShards = 65535
)

type Sharding struct {
	shards uint16
}

// A simplified implementation of 32-bit murmur hash 3, which only accepts
// uint32 value as data, and uses 12345 as the seed.
//
// See https://code.google.com/p/smhasher/wiki/MurmurHash3 for details.
func (self *Sharding) hash(value uint32) uint32 {
	// body
	k := value * 0xcc9e2d51   // k = val * c1
	k = (k << 15) | (k >> 17) // k = rotl32(h, 15)
	k *= 0x1b873593           // k *= c2

	h := 12345 ^ k            // seed ^ k
	h = (h << 13) | (h >> 19) // k = rotl32(h, 13)
	h = h*5 + 0xe6546b64

	// finalize (NOTE: there's no tail)
	h = h ^ 4

	// fmix32
	h ^= h >> 16
	h *= 0x85ebca6b
	h ^= h >> 13
	h *= 0xc2b2ae35
	h ^= h >> 16

	return h
}

func (self *Sharding) Shard(key string) uint16 {
	if self.shards < 2 {
		return 0
	}
	// Converts any string Key into binary unsigned 64-bits integer.
	// *NOTE* In order to control the stablity of the flow (as the binary
	// encoder will throw an error if the given byte is too small), it encode
	// the string via md5 algorithm to fix the length of the source at the
	// very beginning.
	bytes := hmac.New(md5.New, []byte(strings.ToLower(key))).Sum(nil)
	value := binary.BigEndian.Uint64(bytes)
	hash := uint32(value) ^ uint32(value>>32)

	var closestShard uint16 = 0
	var minPosition uint16 = maxShards

	selectClosestShard := func(shard, pos uint16) {
		pos %= (maxShards - shard)
		if pos < minPosition {
			closestShard = shard
			minPosition = pos
		}
	}

	numBlocks := self.shards >> 1
	for i := uint16(0); i < numBlocks; i++ {
		// Each hash can generate 2 permutation positions.  Implementation
		// note: we can replace murmur hash with any other pseudorandom
		// generator, as long as it's sufficiently "random".
		hash = self.hash(hash)

		shard := i << 1

		selectClosestShard(shard, uint16(hash))
		if minPosition == 0 {
			return closestShard
		}

		selectClosestShard(shard+1, uint16(hash>>16))
		if minPosition == 0 {
			return closestShard
		}
	}

	if (self.shards & 0x1) == 1 {
		hash = self.hash(hash)
		selectClosestShard(self.shards-1, uint16(hash))
	}

	return closestShard
}
