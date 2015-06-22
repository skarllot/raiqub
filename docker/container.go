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

/*
 * Based on https://github.com/niilo/golib/blob/master/test/dockertest/docker.go
 * Copyright 2014 The Camlistore Authors
 */

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"time"
)

// A Container represents a Docker container.
type Container struct {
	docker *Docker
	id     string
}

// Inspect returns container configuration.
func (s *Container) Inspect() ([]Inspect, error) {
	out, err := exec.Command(s.docker.binCmd, "inspect", s.id).Output()
	if err != nil {
		return nil, err
	}
	
	var list []Inspect
	if err := json.NewDecoder(bytes.NewReader(out)).Decode(&list); err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("Empty output when inspecting container")
	}
	
	return list, nil
}

// Kill terminates current container process.
func (s *Container) Kill() error {
	return exec.Command(s.docker.binCmd, "kill", s.id).Run()
}

// IP returns the network address of current container.
func (s *Container) IP() (string, error) {
	inspect, err := s.Inspect()
	if err != nil {
		return "", err
	}
	
	ip := inspect[0].NetworkSettings.IPAddress
	if ip == "" {
		return "", fmt.Errorf("Could not find an IP. Not running?")
	}
	
	return ip, nil
}

// Remove erases current container from Docker.
func (s *Container) Remove() error {
	return exec.Command(s.docker.binCmd, "rm", s.id).Run()
}

// WaitStartup blocks current thread until current container begins to listen
// exposed port.
func (s *Container) WaitStartup(timeout time.Duration) error {
	inspect, err := s.Inspect()
	if err != nil {
		return fmt.Errorf("Error inspecting container: %v", err)
	}

	port, _ := inspect[0].NetworkSettings.SplitPort(0)
	addr := fmt.Sprintf("%s:%d", inspect[0].NetworkSettings.IPAddress, port)
	return awaitPeer(addr, timeout)
}

func awaitPeer(addr string, timeout time.Duration) error {
	// Based on http://camlistore.org/pkg/netutil/#AwaitReachable
	max := time.Now().Add(timeout)
	for time.Now().Before(max) {
		c, err := net.Dial("tcp", addr)
		if err == nil {
			c.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("%v unreachable for %v", addr, timeout)
}
