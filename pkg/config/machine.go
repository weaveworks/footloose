package config

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Volume is a volume that can be attached to a Machine.
type Volume struct {
	// Type is the volume type. One of "bind" or "volume".
	Type string `json:"type" yaml:"type,omitempty"`
	// Source is the volume source.
	// With type=bind, the volume source is a directory or a file in the host
	// filesystem.
	// With type=volume, source is either the name of a docker volume or "" for
	// anonymous volumes.
	Source string `json:"source" yaml:"source,omitempty"`
	// Destination is the mount point inside the container.
	Destination string `json:"destination" yaml:"destination,omitempty"`
	// ReadOnly specifies if the volume should be read-only or not.
	ReadOnly bool `json:"readOnly" yaml:"readOnly,omitempty"`
}

// PortMapping describes mapping a port from the machine onto the host.
type PortMapping struct {
	// Protocol is the layer 4 protocol for this mapping. One of "tcp" or "udp".
	// Defaults to "tcp".
	Protocol string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	// Address is the host address to bind to. Defaults to "0.0.0.0".
	Address string `json:"address,omitempty" yaml:"address,omitempty"`
	// HostPort is the base host port to map the containers ports to. As we
	// configure a number of machine replicas, each machine will use HostPort+i
	// where i is between 0 and N-1, N being the number of machine replicas. If 0,
	// a local port will be automatically allocated.
	HostPort int `json:"hostPort,omitempty" yaml:"hostPort,omitempty"`
	// ContainerPort is the container port to map.
	ContainerPort int `json:"containerPort" yaml:"containerPort,omitempty"`
}

// Machine is the machine configuration.
type Machine struct {
	// Name is the machine name. This is a format string with %d as the machine
	// index, a number between 0 and N-1, N being the number of machines in the
	// cluster. This name will also be used as the machine hostname. Defaults to
	// "node%d".
	Name string `json:"name" yaml:"name,omitempty"`
	// Image is the container image to use for this machine.
	Image string `json:"image" yaml:"image,omitempty"`
	// Privileged controls whether to start the Machine as a privileged container
	// or not. Defaults to false.
	Privileged bool `json:"privileged,omitempty" yaml:"privileged,omitempty"`
	// Volumes is the list of volumes attached to this machine.
	Volumes []Volume `json:"volumes,omitempty" yaml:"volumes,omitempty"`
	// PortMappings is the list of ports to expose to the host.
	PortMappings []PortMapping `json:"portMappings,omitempty" yaml:"portMappings,omitempty"`
	// Cmd is a cmd which will be run in the container.
	Cmd string `json:"cmd,omitempty" yaml:"cmd,omitempty"`
}

// validate checks basic rules for Machine's fields
func (conf Machine) validate() (rerr error) {
	rerr = nil
	validName := strings.Contains(conf.Name, "%d")
	if validName != true {
		log.Warnf("Machine conf validation: machine name %v is not valid, it should contains %%d", conf.Name)
		rerr = fmt.Errorf("Machine configuration not valid")
	}
	for _, pmapping := range conf.PortMappings {
		if err := pmapping.validate(); err != nil {
			log.Warn(err)
			rerr = fmt.Errorf("Machine configuration not valid")
		}
	}
	return rerr
}

func (conf PortMapping) validate() error {
	if conf.HostPort > maxPort || conf.HostPort < minPort {
		return fmt.Errorf("Machine conf validation: hostPort %v is not valid, it cannot be hight than %v or lesser than %v",
			conf.HostPort,
			maxPort,
			minPort)
	}
	if conf.ContainerPort > maxPort || conf.ContainerPort < minPort {
		return fmt.Errorf("Machine conf validation: containerPort %v is not valid, it cannot be hight than %v or lesser than %v",
			conf.ContainerPort,
			maxPort,
			minPort)
	}
	return nil
}
