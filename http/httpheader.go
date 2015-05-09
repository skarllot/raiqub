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

package http

import (
	"net/http"
)

// A HttpHeader represents a key-value pair in a HTTP header.
type HttpHeader struct {
	// HTTP header field name.
	Name string
	// HTTP header field value.
	Value string
}

// Clone make a copy of current instance.
func (s HttpHeader) Clone() *HttpHeader {
	return &s
}

// GetReader gets HTTP header value, as defined by current instance, from
// Request Header and sets to current instance.
func (s *HttpHeader) GetReader(h http.Header) *HttpHeader {
	s.Value = h.Get(s.Name)
	return s
}

// SetName sets header name of current instance.
func (s *HttpHeader) SetName(name string) *HttpHeader {
	s.Name = name
	return s
}

// SetValue sets header value of current instance.
func (s *HttpHeader) SetValue(value string) *HttpHeader {
	s.Value = value
	return s
}

// SetWriter sets HTTP header, as defined by current instance, to ResponseWriter
// Header.
func (s *HttpHeader) SetWriter(h http.Header) *HttpHeader {
	h.Set(s.Name, s.Value)
	return s
}
