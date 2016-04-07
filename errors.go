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
	"fmt"
)

// A DuplicatedKeyError represents an error when the provided new key to insert
// to an dictionary already exists.
type DuplicatedKeyError string

// Error returns string representation of current instance error.
func (e DuplicatedKeyError) Error() string {
	return fmt.Sprintf(
		"Could not create the '%s' key because it already exists", string(e))
}

// A InvalidKeyError represents an error when the provided key to retrieve of
// and dictionary does not exist.
type InvalidKeyError string

// Error returns string representation of current instance error.
func (e InvalidKeyError) Error() string {
	return fmt.Sprintf(
		"Could not get the '%s' key because it does not exist or it is expired",
		string(e))
}
