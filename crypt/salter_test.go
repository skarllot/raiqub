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
	"testing"
)

func TestSaltUnpredictability(t *testing.T) {
	dict := make(map[string]bool)
	s := NewSalter(NewRandomSourceList(), nil)
	count := 0

	for i := 0; i < DEFAULT_UNPRED_ROUNDS; i++ {
		val := s.DefaultToken()

		if _, ok := dict[val]; ok {
			count++
		} else {
			dict[val] = true
		}
	}

	if count > 0 {
		t.Errorf(
			"Salter class could not generate unpredictable data: %d of %d",
			count, DEFAULT_UNPRED_ROUNDS)
	}
}

func BenchmarkSalter(b *testing.B) {
	salter := NewSalter(NewRandomSourceList(), nil)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		salter.DefaultToken()
	}
}
