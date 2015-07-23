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

import "strings"

// A StringSlice represents an array of string.
type StringSlice []string

// IndexOf looks for specified string into current slice.
func (s StringSlice) IndexOf(str string) int {
	for i, v := range s {
		if str == v {
			return i
		}
	}

	return -1
}

// IndexOfIgnoreCase looks for specified string, disregarding letter casing,
// into current slice.
func (s StringSlice) IndexOfIgnoreCase(str string) int {
	str = strings.ToLower(str)
	for i, v := range s {
		if str == strings.ToLower(v) {
			return i
		}
	}

	return -1
}

// Exists determines whether specified string exists into current slice.
func (s StringSlice) Exists(str string) bool {
	return s.IndexOf(str) != -1
}

// ExistsIgnoreCase determines whether specified string exists into current
// slice (ignores letter casing).
func (s StringSlice) ExistsIgnoreCase(str string) bool {
	return s.IndexOfIgnoreCase(str) != -1
}

//ExistsAllIgnoreCase determine whether all specified strings exists into
// current slice (ignores letter casing).
func (s StringSlice) ExistsAllIgnoreCase(str []string) bool {
	for _, v := range str {
		if !s.ExistsIgnoreCase(v) {
			return false
		}
	}

	return true
}
