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

func TestValueExpiration(t *testing.T) {
	ts := NewTimedStore(time.Millisecond * 10)

	ts.NewValue("v1", nil)
	ts.NewValue("v2", nil)

	if _, err := ts.GetValue("v1"); err != nil {
		t.Error("The value v1 was not stored")
	}
	if _, err := ts.GetValue("v2"); err != nil {
		t.Error("The value v2 was not stored")
	}

	time.Sleep(time.Millisecond * 20)

	if _, err := ts.GetValue("v1"); err == nil {
		t.Error("The value v1 was not expired")
	}
	if _, err := ts.GetValue("v2"); err == nil {
		t.Error("The value v2 was not expired")
	}
}

func TestValueIdCollision(t *testing.T) {
	ts := NewTimedStore(time.Millisecond)

	if _, err := ts.NewValue("v1", nil); err != nil {
		t.Error("The value v1 could not be stored")
	}
	if _, err := ts.NewValue("v1", nil); err == nil {
		t.Error("The duplicated v1 could be stored")
	}
}

func TestValueSetExpiration(t *testing.T) {
	ts := NewTimedStore(time.Millisecond)

	ts.NewValue("v1", nil)
	ts.SetValueDuration("v1", time.Second)

	time.Sleep(time.Millisecond * 10)

	if _, err := ts.GetValue("v1"); err != nil {
		t.Error("The value v1 is expired before expected")
	}
}

func BenchmarkValueCreation(b *testing.B) {
	ts := NewTimedStore(time.Millisecond)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ts.NewValue(time.Now().Format(time.ANSIC), time.Now())
	}
}
