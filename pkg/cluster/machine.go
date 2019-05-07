package cluster

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/weaveworks/footloose/pkg/config"
	"github.com/weaveworks/footloose/pkg/docker"
)

type machineCache struct {
	// maps containerPort -> hostPort.
	ports map[int]int
}

// Machine is a single machine.
type Machine struct {
	spec *config.Machine

	// container name.
	name string
	// container hostname.
	hostname string
	// container ip.
	ip string

	// Fields that are cached from the docker daemon.
	machineCache
}

// ContainerName is the name of the running container corresponding to this
// Machine.
func (m *Machine) ContainerName() string {
	return m.name
}

// Hostname is the machine hostname.
func (m *Machine) Hostname() string {
	return m.hostname
}

// IsCreated returns if a machine is has been created. A created machine could
// either be running or stopped.
func (m *Machine) IsCreated() bool {
	res, _ := docker.Inspect(m.name, "{{.Name}}")
	if len(res) > 0 && len(res[0]) > 0 {
		return true
	}
	return false
}

// IsStarted returns if a machine is currently started or not.
func (m *Machine) IsStarted() bool {
	res, _ := docker.Inspect(m.name, "{{.State.Running}}")
	parsed, _ := strconv.ParseBool(strings.Trim(res[0], `'`))
	if parsed {
		return true
	}
	return false
}

// HostPort returns the host port corresponding to the given container port.
func (m *Machine) HostPort(containerPort int) (hostPort int, err error) {
	// Use the cached version first.
	if hostPort, ok := m.ports[containerPort]; ok {
		return hostPort, nil
	}
	// retrieve the specific port mapping using docker inspect
	lines, err := docker.Inspect(m.name, fmt.Sprintf("{{(index (index .NetworkSettings.Ports \"%d/tcp\") 0).HostPort}}", containerPort))
	if err != nil {
		return -1, errors.Wrap(err, "hostport: failed to inspect container")
	}
	if len(lines) != 1 {
		return -1, errors.Errorf("hostport: should only be one line, got %d lines", len(lines))
	}

	if m.ports == nil {
		m.ports = make(map[int]int)
	}

	port := strings.Replace(lines[0], "'", "", -1)
	m.ports[containerPort], err = strconv.Atoi(port)
	if err != nil {
		return -1, errors.Wrap(err, "hostport: failed to parse string to int")
	}
	return m.ports[containerPort], nil
}
