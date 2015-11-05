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
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/raiqub/dot"
)

// A Container represents a Docker container.
type Container struct {
	docker *Docker
	id     string
}

// HasExposedPorts returns whether current container has exposed ports.
func (s *Container) HasExposedPorts() bool {
	nodes, err := s.NetworkNodes()
	if err != nil {
		return false
	}

	for _, v := range nodes {
		if v.Port > 0 || len(v.Protocol) > 0 {
			return true
		}
	}

	return false
}

// Inspect returns container configuration.
func (s *Container) Inspect() ([]Inspect, error) {
	out, err := s.docker.Run("inspect", s.id)
	if err != nil {
		return nil, err
	}

	var list []Inspect
	err = json.NewDecoder(strings.NewReader(out)).Decode(&list)
	if err != nil {
		return nil, UnexpectedOutputError(fmt.Sprintf(
			"Error parsing output when inspecting container: %v", err))
	}
	if len(list) == 0 {
		return nil, UnexpectedOutputError(
			"Empty output when inspecting container")
	}

	return list, nil
}

// Kill terminates current container process.
func (s *Container) Kill() error {
	err := exec.Command(s.docker.binCmd, "kill", s.id).Run()
	if err != nil {
		return err
	}

	stopped := dot.WaitFunc(
		250*time.Millisecond, 30*time.Second, func() bool {
			inspect, err := s.Inspect()
			if err != nil || len(inspect) == 0 {
				return true
			}
			if !inspect[0].State.Running {
				return true
			}

			return false
		})
	if !stopped {
		err = fmt.Errorf("Timeout waiting '%s' container to stop", s.id)
	}

	return err
}

// NetworkNodes returns the network addresses and exposed ports of current
// container.
func (s *Container) NetworkNodes() ([]NetworkNode, error) {
	inspect, err := s.Inspect()
	if err != nil {
		return nil, err
	}

	if len(inspect) == 0 {
		return nil, fmt.Errorf("Container inspect is empty")
	}

	nodes := make([]NetworkNode, 0)
	for _, i := range inspect {
		ip := i.NetworkSettings.IPAddress
		if ip == "" {
			continue
		}

		if len(i.NetworkSettings.Ports) == 0 {
			nodes = append(nodes, NetworkNode{
				IpAddress: ip,
			})
		} else {
			for k, _ := range i.NetworkSettings.Ports {
				node := NetworkNode{}
				node.IpAddress = ip
				node.SetFromDocker(k)
				nodes = append(nodes, node)
			}
		}
	}

	return nodes, nil
}

// Remove erases current container from Docker.
func (s *Container) Remove() error {
	return exec.Command(s.docker.binCmd, "rm", s.id).Run()
}

// WaitStartup blocks current thread until current container begins to listen
// exposed port.
func (s *Container) WaitStartup(timeout time.Duration) error {
	nodes, err := s.NetworkNodes()
	if err != nil {
		return fmt.Errorf("Error getting network nodes: %v", err)
	}

	if !s.HasExposedPorts() {
		return fmt.Errorf("Current container has no exposed ports")
	}

	ok := dot.WaitPeerListening(
		nodes[0].Protocol, nodes[0].FormatDialAddress(), timeout)
	if !ok {
		err = fmt.Errorf("%s unreachable for %v",
			nodes[0].FormatDialAddress(), timeout)
	}

	return err
}
