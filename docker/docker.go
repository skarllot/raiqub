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
	"os/exec"
)

const (
	// Defines the default name for Docker binary.
	DOCKER_BIN = "docker"
)

// A Docker represents a Docker binary.
type Docker struct {
	binCmd string
}

// NewDocker creates a new instance of Docker with default binary name.
func NewDocker() *Docker {
	return &Docker{
		binCmd: DOCKER_BIN,
	}
}

// HasBin check whether Docker binary is available.
func (s *Docker) HasBin() bool {
	_, err := exec.LookPath(s.binCmd)
	return err == nil
}
