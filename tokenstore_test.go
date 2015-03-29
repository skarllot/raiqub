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
	"testing"
	"time"
)

const TOKEN_SALT = "CvoTVwDw685Ve0qjGn//zmHGKvoCcslYNQT4AQ9FygSk9t6NuzBHuohyOHhqb/1omn6c"

func TestTokenAuthentication(t *testing.T) {
	ts := NewTokenStore(time.Millisecond*10, time.Second, TOKEN_SALT)

	t1 := ts.NewToken()
	if err := ts.SetTokenAsAuthenticated(t1); err != nil {
		t.Error("The token t1 could not authenticate")
	}

	time.Sleep(time.Millisecond * 20)

	if _, err := ts.GetValue(t1); err != nil {
		t.Error("The token t1 expired sooner than expected")
	}
}
func TestTokenExpiration(t *testing.T) {
	ts := NewTokenStore(time.Millisecond*10, time.Millisecond*30, TOKEN_SALT)

	t1 := ts.NewToken()
	t2 := ts.NewToken()

	if _, err := ts.GetValue(t1); err != nil {
		t.Error("The token t1 was not stored")
	}
	if _, err := ts.GetValue(t2); err != nil {
		t.Error("The token t2 was not stored")
	}

	time.Sleep(time.Millisecond * 20)

	if _, err := ts.GetValue(t1); err == nil {
		t.Error("The token t1 was not expired")
	}
	if _, err := ts.GetValue(t2); err == nil {
		t.Error("The token t2 was not expired")
	}

	if err := ts.RemoveToken(t1); err == nil {
		t.Error("The expired token t1 should not be removable")
	}
	if err := ts.SetValue(t2, nil); err == nil {
		t.Error("The expired token t2 should not be changeable")
	}
	if err := ts.SetTokenAsAuthenticated(t1); err == nil {
		t.Error("The expired token t1 should not be authenticable")
	}
}

func TestTokenHandling(t *testing.T) {
	testValues := []struct {
		ref   string
		token string
		value int
	}{
		{"t1", "", 3},
		{"t2", "", 6},
		{"t3", "", 83679},
		{"t4", "", 2748},
		{"t5", "", 54},
		{"t6", "", 6},
		{"t7", "", 2},
		{"t8", "", 8},
		{"t9", "", 7},
		{"t10", "", 8},
	}
	rmTestIndex := 6
	changeValues := map[int]int{
		2: 5062,
		9: 4099,
	}

	ts := NewTokenStore(time.Millisecond*10, time.Millisecond*30, TOKEN_SALT)

	for i, _ := range testValues {
		item := &testValues[i]
		item.token = ts.NewToken()
		err := ts.SetValue(item.token, item.value)
		if err != nil {
			t.Errorf("The token %s could not be set", item.ref)
		}
	}

	if ts.Count() != len(testValues) {
		t.Error("The tokens count do not match")
	}

	for _, i := range testValues {
		v, err := ts.GetValue(i.token)
		if err != nil {
			t.Errorf("The token %s could not be read", i.ref)
		}
		if v != i.value {
			t.Errorf("The token %s was stored incorrectly", i.ref)
		}
	}

	rmTestKey := testValues[rmTestIndex]
	if err := ts.RemoveToken(rmTestKey.token); err != nil {
		t.Errorf("The token %s could not be removed", rmTestKey.ref)
	}
	if _, err := ts.GetValue(rmTestKey.token); err == nil {
		t.Errorf("The removed token %s should not be retrieved", rmTestKey.ref)
	}
	if ts.Count() == len(testValues) {
		t.Error("The tokens count should not match")
	}

	for k, v := range changeValues {
		item := testValues[k]
		err := ts.SetValue(item.token, v)
		if err != nil {
			t.Errorf("The token %s could not be changed", item.ref)
		}
	}
	for k, v := range changeValues {
		item := testValues[k]
		v2, err := ts.GetValue(item.token)
		if err != nil {
			t.Errorf("The token %s could not be read", item.ref)
		}
		if v2 != v {
			t.Errorf("The token %s was not changed", item.ref)
		}
	}
}

func BenchmarkTokenCreation(b *testing.B) {
	ts := NewTokenStore(time.Millisecond, time.Second, TOKEN_SALT)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ts.NewToken()
	}
}
