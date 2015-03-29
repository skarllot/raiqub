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

// A TimedStore provides values that expires after defined duration of time.
type TimedStore struct {
	values   map[string]*TimedValue
	duration time.Duration
	mutex    *sync.Mutex
}

// NewTimedStore creates a new instance of TimedStore.
func NewTimedStore(d time.Duration) *TimedStore {
	return &TimedStore{
		values:   make(map[string]*TimedValue),
		duration: d,
		mutex:    &sync.Mutex{},
	}
}

// AddValue adds a new value to current TimedStore instance.
func (s *TimedStore) AddValue(id string, value interface{}) (*TimedValue, error) {
	s.removeExpired()

	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, ok := s.values[id]; ok {
		return nil, errors.New(
			"Could not allocate the new value because of duplicated id")
	}

	i := &TimedValue{
		expireAt: time.Now().Add(s.duration),
		duration: s.duration,
		value:    value,
	}

	s.values[id] = i
	return i, nil
}

// Count gets the number of stored values by current instance.
func (s *TimedStore) Count() int {
	s.removeExpired()

	s.mutex.Lock()
	defer s.mutex.Unlock()
	return len(s.values)
}

// GetValue gets the value stored by specified ID.
func (s *TimedStore) GetValue(id string) (interface{}, error) {
	s.removeExpired()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	v, err := s.unsafeGet(id)
	if err != nil {
		return nil, err
	}
	v.UpdateExpiration()
	return v.value, nil
}

// removeExpired remove all expired values from current TimedStore instance
// list.
func (s *TimedStore) removeExpired() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i := range s.values {
		if s.values[i].IsExpired() {
			delete(s.values, i)
		}
	}
}

// RemoveValue removes the value of specified ID.
func (s *TimedStore) RemoveValue(id string) error {
	s.removeExpired()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, err := s.unsafeGet(id)
	if err != nil {
		return err
	}

	delete(s.values, id)
	return nil
}

// SetValue sets the value of specified ID.
func (s *TimedStore) SetValue(id string, value interface{}) error {
	s.removeExpired()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	v, err := s.unsafeGet(id)
	if err != nil {
		return err
	}

	v.UpdateExpiration()
	v.value = value
	return nil
}

// SetValueDuration modifies the lifetime of specified value.
func (s *TimedStore) SetValueDuration(id string, d time.Duration) error {
	s.removeExpired()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	v, err := s.unsafeGet(id)
	if err != nil {
		return err
	}

	v.duration = d
	v.UpdateExpiration()
	return nil
}

// unsafeGet gets one TimedValue instance from its ID without locking.
func (s *TimedStore) unsafeGet(id string) (*TimedValue, error) {
	v, ok := s.values[id]
	if !ok {
		return nil, errors.New(
			fmt.Sprintf("The requested id '%s' does not exist or is expired", id))
	}
	return v, nil
}
