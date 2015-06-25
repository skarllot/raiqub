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
	out, err := s.docker.Run("images", "--no-trunc")
	if err != nil {
		return false
	}

	return strings.Contains(out, s.name)
}

// Name returns image name from current instance.
func (s *Image) Name() string {
	return s.name
}

// Pull retrieves current Docker image from Docker repository.
//
// Returns ExternalCmdError on error.
func (s *Image) Pull() error {
	_, err := s.docker.Run("pull", s.name)
	if err != nil {
		return err
	}

	return nil
}

// Run creates a new Docker container as defined by current image and container
// template.
//
// Returns ExternalCmdError or UnexpectedOutputError on error.
func (s *Image) Run(cfg *RunConfig) (*Container, error) {
	args := make([]string, 0, len(cfg.Options)+len(cfg.Args)+1)
	args = append(args, cfg.Options...)
	args = append(args, s.name)
	args = append(args, cfg.Args...)

	out, err := s.docker.Run("run", args...)
	if err != nil {
		return nil, err
	}

	container := &Container{
		docker: s.docker,
		id:     strings.TrimSpace(out),
	}
	if container.id == "" {
		return nil, UnexpectedOutputError(
			"Unexpected empty output when running docker container")
	}

	return container, nil
}

// Setup check if Docker binary is available and pull current image from Docker
// repository in case it is not already available.
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
