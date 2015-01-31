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
package contextstore

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type TimedStore struct {
	items    map[string]*TimedItem
	duration time.Duration
	mutex    *sync.Mutex
}

func (s *TimedStore) New(d time.Duration) *TimedStore {
	return &TimedStore{
		items:    make(map[string]*TimedItem),
		duration: d,
		mutex:    &sync.Mutex{},
	}
}

func (s *TimedStore) NewItem(id string, value interface{}) *TimedItem {
	s.removeExpired()
	i := &TimedItem{
		expireAt: time.Now().Add(s.duration),
		duration: 0,
		value:    value,
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.items[id] = i

	return i
}

func (s *TimedStore) GetItem(id string) (interface{}, error) {
	s.removeExpired()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	v, err := s.unsafeGet(id)
	if err != nil {
		return nil, err
	}
	s.refreshItem(v)
	return v.value, nil
}

func (s *TimedStore) refreshItem(i *TimedItem) {
	d := s.duration
	if i.duration != 0 {
		d = i.duration
	}
	i.expireAt = time.Now().Add(d)
}

func (s *TimedStore) removeExpired() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i := range s.items {
		if s.items[i].IsExpired() {
			delete(s.items, i)
		}
	}
}

func (s *TimedStore) RemoveItem(id string) error {
	s.removeExpired()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, err := s.unsafeGet(id)
	if err != nil {
		return err
	}

	delete(s.items, id)
	return nil
}

func (s *TimedStore) SetItem(id string, value interface{}) error {
	s.removeExpired()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	v, err := s.unsafeGet(id)
	if err != nil {
		return err
	}

	s.refreshItem(v)
	v.value = value
	return nil
}

func (s *TimedStore) SetItemDuration(id string, d time.Duration) error {
	s.removeExpired()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	v, err := s.unsafeGet(id)
	if err != nil {
		return err
	}

	v.expireAt = time.Now().Add(d)
	v.duration = d
	return nil
}

func (s *TimedStore) unsafeGet(id string) (*TimedItem, error) {
	v, ok := s.items[id]
	if !ok {
		return nil, errors.New(
			fmt.Sprintf("The requested id '%s' does not exist or is expired", id))
	}
	return v, nil
}
