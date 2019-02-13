package config

// Volume is a volume that can be attached to a Machine.
type Volume struct {
	// Type is the volume type. One of "bind" or "volume".
	Type string `json:"type"`
	// Source is the volume source.
	// With type=bind, the volume source is a directory or a file in the host
	// filesystem.
	// With type=volume, source is either the name of a docker volume or "" for
	// anonymous volumes.
	Source string `json:"source"`
	// Destination is the mount point inside the container.
	Destination string `json:"destination"`
	// ReadOnly specifies if the volume should be read-only or not.
	ReadOnly bool `json:"readOnly"`
}

// PortMapping describes mapping a port from the machine onto the host.
type PortMapping struct {
	// Protocol is the layer 4 protocol for this mapping. One of "tcp" or "udp".
	// Defaults to "tcp".
	Protocol string `json:"protocol,omitempty"`
	// Address is the host addres to bind to. Defaults to "0.0.0.0".
	Address string `json:"address,omitempty"`
	// HostPort is the base host port to map the containers ports to. As we
	// configure a number of machine replicas, each machine will use HostPort+i
	// where i is between 0 and N-1, N being the number of machine replicas. If 0,
	// a local port will be automatically allocated.
	HostPort uint16 `json:"hostPort,omitempty"`
	// ContainerPort is the container port to map.
	ContainerPort uint16 `json:"containerPort"`
}

// Machine is the machine configuration.
type Machine struct {
	// Name is the machine name. This is a format string with %d as the machine
	// index, a number between 0 and N-1, N being the number of machines in the
	// cluster. This name will also be used as the machine hostname. Defaults to
	// "node%d".
	Name string `json:"name"`
	// Image is the container image to use for this machine.
	Image string `json:"image"`
	// Privileged controls whether to start the Machine as a privileged container
	// or not. Defaults to false.
	Privileged bool `json:"privileged,omitempty"`
	// Volumes is the list of volumes attached to this machine.
	Volumes []Volume `json:"volumes,omitempty"`
	// PortMappings is the list of ports to expose to the host.
	PortMappings []PortMapping `json:"portMappings,omitempty"`
}
