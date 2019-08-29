package cluster

import (
	"fmt"
	"strconv"
	"strings"
	"syscall"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/footloose/pkg/config"
	"github.com/weaveworks/footloose/pkg/docker"
	"github.com/weaveworks/footloose/pkg/exec"
	"github.com/weaveworks/footloose/pkg/ignite"
)

// Machine is a single machine.
type Machine struct {
	spec *config.Machine

	// container name.
	name string
	// container hostname.
	hostname string
	// container ip.
	ip string

	runtimeNetworks []*RuntimeNetwork
	// Fields that are cached from the docker daemon.

	ports map[int]int
	// maps containerPort -> hostPort.
}

// ContainerName is the name of the running container corresponding to this
// Machine.
func (m *Machine) ContainerName() string {
	if m.IsIgnite() {
		filter := fmt.Sprintf(`label=ignite.name=%s`, m.name)
		cid, err := exec.ExecuteCommand("docker", "ps", "-q", "-f", filter)
		if err != nil || len(cid) == 0 {
			return m.name
		}
		return cid
	}
	return m.name
}

// Hostname is the machine hostname.
func (m *Machine) Hostname() string {
	return m.hostname
}

// IsCreated returns if a machine is has been created. A created machine could
// either be running or stopped.
func (m *Machine) IsCreated() bool {
	if m.IsIgnite() {
		return ignite.IsCreated(m.name)
	}

	res, _ := docker.Inspect(m.name, "{{.Name}}")
	if len(res) > 0 && len(res[0]) > 0 {
		return true
	}
	return false
}

// IsStarted returns if a machine is currently started or not.
func (m *Machine) IsStarted() bool {
	if m.IsIgnite() {
		return ignite.IsStarted(m.name)
	}

	res, _ := docker.Inspect(m.name, "{{.State.Running}}")
	parsed, _ := strconv.ParseBool(strings.Trim(res[0], `'`))
	return parsed
}

// HostPort returns the host port corresponding to the given container port.
func (m *Machine) HostPort(containerPort int) (hostPort int, err error) {
	// Use the cached version first.
	if hostPort, ok := m.ports[containerPort]; ok {
		return hostPort, nil
	}
	// retrieve the specific port mapping using docker inspect
	lines, err := docker.Inspect(m.ContainerName(), fmt.Sprintf("{{(index (index .NetworkSettings.Ports \"%d/tcp\") 0).HostPort}}", containerPort))
	if err != nil {
		return -1, errors.Wrapf(err, "hostport: failed to inspect container: %v", lines)
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

// Status returns the machine status.
func (m *Machine) Status() *MachineStatus {
	s := MachineStatus{}

	s.Container = m.ContainerName()
	s.Image = m.spec.Image
	s.Command = m.spec.Cmd
	s.Spec = m.spec
	state := NotCreated
	if m.IsCreated() {
		state = Stopped
		if m.IsStarted() {
			state = Running
		}
	}
	s.State = state
	var ports []port
	for k, v := range m.ports {
		p := port{
			Host:  v,
			Guest: k,
		}
		ports = append(ports, p)
	}
	if len(ports) < 1 {
		for _, p := range m.spec.PortMappings {
			ports = append(ports, port{Host: int(p.ContainerPort), Guest: 0})
		}
	}
	s.Ports = ports
	s.RuntimeNetworks = m.runtimeNetworks

	return &s
}

// Only check for Ignite prerequisites once
var igniteChecked bool

func (m *Machine) IsIgnite() (b bool) {
	b = m.spec.Backend == ignite.BackendName

	if !igniteChecked && b {
		if syscall.Getuid() != 0 {
			log.Fatalf("Footloose needs to run as root to use the %q backend", ignite.BackendName)
		}

		ignite.CheckVersion()
		igniteChecked = true
	}

	return
}
