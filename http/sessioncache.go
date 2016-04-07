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

package http

import (
	"errors"
	"fmt"
	"github.com/skarllot/raiqub/crypt"
	"github.com/skarllot/raiqub/data"
	"time"
)

// A SessionCache provides a temporary token to uniquely identify an user
// session.
type SessionCache struct {
	cache  *data.Cache
	salter *crypt.Salter
}

// NewSessionCache creates a new instance of SessionCache and defines a lifetime
// for sessions and a initial salt for random input.
func NewSessionCache(d time.Duration, salt string) *SessionCache {
	return &SessionCache{
		cache: data.NewCache(d),
		salter: crypt.NewSalter(
			crypt.NewRandomSourceListSecure(), []byte(salt)),
	}
}

// Count gets the number of tokens stored by current instance.
func (s *SessionCache) Count() int {
	return s.cache.Count()
}

// getInvalidTokenError gets the default error when an invalid or expired token
// is requested.
func (s *SessionCache) getInvalidTokenError(token string) error {
	return errors.New(fmt.Sprintf(
		"The requested token '%s' is invalid or is expired", token))
}

// Get gets the value stored by specified token.
func (s *SessionCache) Get(token string) (interface{}, error) {
	v, err := s.cache.Get(token)
	if err != nil {
		return nil, s.getInvalidTokenError(token)
	}
	return v, err
}

// Add creates a new unique token and stores it into current SessionCache
// instance.
//
// The token creation will take at least 200 microseconds, but could normally
// take 2.5 milliseconds. The token generation function it is built with
// security over performance.
func (s *SessionCache) Add() string {
	strSum := s.salter.DefaultToken()

	err := s.cache.Add(strSum, nil)
	if err != nil {
		panic("Something is seriously wrong, a duplicated token was generated")
	}

	return strSum
}

// Delete deletes specified token from current SessionCache instance.
func (s *SessionCache) Delete(token string) error {
	err := s.cache.Delete(token)
	if err != nil {
		return s.getInvalidTokenError(token)
	}
	return nil
}

// Set store a value to specified token.
func (s *SessionCache) Set(token string, value interface{}) error {
	err := s.cache.Set(token, value)
	if err != nil {
		return s.getInvalidTokenError(token)
	}
	return nil
}
