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
	"testing"
)

const (
	DEFAULT_UNPRED_ROUNDS            = 1000
	DEFAULT_PREDICTABILITY_THRESHOLD = .05
)

func testUnpredictability(r io.Reader) int {
	dict := make(map[int16]bool)
	buff := make([]byte, 2)
	count := 0

	for i := 0; i < DEFAULT_UNPRED_ROUNDS; i++ {
		r.Read(buff)
		val := int16(buff[0]) + int16(buff[1])*256

		if _, ok := dict[val]; ok {
			count++
		} else {
			dict[val] = true
		}
	}

	return count
}

func TestSystemRandomUnpredictability(t *testing.T) {
	count := testUnpredictability(rand.Reader)

	if count > DEFAULT_UNPRED_ROUNDS*DEFAULT_PREDICTABILITY_THRESHOLD {
		t.Errorf(
			"System random generator could not generate unpredictable data: %d of %d",
			count, DEFAULT_UNPRED_ROUNDS)
	}
	t.Logf(
		"System random generator predictability: %.2f%%",
		(float32(count)/float32(DEFAULT_UNPRED_ROUNDS))*100)
}

func TestRaiqubUnpredictability(t *testing.T) {
	count := testUnpredictability(NewRandom())

	if count > DEFAULT_UNPRED_ROUNDS*DEFAULT_PREDICTABILITY_THRESHOLD {
		t.Errorf(
			"Raiqub random generator could not generate unpredictable data: %d of %d",
			count, DEFAULT_UNPRED_ROUNDS)
	}
	t.Logf(
		"Raiqub random generator predictability: %.2f%%",
		(float32(count)/float32(DEFAULT_UNPRED_ROUNDS))*100)
}

func BenchmarkRaiqubRandom(b *testing.B) {
	rnd := NewRandom()
	buff := make([]byte, 1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rnd.Read(buff)
	}
}

func BenchmarkRaiqubRandomNL(b *testing.B) {
	rnd := NewRandom()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rnd.readByte()
	}
}
