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
	"sync"
	"time"
)

const (
	// Defines default sleep time to create unpredictability.
	DEFAULT_SLEEP_TIME = time.Microsecond * 10
)

// A Random provides a pseudo-random generator based on syscall time deltas of
// Sleep calls.
type Random struct {
	val   [2]byte
	index byte
	mutex *sync.Mutex
}

// NewRandom creates a new instance of Random.
func NewRandom() *Random {
	return &Random{
		val:   [...]byte{0, 0},
		index: 0,
		mutex: &sync.Mutex{},
	}
}

// Read fills specified byte array with random data.
// Always return parameter array length and no errors.
func (s *Random) Read(b []byte) (n int, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i, _ := range b {
		b[i] = s.readByte()
	}

	return len(b), nil
}

// readByte creates a single random byte.
func (s *Random) readByte() byte {
	before := time.Now()
	rndDuration := time.Duration(getUInt16FromBytes(s.val))

	time.Sleep(DEFAULT_SLEEP_TIME + rndDuration)
	diff := time.Now().Sub(before)
	n := byte(diff.Nanoseconds())

	s.val[s.index] = n
	s.index ^= 1

	return n
}

// getUInt16FromBytes convert a 2-byte array to 16-bit unsigned integer.
func getUInt16FromBytes(input [2]byte) uint16 {
	return uint16(input[0]) + uint16(input[1])*256
}
