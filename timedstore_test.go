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

	ts.AddValue("v1", nil)
	ts.AddValue("v2", nil)

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

	if err := ts.RemoveValue("v1"); err == nil {
		t.Error("The expired value v1 should not be removable")
	}
	if err := ts.SetValue("v2", nil); err == nil {
		t.Error("The expired value v2 should not be changeable")
	}
}

func TestValueHandling(t *testing.T) {
	testValues := map[string]int{
		"v1":  3,
		"v2":  6,
		"v3":  83679,
		"v4":  2748,
		"v5":  54,
		"v6":  6,
		"v7":  2,
		"v8":  8,
		"v9":  7,
		"v10": 8,
	}
	rmTestKey := "v5"
	changeValues := map[string]int{
		"v4": 5062,
		"v9": 4099,
	}

	ts := NewTimedStore(time.Millisecond * 10)

	for k, v := range testValues {
		_, err := ts.AddValue(k, v)
		if err != nil {
			t.Errorf("The value %s could not be added", k)
		}
	}

	if ts.Count() != len(testValues) {
		t.Error("The values count do not match")
	}

	for k, v := range testValues {
		v2, err := ts.GetValue(k)
		if err != nil {
			t.Errorf("The value %s could not be read", k)
		}
		if v2 != v {
			t.Errorf("The value %s was stored incorrectly", k)
		}
	}

	if err := ts.RemoveValue(rmTestKey); err != nil {
		t.Errorf("The value %s could not be removed", rmTestKey)
	}
	if _, err := ts.GetValue(rmTestKey); err == nil {
		t.Errorf("The removed value %s should not be retrieved", rmTestKey)
	}
	if ts.Count() == len(testValues) {
		t.Error("The values count should not match")
	}

	for k, v := range changeValues {
		err := ts.SetValue(k, v)
		if err != nil {
			t.Errorf("The value %s could not be changed", k)
		}
	}
	for k, v := range changeValues {
		v2, err := ts.GetValue(k)
		if err != nil {
			t.Errorf("The value %s could not be read", k)
		}
		if v2 != v {
			t.Errorf("The value %s was not changed", k)
		}
	}
}

func TestValueIdCollision(t *testing.T) {
	ts := NewTimedStore(time.Millisecond)

	if _, err := ts.AddValue("v1", nil); err != nil {
		t.Error("The value v1 could not be stored")
	}
	if _, err := ts.AddValue("v1", nil); err == nil {
		t.Error("The duplicated v1 could be stored")
	}
}

func TestValueSetExpiration(t *testing.T) {
	ts := NewTimedStore(time.Millisecond)

	ts.AddValue("v1", nil)
	ts.SetValueDuration("v1", time.Second)

	time.Sleep(time.Millisecond * 10)

	if _, err := ts.GetValue("v1"); err != nil {
		t.Error("The value v1 is expired before expected")
	}

	if err := ts.SetValueDuration("v2", time.Second); err == nil {
		t.Error("Should not be possible to set duration for invalid value")
	}
}

func BenchmarkValueCreation(b *testing.B) {
	ts := NewTimedStore(time.Millisecond)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ts.AddValue(time.Now().Format(time.RFC3339Nano), time.Now())
	}
}
