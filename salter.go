/*
 * Copyright (C) 2015 Fabrício Godoy <skarllot@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 59 Temple Place - Suite 330, Boston, MA  02111-1307, USA.
 */

package raiqub

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"
)

const (
	// Defines default salt size to 1024-bit
	DEFAULT_SALT_SIZE = 128
	// Defines default sleep time to create unpredictability.
	DEFAULT_SLEEP_TIME = time.Microsecond * 10
	// Defines default number of rounds to generate time-based random data.
	DEFAULT_SLEEP_ROUNDS = 20
)

// A Salter provides a random data generator to password salt and unique session
// IDs.
type Salter struct {
	salt     []byte
	lastDate time.Time
}

// NewSalter creates a new instance of Salter.
func NewSalter() *Salter {
	return NewSalterInput(getRandomBytes(DEFAULT_SALT_SIZE))
}

// NewSalterInput creates a new instance of Salter and provides an initial input
// to salt next token.
func NewSalterInput(input []byte) *Salter {
	hash := sha256.New()
	hash.Write(input)

	return &Salter{
		salt:     hash.Sum(nil),
		lastDate: time.Now(),
	}
}

// BToken generates an array of random bytes with length as specified by size
// parameter.
func (self *Salter) BToken(size int) []byte {
	mac := hmac.New(sha256.New, self.salt)

	// Tries to create unpredictable token
	// Most strength comes from 'rand.Read'
	// Another bits are used to avoid the chance of system random genarator
	//   is compromissed by internal issue
	mac.Write(getRandomBytes(size))
	mac.Write(getRandomBytes(time.Now().Second() / 6))
	mac.Write(getRandomBytes2(self.lastDate))
	macSum := mac.Sum(nil)

	self.salt = macSum
	self.lastDate = time.Now()
	return macSum
}

// DefaultBToken generates an array of random bytes with default length.
func (self *Salter) DefaultBToken() []byte {
	return self.BToken(DEFAULT_SALT_SIZE)
}

// DefaultToken generates a base-64 string of random bytes with default length.
func (self *Salter) DefaultToken() string {
	return self.Token(DEFAULT_SALT_SIZE)
}

// Token generates a base-64 string of random bytes with length as specified by
// size parameter.
func (self *Salter) Token(size int) string {
	return base64.URLEncoding.EncodeToString(self.BToken(size))
}

// getRandomBytes gets secure random bytes.
func getRandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic("Could not access secure random generator")
	}

	return b
}

// getRandomBytes2 gets a time-based random bytes.
func getRandomBytes2(input time.Time) []byte {
	binput := []byte(input.Format(time.RFC3339Nano))
	diff := time.Now().Sub(input)
	now := time.Now()

	// Tries to create unpredictable bytes limited by time.
	// Most strength comes from time differences.
	hash := hmac.New(sha256.New, binput)
	hash.Write(getBytesFromInt64(diff.Nanoseconds(), 2))
	hash.Write([]byte(now.Format(time.RFC3339Nano)))

	// Each round generates 2 significant random bytes. But only 1 byte is
	// really strong.
	for i := 0; i < DEFAULT_SLEEP_ROUNDS; i++ {
		// The "sleep" time is almost unpredictable.
		time.Sleep(DEFAULT_SLEEP_TIME)
		now = time.Now()
		diff = now.Sub(input)

		hash.Write(getBytesFromInt64(diff.Nanoseconds(), 2))
		hash.Write([]byte(now.Format(time.RFC3339Nano)))
	}

	return hash.Sum(nil)
}

// getBytesFromInt64 decomposes each Int64 byte to byte array.
func getBytesFromInt64(val int64, size int) []byte {
	result := make([]byte, size)

	for i := 0; i < size; i++ {
		result[i] = byte(val >> uint(8*i))
	}
	return result
}
