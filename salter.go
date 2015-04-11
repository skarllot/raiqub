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
	"io"
	"time"
)

const (
	// Defines default input size from system random generator to 512-bit
	DEFAULT_SYS_RAND_SIZE = 64
	// Defines default input size from raiqub random generator to 96-bit.
	DEFAULT_RAIQUB_RAND_SIZE = 12
)

// A Salter provides a random data generator to password salt and unique session
// IDs.
type Salter struct {
	salt []byte
	rnd2 *Random
}

// NewSalter creates a new instance of Salter.
func NewSalter() *Salter {
	return NewSalterInput(getRandomBytes(rand.Reader, DEFAULT_SYS_RAND_SIZE))
}

// NewSalterInput creates a new instance of Salter and provides an initial input
// to salt next token.
func NewSalterInput(input []byte) *Salter {
	hash := sha256.New()
	hash.Write(input)

	return &Salter{
		salt: hash.Sum(nil),
		rnd2: NewRandom(),
	}
}

// BToken generates an array of random bytes with length as specified by size
// parameter.
func (self *Salter) BToken(size int) []byte {
	mac := hmac.New(sha256.New, self.salt)

	// Tries to create unpredictable token
	// Most strength comes from 'rand.Reader'
	// Another bits are used to avoid the chance of system random genarator
	//   is compromissed by internal issue
	mac.Write(getRandomBytes(rand.Reader, size))
	mac.Write(getRandomBytes(rand.Reader, time.Now().Second()/6))
	mac.Write(getRandomBytes(self.rnd2, DEFAULT_RAIQUB_RAND_SIZE))
	macSum := mac.Sum(nil)

	self.salt = macSum
	return macSum
}

// DefaultBToken generates an array of random bytes with default length.
func (self *Salter) DefaultBToken() []byte {
	return self.BToken(DEFAULT_SYS_RAND_SIZE)
}

// DefaultToken generates a base-64 string of random bytes with default length.
func (self *Salter) DefaultToken() string {
	return self.Token(DEFAULT_SYS_RAND_SIZE)
}

// Token generates a base-64 string of random bytes with length as specified by
// size parameter.
func (self *Salter) Token(size int) string {
	return base64.URLEncoding.EncodeToString(self.BToken(size))
}

// getRandomBytes gets secure random bytes.
func getRandomBytes(r io.Reader, n int) []byte {
	b := make([]byte, n)
	_, err := r.Read(b)
	if err != nil {
		panic("Could not access secure random generator")
	}

	return b
}
