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
	"strconv"
	"strings"
	"time"
)

// A Inspect represents a return value from Docker inspect command.
type Inspect struct {
	Created         time.Time
	Driver          string
	ExecDriver      string
	HostnamePath    string
	HostsPath       string
	Id              string
	Image           string
	Name            string
	NetworkSettings NetworkSettings
	Path            string
	RestartCount    int
	State           State
}

// A NetworkSettings represents the network settings returned from inspect
// command.
type NetworkSettings struct {
	Bridge      string
	Gateway     string
	IPAddress   string
	IPPrefixLen int
	MacAddress  string
	Ports       map[string]*HostPublish
}

// A HostPublish represents publishing configuration for exposed port.
type HostPublish struct {
	HostIp   string
	HostPort string
}

// A State represents current container state
type State struct {
	Error      string
	ExitCode   int
	FinishedAt time.Time
	OOMKilled  bool
	Paused     bool
	Pid        int
	Restarting bool
	Running    bool
	StartedAt  time.Time
}

// PortAt returns expose and publish port information at especified index.
func (s NetworkSettings) PortAt(index int) (string, *HostPublish) {
	counter := 0
	for k, v := range s.Ports {
		if counter != index {
			counter++
			continue
		}

		return k, v
	}

	return "", nil
}

// SplitPort splits port number and protocol name from Ports string.
func (s NetworkSettings) SplitPort(index int) (uint16, string) {
	ptProt, _ := s.PortAt(index)
	vect := strings.Split(ptProt, "/")
	port, _ := strconv.ParseUint(vect[0], 10, 16)
	return uint16(port), vect[1]
}
