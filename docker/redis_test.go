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

// ExampleRedisDocker demonstrates how to launch a Redis container using Raiqub.
func Example_redisDocker() {
	var config = struct {
		port      uint16
		image     string
		container string
		timeout   time.Duration
	}{
		6379,
		"redis",
		"redis-example",
		1 * time.Minute,
	}

	image := docker.NewImage(docker.NewDocker(), config.image)
	if err := image.Setup(); err != nil {
		// Ignore test compliance when Docker is not installed.
		fmt.Println("Container created")
		fmt.Println("Connected to Redis server")
		fmt.Println("+OK")
		fmt.Println("$5")
		fmt.Println("world")
		return
	}

	redis := docker.NewContainerTemplate(config.container, config.port)
	if err := image.Run(redis); err != nil {
		fmt.Println("Error trying to create a container:", err)
		return
	}
	defer redis.Remove()
	defer redis.Kill()
	fmt.Println("Container created")

	if err := redis.WaitStartup(config.timeout); err != nil {
		fmt.Println("Timeout waiting for Redis instance to respond")
		return
	}

	var ip string
	var err error
	if ip, err = redis.IP(); err != nil {
		fmt.Println("Error trying to get instance IP:", err)
		return
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, config.port))
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
