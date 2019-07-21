/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ignite

import (
	"fmt"

	"github.com/weaveworks/footloose/pkg/config"
	"github.com/weaveworks/footloose/pkg/exec"
)

const (
	IgniteName = "ignite"
)

// Create creates a container with "docker create", with some error handling
// it will return the ID of the created container if any, even on error
func Create(name string, spec *config.Machine, pubKeyPath string) (id string, err error) {

	runArgs := []string{
		"run",
		spec.Image,
		fmt.Sprintf("--name=%s", name),
		fmt.Sprintf("--cpus=%d", spec.IgniteConfig().CPUs),
		fmt.Sprintf("--memory=%s", spec.IgniteConfig().Memory),
		fmt.Sprintf("--size=%s", spec.IgniteConfig().Disk),
		fmt.Sprintf("--kernel-image=%s", spec.IgniteConfig().Kernel),
	}

	copyFiles := spec.IgniteConfig().CopyFiles
	if copyFiles == nil {
		copyFiles = make(map[string]string)
	}
	copyFiles[pubKeyPath] = "/root/.ssh/authorized_keys"
	for _, v := range setupCopyFiles(copyFiles) {
		runArgs = append(runArgs, v)
	}

	for i, mapping := range spec.PortMappings {
		if mapping.HostPort == 0 {
			// TODO: should warn here as containerPort is dropped
			continue
		}
		runArgs = append(runArgs, fmt.Sprintf("--ports=%d:%d", int(mapping.HostPort)+i, mapping.ContainerPort))
	}

	_, err = exec.ExecuteCommand(execName, runArgs...)
	return "", err
}

func setupCopyFiles(copyFiles map[string]string) []string {
	ret := []string{}
	for k, v := range copyFiles {
		s := fmt.Sprintf("--copy-files=%s:%s", k, v)
		ret = append(ret, s)
	}
	return ret
}

func IsCreated(name string) bool {
	exitCode, err := exec.ExecForeground(execName, "logs", name)
	if err != nil || exitCode != 0 {
		return false
	}
	return true
}
