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

package crypt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
)

// A Salter provides a random data generator to password salt and unique session
// IDs.
type Salter struct {
	salt       []byte
	rndSources RandomSourceList
	sumWeight  int
}

// NewSalter creates a new instance of Salter. It requires a list of sources of
// random data and optionally an input to salt next token.
func NewSalter(sources RandomSourceList, input []byte) *Salter {
	if input == nil {
		input = []byte("")
	}

	hash := sha256.New()
	hash.Write(input)

	sum := 0
	for _, v := range sources {
		sum += v.Weight
	}

	return &Salter{
		salt:       hash.Sum(nil),
		rndSources: sources,
		sumWeight:  sum,
	}
}

// BToken generates an array of random bytes with length as specified by size
// parameter.
func (self *Salter) BToken(size int) []byte {
	mac := hmac.New(sha256.New, self.salt)

	for _, v := range self.rndSources {
		itemSize := float32(size) * (float32(v.Weight) / float32(self.sumWeight))
		mac.Write(getRandomBytes(v.Reader, int(itemSize)))
	}
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
