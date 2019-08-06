/*
 * Minio Client (C) 2014, 2015, 2016 Minio, Inc.
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

package cmd

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/didiyun/mc/pkg/console"
)

// check if minimum Go version is met.
func checkGoVersion() {
	runtimeVersion := runtime.Version()

	// Checking version is always successful with go tip
	if strings.HasPrefix(runtimeVersion, "devel") {
		return
	}

	// Parsing golang version
	curVersion, e := version.NewVersion(runtimeVersion[2:])
	if e != nil {
		console.Fatalln("Unable to determine current go version.", e)
	}

	// Prepare version constraint.
	constraints, e := version.NewConstraint(minGoVersion)
	if e != nil {
		console.Fatalln("Unable to check go version.")
	}

	// Check for minimum version.
	if !constraints.Check(curVersion) {
		console.Fatalln(fmt.Sprintf("Please recompile Minio with Golang version %s.", minGoVersion))
	}
}
