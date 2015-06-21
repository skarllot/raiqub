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
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	// Defines the maximum data sent by client to 10 MB
	HTTP_BODY_MAX_LENGTH = 1048576
	// WebDAV; RFC 4918
	StatusUnprocessableEntity = 422
)

// JsonWrite sets response content type to JSON, sets HTTP status and serializes
// defined content to JSON format.
func JsonWrite(w http.ResponseWriter, status int, content interface{}) {
	HttpHeader_ContentType_Json().SetWriter(w.Header())
	w.WriteHeader(status)
	if content != nil {
		json.NewEncoder(w).Encode(content)
	}
}

// JsonRead tries to read client sent content using JSON deserialization and
// writes it to defined object.
func JsonRead(body io.ReadCloser, obj interface{}, w http.ResponseWriter) bool {
	content, err := ioutil.ReadAll(io.LimitReader(body, HTTP_BODY_MAX_LENGTH))
	if err != nil {
		jerr := NewJsonErrorFromError(http.StatusInternalServerError, err)
		JsonWrite(w, jerr.Status, jerr)
		return false
	}

	if err := body.Close(); err != nil {
		jerr := NewJsonErrorFromError(http.StatusInternalServerError, err)
		JsonWrite(w, jerr.Status, jerr)
		return false
	}

	if err := json.Unmarshal(content, obj); err != nil {
		jerr := NewJsonErrorFromError(StatusUnprocessableEntity, err)
		JsonWrite(w, jerr.Status, jerr)
		return false
	}

	return true
}
