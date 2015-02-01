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
Package appcontext provides some implementations of types that stores values
shared to local application.

Creating an application context its the recommended way to avoid global
variables and strict the access to your variables to selected functions.

TimedStore

A TimedStore provides values that expires after defined duration of time. That
duration is defined when a new instance is initialized calling
'contextstore.NewTimedStore()' function and it is used to all stored values.

The lifetime of a value can be modified calling 'SetItemDuration()'. The
expiration time of a value is automatically updated when its value is retrieved
by the following methods: 'GetItem()', 'SetItem()' and 'SetItemDuration()'.

TokenStore

A TokenStore provides session tokens to uniquely identify an user session. Each
token expires automatically if it is not used after defined time. The tokens for
authenticated sessions may have a distinct expiration time.
*/
package appcontext
