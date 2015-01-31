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
package contextstore

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type TokenStore struct {
	*TimedStore
	salt         string
	authDuration time.Duration
}

func (s *TokenStore) Count() int {
	return len(s.items)
}

func (s *TokenStore) getInvalidTokenError(token string) error {
	return errors.New(fmt.Sprintf(
		"The requested token '%s' is invalid or is expired", token))
}

func (s *TokenStore) GetValue(token string) (interface{}, error) {
	v, err := s.GetItem(token)
	if err != nil {
		return nil, s.getInvalidTokenError(token)
	}
	return v, err
}

func (s *TokenStore) New(noAuth, auth time.Duration, salt string) *TokenStore {
	ts := (&TimedStore{}).New(noAuth)
	return &TokenStore{
		ts,
		salt,
		auth,
	}
}

func (s *TokenStore) NewToken() string {
	hash := sha256.New()
	now := time.Now().Format(time.ANSIC)

	hash.Write([]byte(now))
	hash.Write([]byte(strconv.Itoa(time.Now().Nanosecond())))
	hash.Write([]byte(s.salt))
	hash.Write(getRandomBytes(32 + time.Now().Second()))
	strSum := base64.URLEncoding.EncodeToString(hash.Sum(nil))

	s.NewItem(strSum, nil)

	return strSum
}

func (s *TokenStore) RemoveToken(token string) error {
	err := s.RemoveItem(token)
	if err != nil {
		return s.getInvalidTokenError(token)
	}
	return nil
}

func (s *TokenStore) SetTokenAsAuthenticated(token string) error {
	err := s.SetItemDuration(token, s.authDuration)
	if err != nil {
		return s.getInvalidTokenError(token)
	}
	return nil
}

func (s *TokenStore) SetValue(token string, value interface{}) error {
	err := s.SetItem(token, value)
	if err != nil {
		return s.getInvalidTokenError(token)
	}
	return nil
}

func getRandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic("Could not access secure random generator")
	}

	return b
}

func getRandomString(n int) string {
	b := getRandomBytes(n)
	return base64.URLEncoding.EncodeToString(b)
}
