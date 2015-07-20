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
	"net"
	"time"
)

// WaitPeerListening wait for specified peer be ready for new connections.
func WaitPeerListening(network, address string, timeout time.Duration) bool {
	return WaitFunc(100*time.Millisecond, timeout, func() bool {
		c, err := net.Dial(network, address)
		if err == nil {
			c.Close()
			return true
		}

		return false
	})
}
