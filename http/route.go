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

// A Route represents a HTTP route.
type Route struct {
	// Unique name to current route.
	Name string
	// HTTP method handled by this route.
	Method string
	// HTTP path handled by this route.
	Path string
	// Indicates whether authentication is required to call this route.
	MustAuth bool
	// Defines which method is called to handle this route.
	ActionFunc http.HandlerFunc
}

// Routes represents a slice of Route objects.
type Routes []Route

// The Routable interface is implemented by objects that has HTTP routes.
type Routable interface {
	Routes() Routes
}

// MergeRoutes returns a slice with all routes returned by Routable objects.
func MergeRoutes(r ...Routable) Routes {
	routes := make(Routes, 0)
	for _, v := range r {
		routes = append(routes, v.Routes()...)
	}
	return routes
}
