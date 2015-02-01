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

package appcontext

import (
	"testing"
	"time"
)

const TOKEN_SALT = "CvoTVwDw685Ve0qjGn//zmHGKvoCcslYNQT4AQ9FygSk9t6NuzBHuohyOHhqb/1omn6c"

func TestTokenExpiration(t *testing.T) {
	ts := NewTokenStore(time.Millisecond*10, time.Millisecond*30, TOKEN_SALT)

	t1 := ts.NewToken()
	t2 := ts.NewToken()

	if _, err := ts.GetValue(t1); err != nil {
		t.Error("The value t1 was not stored")
	}
	if _, err := ts.GetValue(t2); err != nil {
		t.Error("The value t2 was not stored")
	}

	time.Sleep(time.Millisecond * 20)

	if _, err := ts.GetValue(t1); err == nil {
		t.Error("The value t1 was not expired")
	}
	if _, err := ts.GetValue(t2); err == nil {
		t.Error("The value t2 was not expired")
	}
}

func BenchmarkTokenCreation(b *testing.B) {
	ts := NewTokenStore(time.Millisecond, time.Second, TOKEN_SALT)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ts.NewToken()
	}
}
