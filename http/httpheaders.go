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

// HttpHeader_AccessControlAllowCredentials creates a HTTP header to CORS-able
// API indicate that authentication is allowed.
func HttpHeader_AccessControlAllowCredentials() *HttpHeader {
	return &HttpHeader{
		"Access-Control-Allow-Credentials",
		"", // boolean
	}
}

// HttpHeader_AccessControlAllowCredentials creates a HTTP header to CORS-able
// API indicate which HTTP headers are allowed.
func HttpHeader_AccessControlAllowHeaders() *HttpHeader {
	return &HttpHeader{
		"Access-Control-Allow-Headers",
		"", // comma-separated list
	}
}

// HttpHeader_AccessControlAllowMethods creates a HTTP header to CORS-able API
// indicate which HTTP methods are allowed to current resource.
func HttpHeader_AccessControlAllowMethods() *HttpHeader {
	return &HttpHeader{
		"Access-Control-Allow-Methods",
		"", // comma-separated list of HTTP methods
	}
}

// HttpHeader_AccessControlAllowOrigin creates a HTTP header to CORS-able API
// indicate which origin is expected.
func HttpHeader_AccessControlAllowOrigin() *HttpHeader {
	return &HttpHeader{
		"Access-Control-Allow-Origin",
		"", // http-formatted domain or asterisk to any
	}
}

// HttpHeader_AccessControlMaxAge creates a HTTP header to CORS-able API
// indicate how long preflight results should be cached.
func HttpHeader_AccessControlMaxAge() *HttpHeader {
	return &HttpHeader{
		"Access-Control-Max-Age",
		"", // seconds
	}
}

// HttpHeader_AccessControlRequestHeaders creates a HTTP header to CORS-able
// client indicate which headers will be used for request.
func HttpHeader_AccessControlRequestHeaders() *HttpHeader {
	return &HttpHeader{
		"Access-Control-Request-Headers",
		"", // comma-separated list of HTTP headers
	}
}

// HttpHeader_AccessControlRequestMethod creates a HTTP header to CORS-able
// client indicate which HTTP method will be used for request.
func HttpHeader_AccessControlRequestMethod() *HttpHeader {
	return &HttpHeader{
		"Access-Control-Request-Method",
		"", // HTTP method name
	}
}

// HttpHeader_ContentType_Json creates a HTTP header to define JSON content
// type.
func HttpHeader_ContentType_Json() *HttpHeader {
	return &HttpHeader{
		"Content-Type",
		"application/json; charset=UTF-8",
	}
}

// HttpHeader_Location creates a HTTP header to define location of new object.
func HttpHeader_Location() *HttpHeader {
	return &HttpHeader{
		"Location",
		"", // relative http location
	}
}

// HttpHeader_Origin creates a HTTP header to CORS-able client indicate its
// address.
func HttpHeader_Origin() *HttpHeader {
	return &HttpHeader{
		"Origin",
		"", // http-formatted domain
	}
}
