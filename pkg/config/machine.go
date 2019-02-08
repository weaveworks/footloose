package config

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
}
