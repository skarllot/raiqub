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

package test

import (
	"testing"
	"time"

	"github.com/skarllot/raiqub/docker"
)

const (
	IMAGE_REDIS     = "redis"
	STARTUP_TIMEOUT = 3 * time.Minute
)

// A Environment represents a Docker testing environment.
type Environment struct {
	dockerBin *docker.Docker
	image     *docker.Image
	container *docker.Container
	run       func() (*docker.Container, error)
	test      *testing.T
}

// NewMongoDBEnvironment creates a instance that allows to use MongoDB for
// testing
func NewMongoDBEnvironment(t *testing.T) *Environment {
	d := docker.NewDocker()
	mongo := docker.NewImageMongoDB(d)

	return &Environment{
		dockerBin: d,
		image:     &mongo.Image,
		test:      t,
		run: func() (*docker.Container, error) {
			cfg := docker.NewRunConfig()
			cfg.Detach()
			return mongo.RunLight(cfg)
		},
	}
}

// NewRedisEnvironment creates a instance that allows to use Redis for testing.
func NewRedisEnvironment(t *testing.T) *Environment {
	d := docker.NewDocker()
	redis := docker.NewImage(d, IMAGE_REDIS)

	return &Environment{
		dockerBin: d,
		image:     redis,
		test:      t,
		run: func() (*docker.Container, error) {
			cfg := docker.NewRunConfig()
			cfg.Detach()
			return redis.Run(cfg)
		},
	}
}

// Applicability tests whether current testing environment can be run on current
// host.
func (s *Environment) Applicability() bool {
	return s.dockerBin.HasBin()
}

// Network returns network information from current running container.
func (s *Environment) Network() ([]docker.NetworkNode, error) {
	if s.container == nil {
		return nil, NotRunningError(s.image.Name())
	}
	
	nodes, err := s.container.NetworkNodes()
	if err != nil {
		return nil, err
	}
	
	return nodes, nil
}

// Run starts a new Docker instance for testing environment.
func (s *Environment) Run() bool {
	if err := s.image.Setup(); err != nil {
		s.test.Fatal("Error setting up Docker:", err)
		return false
	}

	var err error
	s.container, err = s.run()
	if err != nil {
		s.test.Fatal("Error running a new Docker container:", err)
		return false
	}

	if s.container.HasExposedPorts() {
		if err := s.container.WaitStartup(STARTUP_TIMEOUT); err != nil {
			s.test.Fatal("Timeout waiting Docker instance to respond:", err)
			return false
		}
	}

	return true
}

// Stop removes current running testing environment.
func (s *Environment) Stop() {
	if s.container == nil {
		return
	}

	s.container.Kill()
	s.container.Remove()
}
