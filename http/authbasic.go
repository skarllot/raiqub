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
	"bytes"
	"encoding/base64"
	"net/http"
	"strings"
)

const (
	BASIC_PREFIX = "Basic "
	BASIC_REALM  = BASIC_PREFIX + "realm=\"Restricted\""
)

// A HttpBasicAuthenticator represents a handler for HTTP basic authentication.
type HttpBasicAuthenticator struct {
	HttpAuthenticable
}

// AuthHandler is a HTTP request middleware that enforces authentication.
func (self HttpBasicAuthenticator) AuthHandler(next http.Handler) http.Handler {
	if self.HttpAuthenticable == nil {
		panic("HttpAuthenticable cannot be nil")
	}

	f := func(w http.ResponseWriter, r *http.Request) {
		user, secret := self.parseAuthHeader(r.Header.Get("Authorization"))
		if len(user) > 0 &&
			len(secret) > 0 &&
			self.TryAuthentication(r, user, secret) {
			next.ServeHTTP(w, r)
			return
		}

		HttpHeader_WwwAuthenticate().
			SetValue(BASIC_REALM).
			SetWriter(w.Header())
		http.Error(w, http.StatusText(http.StatusUnauthorized),
			http.StatusUnauthorized)
	}

	return http.HandlerFunc(f)
}

func (self HttpBasicAuthenticator) parseAuthHeader(
	header string,
) (user, secret string) {
	if !strings.HasPrefix(header, BASIC_PREFIX) {
		return
	}
	payload, err := base64.StdEncoding.DecodeString(header[len(BASIC_PREFIX):])
	if err != nil {
		return
	}
	pair := bytes.SplitN(payload, []byte(":"), 2)
	if len(pair) != 2 {
		return
	}

	user, secret = string(pair[0]), string(pair[1])
	user, secret = strings.TrimSpace(user), strings.TrimSpace(secret)
	return
}
