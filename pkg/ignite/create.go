package ignite

import (
	"fmt"
	"path"
	"path/filepath"

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
		s := fmt.Sprintf("--copy-files=%s:%s", toAbs(k), v)
		ret = append(ret, s)
	}
	return ret
}

func toAbs(p string) string {
	ap := p
	if !path.IsAbs(ap) {
		aap, err := filepath.Abs(ap)
		ap = aap
		// if Abs reports an error just return the original path 'p'
		if err != nil {
			ap = p
		}
	}
	return ap
}

func IsCreated(name string) bool {
	exitCode, err := exec.ExecForeground(execName, "logs", name)
	if err != nil || exitCode != 0 {
		return false
	}
	return true
}
