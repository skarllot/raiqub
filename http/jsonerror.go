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
	"reflect"
)

// A JsonError represents an error returned by JSON-based API.
type JsonError struct {
	// HTTP status code.
	Status int `json:"status,omitempty"`
	// Error code.
	Code int `json:"code,omitempty"`
	// Error type.
	Type string `json:"type,omitempty"`
	// A message with error details.
	Message string `json:"message,omitempty"`
	// A URL for reference.
	MoreInfo string `json:"moreInfo,omitempty"`
}

func NewJsonErrorFromError(status int, e error) JsonError {
	errType := reflect.TypeOf(e)
	var typeName string
	if errType.Kind() == reflect.Ptr {
		typeName = errType.Elem().Name()
	} else {
		typeName = errType.Name()
	}

	return JsonError{
		Status:  status,
		Type:    typeName,
		Message: e.Error(),
	}
}
