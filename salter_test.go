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
	"encoding/base64"
	"testing"
	"time"
)

const (
	DEFAULT_UNPRED_ROUNDS = 10
)

func TestSleepUnpredictability(t *testing.T) {
	now := time.Now()
	for i := 0; i < DEFAULT_UNPRED_ROUNDS; i++ {
		time.Sleep(DEFAULT_SLEEP_TIME)
		diff := time.Now().Sub(now)
		t.Logf("%d: %#v (%d)", i,
			getBytesFromInt64(diff.Nanoseconds(), 2),
			diff.Nanoseconds())
		now = time.Now()
	}
}

func TestTimeRandomUnpredictability(t *testing.T) {
	for i := 0; i < DEFAULT_UNPRED_ROUNDS; i++ {
		t.Logf("%d: %s", i,
			base64.StdEncoding.EncodeToString(getRandomBytes2(time.Now())))
	}
}

func TestSaltUnpredictability(t *testing.T) {
	s := NewSalter()
	for i := 0; i < DEFAULT_UNPRED_ROUNDS; i++ {
		t.Logf("%d: %s", i, s.DefaultToken())
	}
}
