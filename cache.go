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

package raiqub

import (
	"sync"
	"time"
)

// A Cache provides in-memory key:value cache that expires after defined
// duration of time.
type Cache struct {
	values   map[string]*cacheItem
	lifetime time.Duration
	sync.RWMutex
}

// NewCache creates a new instance of Cache and defines the default lifetime for
// new cached items.
func NewCache(d time.Duration) *Cache {
	return &Cache{
		values:   make(map[string]*cacheItem),
		lifetime: d,
	}
}

// Add adds a new key:value to current Cache instance.
//
// Errors:
// - DuplicatedKeyError when requested key already exists.
func (s *Cache) Add(key string, value interface{}) error {
	lckStatus := s.removeExpired()

	i := &cacheItem{
		expireAt: time.Now().Add(s.lifetime),
		lifetime: s.lifetime,
		value:    value,
	}

	if lckStatus == ReadLocked {
		s.RUnlock()
		s.Lock()
	}
	defer s.Unlock()

	if _, ok := s.values[key]; ok {
		return DuplicatedKeyError(key)
	}

	s.values[key] = i
	return nil
}

// Count gets the number of cached values by current instance.
func (s *Cache) Count() int {
	if s.removeExpired() == WriteLocked {
		defer s.Unlock()
	} else {
		defer s.RUnlock()
	}

	return len(s.values)
}

// Flush deletes any cached value into current instance.
func (s *Cache) Flush() {
	s.Lock()
	defer s.Unlock()

	s.values = make(map[string]*cacheItem)
}

// Get gets the value cached by specified key.
//
// Errors:
// - InvalidKeyError when requested key could not be found.
func (s *Cache) Get(key string) (interface{}, error) {
	if s.removeExpired() == WriteLocked {
		s.Unlock()
		s.RLock()
	}
	defer s.RUnlock()

	v, err := s.unsafeGet(key)
	if err != nil {
		return nil, err
	}
	v.Postpone()
	return v.value, nil
}

// Delete deletes the specified key:value.
//
// Errors:
// - InvalidKeyError when requested key could not be found.
func (s *Cache) Delete(key string) error {
	lckStatus := s.removeExpired()

	_, err := s.unsafeGet(key)
	if err != nil {
		return err
	}

	if lckStatus == ReadLocked {
		s.RUnlock()
		s.Lock()
	}
	defer s.Unlock()

	delete(s.values, key)
	return nil
}

// Set sets the value of specified key.
//
// Errors:
// - InvalidKeyError when requested key could not be found.
func (s *Cache) Set(key string, value interface{}) error {
	if s.removeExpired() == WriteLocked {
		s.Unlock()
		s.RLock()
	}
	defer s.RUnlock()

	v, err := s.unsafeGet(key)
	if err != nil {
		return err
	}

	v.Postpone()
	v.value = value
	return nil
}

// SetLifetime modifies the lifetime of specified key:value.
//
// Errors:
// - InvalidKeyError when requested key could not be found.
func (s *Cache) SetLifetime(key string, d time.Duration) error {
	if s.removeExpired() == WriteLocked {
		s.Unlock()
		s.RLock()
	}
	defer s.RUnlock()

	v, err := s.unsafeGet(key)
	if err != nil {
		return err
	}

	v.lifetime = d
	v.Postpone()
	return nil
}

// removeExpired remove all expired values from current Cache instance list.
//
// Returns the locking status of current instance.
func (s *Cache) removeExpired() LockStatus {
	writeLocked := false
	s.RLock()
	for i := range s.values {
		if s.values[i].IsExpired() {
			if !writeLocked {
				s.RUnlock()
				s.Lock()
				writeLocked = true
			}
			delete(s.values, i)
		}
	}

	if writeLocked {
		return WriteLocked
	} else {
		return ReadLocked
	}
}

// unsafeGet gets one cacheItem instance from its key without locking.
//
// Errors:
// - InvalidKeyError when requested key could not be found.
func (s *Cache) unsafeGet(key string) (*cacheItem, error) {
	v, ok := s.values[key]
	if !ok {
		return nil, InvalidKeyError(key)
	}
	return v, nil
}
