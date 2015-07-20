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

// A DockerBinNotFound represents an error when Docker binary could not be
// found.
type DockerBinNotFound string

// Error returns string representation of current instance error.
func (e DockerBinNotFound) Error() string {
	return fmt.Sprintf("The docker binary '%s' was not found", string(e))
}

// A ExternalCmdError represents an error when running external command fails.
type ExternalCmdError struct {
	// The error returned native library.
	InnerError error
	// Standard error output from command-line.
	Stderr string
	// Standard output from command-line.
	Stdout string
}

// Error returns string representation of current instance error.
func (e ExternalCmdError) Error() string {
	return fmt.Sprintf("Error running external command: %v", e.InnerError)
}

// A UnexpectedOutputError represents an error when an external command prints
// an unexpected content.
type UnexpectedOutputError string

// Error returns string representation of current instance error.
func (e UnexpectedOutputError) Error() string {
	return fmt.Sprint(string(e))
}
