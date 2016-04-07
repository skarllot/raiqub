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

/*
Package http provides operations to help HTTP server implementation.

Chain

A Chain provides a function to chain HTTP handlers, also know as middlewares,
before a specified HTTP handler. A Chain is basically a slice of middlewares.

HttpHeader

A HttpHeader provides functions to help handle HTTP headers, both for reading
from client request and write to server response. Additionally provides some
pre-defined headers.

Route

A Route provides a easy way to define HTTP routes and handle it.

SessionCache

A SessionCache provides session tokens to uniquely identify an user session and
links it to specified data. Each token expires automatically if it is not used
after defined time.
*/
package http
