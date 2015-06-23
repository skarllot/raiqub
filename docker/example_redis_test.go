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

package docker_test

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/skarllot/raiqub/docker"
)

const (
	REDIS_PORT     = 6379
	REDIS_IMAGE    = "redis"
	LISTEN_TIMEOUT = 45 * time.Second
)

// ExampleRedisDocker demonstrates how to launch a Redis container using Raiqub.
func Example_redisDocker() {
	dockerBin := docker.NewDocker()
	if !dockerBin.HasBin() {
		// Ignore test compliance when Docker is not installed.
		fmt.Println("Container created")
		fmt.Println("Connected to Redis server")
		fmt.Println("+OK")
		fmt.Println("$5")
		fmt.Println("world")
		return
	}

	image := docker.NewImage(dockerBin, REDIS_IMAGE)
	if err := image.Setup(); err != nil {
		fmt.Println("Error setting up Docker environment:", err)
		return
	}

	cfg := docker.NewRunConfig()
	cfg.Detach()
	redis, err := image.Run(cfg)
	if err != nil {
		fmt.Println("Error trying to create a container:", err)
		return
	}
	defer redis.Remove()
	defer redis.Kill()
	fmt.Println("Container created")

	if err := redis.WaitStartup(LISTEN_TIMEOUT); err != nil {
		fmt.Println("Timeout waiting for Redis instance to respond")
		return
	}

	nodes, err := redis.NetworkNodes()
	if err != nil {
		fmt.Println("Error trying to get network information:", err)
		return
	}

	conn, err := net.Dial(nodes[0].Protocol, nodes[0].FormatDialAddress())
	if err != nil {
		fmt.Println("Could not connect to Redis server:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected to Redis server")

	reader := bufio.NewReader(conn)
	fmt.Fprintln(conn, "SET hello world")
	out, _ := reader.ReadString('\n')
	fmt.Println(out[:3])

	fmt.Fprintln(conn, "GET hello")
	out, _ = reader.ReadString('\n')
	fmt.Println(out[:2])
	out, _ = reader.ReadString('\n')
	fmt.Println(out[:5])

	// Output:
	// Container created
	// Connected to Redis server
	// +OK
	// $5
	// world
}
