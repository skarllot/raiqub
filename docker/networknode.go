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

package docker

import (
	"fmt"
	"strconv"
	"strings"
)

// A NetworkNode represents a network node from Docker container.
type NetworkNode struct {
	Protocol  string
	IpAddress string
	Port      uint16
}

// FormatDialAddress returns an address as expected by net.Dial.
func (s NetworkNode) FormatDialAddress() string {
	return fmt.Sprintf("%s:%d", s.IpAddress, s.Port)
}

// SetFromDocker parses port as provided by Docker and extract protocol and port
// number.
func (s *NetworkNode) SetFromDocker(port string) error {
	split := strings.Split(port, "/")
	iPort, err := strconv.ParseUint(split[0], 10, 16)
	if err != nil {
		return err
	}

	s.Port, s.Protocol = uint16(iPort), split[1]
	return nil
}
