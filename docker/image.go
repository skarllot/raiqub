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
	"fmt"
	"os/exec"
	"strings"
)

// A Image represents a Docker image.
type Image struct {
	docker *Docker
	name   string
}

// NewImage creates a new instance of Image.
func NewImage(d *Docker, name string) *Image {
	return &Image{
		docker: d,
		name:   name,
	}
}

// Exists check whether current Docker image is available.
func (s *Image) Exists() bool {
	out, err := exec.Command(s.docker.binCmd, "images", "--no-trunc").Output()
	if err != nil {
		return false
	}

	return bytes.Contains(out, []byte(s.name))
}

// Pull retrieves current Docker image from Docker repository.
//
// Returns DockerPullError on error.
func (s *Image) Pull() error {
	out, err := exec.Command(s.docker.binCmd, "pull", s.name).CombinedOutput()
	if err != nil {
		return DockerPullError{err, string(out)}
	}

	return nil
}

// Run creates a new Docker container as defined by current image and container
// template.
func (s *Image) Run(cfg *RunConfig) (*Container, error) {
	args := []string{"run"}
	args = append(args, cfg.Options...)
	args = append(args, s.name)
	args = append(args, cfg.Args...)

	cmd := exec.Command(s.docker.binCmd, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%v%v", stderr.String(), err)
	}

	container := &Container{
		docker: s.docker,
		id:     strings.TrimSpace(stdout.String()),
	}
	if container.id == "" {
		return nil, fmt.Errorf(
			"Unexpected empty output when running docker container")
	}

	return container, nil
}

// Setup check if Docker binary is available and pull current image from
// Docker repository whether not available.
func (s *Image) Setup() error {
	if !s.docker.HasBin() {
		return DockerBinNotFound(s.docker.binCmd)
	}

	if !s.Exists() {
		if err := s.Pull(); err != nil {
			return err
		}
	}

	return nil
}
