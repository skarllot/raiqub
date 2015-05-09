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

// A HttpMiddlewareFunc represents a HTTP handler that is called before final
// HTTP handler.
type HttpMiddlewareFunc func(http.Handler) http.Handler

// A Chain represents a chain of HTTP middlewares that are called before target
// HTTP handler.
type Chain []HttpMiddlewareFunc

// Get returns a HTTP handler which is a chain of middlewares and then the
// requested handler.
func (s Chain) Get(handler http.Handler) http.Handler {
	for i := len(s) - 1; i >= 0; i-- {
		handler = s[i](handler)
	}
	return handler
}
