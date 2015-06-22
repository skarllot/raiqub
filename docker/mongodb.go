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

const (
	// Defines the name of MongoDB image from Docker repository.
	IMAGE_MONGODB_NAME = "mongo"
	// Defines default port for MongoDB.
	MONGODB_PORT_DEFAULT = 27017
)

// A ImageMongoDB represents a MongoDB Docker image.
type ImageMongoDB struct {
	Image
}

// NewImageMongoDB creates a new Image pre-configured to MongoDB image.
func NewImageMongoDB(d *Docker) *ImageMongoDB {
	return &ImageMongoDB{
		*NewImage(d, IMAGE_MONGODB_NAME),
	}
}

// RunLight creates a light instance of MongoDB image.
func (s *ImageMongoDB) RunLight(cfg *RunConfig) (*Container, error) {
	cfg.AddArgs("--smallfiles", "--nojournal")
	return s.Image.Run(cfg)
}
