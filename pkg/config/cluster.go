package config

// MachineReplicas are a number of machine following the same specification.
type MachineReplicas struct {
	Spec  Machine `json:"spec"`
	Count int     `json:"count"`
}

// Cluster is a set of Machines.
type Cluster struct {
	// Name is the cluster name. Defaults to "cluster".
	Name string `json:"name"`

	// PrivateKey is the path to the private SSH key used to login into the cluster
	// machines.
	PrivateKey string `json:"privateKey"`
}

// Config is the top level config object.
type Config struct {
	// Cluster describes cluster-wide configuration.
	Cluster Cluster `json:"cluster"`
	// Machines describe the machines we want created for this cluster.
	Machines []MachineReplicas `json:"machines"`
}
