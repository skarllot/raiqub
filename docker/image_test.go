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
	"testing"
	"time"
)

const (
	REDIS_IMAGE_NAME = "redis"
	LISTEN_TIMEOUT   = 45 * time.Second
)

func TestRedisContainer(t *testing.T) {
	docker := NewDocker()
	if !docker.HasBin() {
		t.Skip("This test connot be run because Docker is not acessible")
	}

	image := NewImage(docker, REDIS_IMAGE_NAME)
	if err := image.Setup(); err != nil {
		t.Fatalf("Error setting up Docker environment: %v\n", err)
	}

	cfg := NewRunConfig()
	cfg.Detach()
	redis, err := image.Run(cfg)
	if err != nil {
		t.Fatalf("Error creating a new Docker container instance: %v\n", err)
	}
	defer redis.Remove()
	defer redis.Kill()

	if err := redis.WaitStartup(LISTEN_TIMEOUT); err != nil {
		t.Error("Timeout waiting for Redis instance to respond")
		return
	}
}
