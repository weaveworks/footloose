package ignite

import (
	"fmt"
	"net"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/footloose/pkg/config"
	"github.com/weaveworks/footloose/pkg/exec"
)

const (
	IgniteName = "ignite"
)

// This offset is incremented for each port so we avoid
// duplicate port bindings (and hopefully port collisions).
var portOffset uint16

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
		fmt.Sprintf("--ssh=%s", pubKeyPath),
	}

	copyFiles := spec.IgniteConfig().CopyFiles
	if copyFiles == nil {
		copyFiles = make(map[string]string)
	}
	for _, v := range setupCopyFiles(copyFiles) {
		runArgs = append(runArgs, v)
	}

	for _, mapping := range spec.PortMappings {
		if mapping.HostPort == 0 {
			// If not defined, set the host port to a random free ephemeral port
			var err error
			if mapping.HostPort, err = freePort(); err != nil {
				return "", err
			}
		} else {
			// If defined, apply an offset so all VMs won't use the same port
			mapping.HostPort += portOffset
		}

		runArgs = append(runArgs, fmt.Sprintf("--ports=%d:%d", int(mapping.HostPort), mapping.ContainerPort))
	}

	// increment portOffset per-machine
	portOffset++

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
	if ap, err := filepath.Abs(p); err == nil {
		return ap
	}
	// if Abs reports an error just return the original path 'p'
	return p
}

func IsCreated(name string) bool {
	err := exec.Command(execName, "inspect", "vm", name).Run()
	if err != nil {
		return false
	}
	return true
}

func IsStarted(name string) bool {
	cmd := exec.Command(execName, "inspect", "vm", name)
	lines, err := exec.CombinedOutputLines(cmd)
	if err != nil {
		log.Errorf("Ignite.IsStarted error:%v\n", err)
		return false
	}

	var sb strings.Builder
	for _, s := range lines {
		sb.WriteString(s)
	}

	data := []byte(sb.String())
	vm, err := toVM(data)
	if err != nil {
		return false
	}

	return vm.Status.Running
}

// freePort requests a free/open ephemeral port from the kernel
// Heavily inspired by https://github.com/phayes/freeport/blob/master/freeport.go
func freePort() (uint16, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()

	return uint16(l.Addr().(*net.TCPAddr).Port), nil
}
