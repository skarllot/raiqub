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
	"crypto/rand"
	"io"
)

const (
	// Defines default input size from system random generator to 512-bit
	DEFAULT_SYS_RAND_SIZE = 64
	// Defines default input size from raiqub random generator to 96-bit.
	DEFAULT_RAIQUB_RAND_SIZE = 12
)

// A RandomSource defines a source of random data and its weight from total.
type RandomSource struct {
	// The reader of random data.
	Reader io.Reader
	// The weight of current random source.
	Weight int
}

// A RandomSourceList defines a slice of RandomSource.
type RandomSourceList []RandomSource

// NewRandomSourceList returns system provided source of random data.
func NewRandomSourceList() RandomSourceList {
	return RandomSourceList{
		RandomSource{
			Reader: rand.Reader,
			Weight: DEFAULT_SYS_RAND_SIZE,
		},
	}
}

// NewRandomSourceListRaiqub returns Raiqub provided source of random data.
func NewRandomSourceListRaiqub() RandomSourceList {
	return RandomSourceList{
		RandomSource{
			Reader: NewRandom(),
			Weight: DEFAULT_RAIQUB_RAND_SIZE,
		},
	}
}

// NewRandomSourceListSecure returns system provided source of random data
// besides of Raiqub Random source.
func NewRandomSourceListSecure() RandomSourceList {
	return RandomSourceList{
		RandomSource{
			Reader: rand.Reader,
			Weight: DEFAULT_SYS_RAND_SIZE,
		},
		RandomSource{
			Reader: NewRandom(),
			Weight: DEFAULT_RAIQUB_RAND_SIZE,
		},
	}
}
