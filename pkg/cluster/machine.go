package cluster

import "github.com/dlespiau/footloose/pkg/config"

// Machine is a single machine.
type Machine struct {
	spec *config.Machine

	// container name.
	name     string
	hostname string
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
