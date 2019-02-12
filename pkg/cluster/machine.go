package cluster

import (
	"fmt"
	"strconv"

	"github.com/dlespiau/footloose/pkg/config"
	"github.com/pkg/errors"
	"sigs.k8s.io/kind/pkg/docker"
)

type machineCache struct {
	// maps containerPort -> hostPort.
	ports map[int]int
}

// Machine is a single machine.
type Machine struct {
	spec *config.Machine

	// container name.
	name     string
	hostname string

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

	m.ports[containerPort], err = strconv.Atoi(lines[0])
	if err != nil {
		return -1, errors.Wrap(err, "hostport: failed to parse string to int")
	}
	return m.ports[containerPort], nil
}
