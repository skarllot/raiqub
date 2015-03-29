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
	"errors"
	"fmt"
	"sync"
	"time"
)

// A Cache provides in-memory key:value cache that expires after defined
// duration of time.
type Cache struct {
	values   map[string]*cacheItem
	lifetime time.Duration
	mutex    *sync.Mutex
}

// NewCache creates a new instance of Cache.
func NewCache(d time.Duration) *Cache {
	return &Cache{
		values:   make(map[string]*cacheItem),
		lifetime: d,
		mutex:    &sync.Mutex{},
	}
}

// Add adds a new key:value to current Cache instance.
func (s *Cache) Add(key string, value interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.removeExpired()

	if _, ok := s.values[key]; ok {
		return errors.New(
			"Could not allocate the new value because of duplicated id")
	}

	i := &cacheItem{
		expireAt: time.Now().Add(s.lifetime),
		lifetime: s.lifetime,
		value:    value,
	}

	s.values[key] = i
	return nil
}

// Count gets the number of cached values by current instance.
func (s *Cache) Count() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.removeExpired()

	return len(s.values)
}

// Get gets the value cached by specified key.
func (s *Cache) Get(key string) (interface{}, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.removeExpired()

	v, err := s.unsafeGet(key)
	if err != nil {
		return nil, err
	}
	v.Postpone()
	return v.value, nil
}

// removeExpired remove all expired values from current Cache instance list.
func (s *Cache) removeExpired() {
	for i := range s.values {
		if s.values[i].IsExpired() {
			delete(s.values, i)
		}
	}
}

// Delete deletes the specified key:value.
func (s *Cache) Delete(key string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.removeExpired()

	_, err := s.unsafeGet(key)
	if err != nil {
		return err
	}

	delete(s.values, key)
	return nil
}

// Set sets the value of specified key.
func (s *Cache) Set(key string, value interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.removeExpired()

	v, err := s.unsafeGet(key)
	if err != nil {
		return err
	}

	v.Postpone()
	v.value = value
	return nil
}

// SetLifetime modifies the lifetime of specified key:value.
func (s *Cache) SetLifetime(key string, d time.Duration) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.removeExpired()

	v, err := s.unsafeGet(key)
	if err != nil {
		return err
	}

	v.lifetime = d
	v.Postpone()
	return nil
}

// unsafeGet gets one cacheItem instance from its key without locking.
func (s *Cache) unsafeGet(key string) (*cacheItem, error) {
	v, ok := s.values[key]
	if !ok {
		return nil, errors.New(
			fmt.Sprintf("The requested id '%s' does not exist or is expired", key))
	}
	return v, nil
}
