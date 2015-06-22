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
)

// A RunConfig represents options and arguments used to run a Docker image.
type RunConfig struct {
	Options []string
	Args    []string
}

// NewRunConfig creates a new instance of RunConfig.
func NewRunConfig() *RunConfig {
	return &RunConfig{
		make([]string, 0),
		make([]string, 0),
	}
}

// Detach run the container in the background.
func (s *RunConfig) Detach() {
	s.Options = append(s.Options, "-d")
}

// PublishPort publish a container's port to the host.
func (s *RunConfig) PublishPort(hostPort, containerPort uint16) {
	s.Options = append(s.Options, "-p",
		fmt.Sprintf("%d:%d", hostPort, containerPort))
}

// Name assign a name to the container.
func (s *RunConfig) Name(name string) {
	s.Options = append(s.Options, "--name", name)
}

// AddArgs appends arguments to Docker container.
func (s *RunConfig) AddArgs(args ...string) {
	s.Args = append(s.Args, args...)
}
