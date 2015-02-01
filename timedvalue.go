/*
 * Copyright 2015 Fabrício Godoy
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package appcontext

import (
	"time"
)

// A TimedValue represents a value that expires after defined time.
type TimedValue struct {
	expireAt time.Time
	duration time.Duration
	value    interface{}
}

// IsExpired returns whether current value is expired.
func (i *TimedValue) IsExpired() bool {
	return time.Now().After(i.expireAt)
}

// UpdateExpiration postpone value expiration time to current time added to its
// lifetime duration.
func (i *TimedValue) UpdateExpiration() {
	i.expireAt = time.Now().Add(i.duration)
}
