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
)

type Salter struct {
	salt []byte
}

func NewSalter() *Salter {
	hash := sha256.New()
	hash.Write(getRandomBytes(DEFAULT_SALT_SIZE))

	return &Salter{
		salt: hash.Sum(nil),
	}
}

func (self *Salter) BToken(size int) []byte {
	mac := hmac.New(sha256.New, self.salt)
	now := time.Now().Format(time.RFC3339Nano)

	// Tries to create unpredictable token
	// Most strength comes from 'rand.Read'
	// Another bits are used to avoid the chance of system random genarator
	//   is compromissed by internal issue
	mac.Write(getRandomBytes(size))
	mac.Write(getRandomBytes(time.Now().Second() / 2))
	mac.Write([]byte(now))
	macSum := mac.Sum(nil)
	self.salt = macSum
	return macSum
}

func (self *Salter) DefaultBToken() []byte {
	return self.BToken(DEFAULT_SALT_SIZE)
}

func (self *Salter) DefaultToken() string {
	return self.Token(DEFAULT_SALT_SIZE)
}

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
