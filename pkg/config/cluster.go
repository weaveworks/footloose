package config

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

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
	// machines. Can be expanded to user homedir if ~ is found. Ex. ~/.ssh/id_rsa
	PrivateKey string `json:"privateKey"`
}

// Config is the top level config object.
type Config struct {
	// Cluster describes cluster-wide configuration.
	Cluster Cluster `json:"cluster"`
	// Machines describe the machines we want created for this cluster.
	Machines []MachineReplicas `json:"machines"`
}

// validate checks basic rules for MachineReplicas's fields
func (conf MachineReplicas) validate() error {
	return conf.Spec.validate()
}

// validate checks basic rules for Cluster's fields
func (conf Cluster) validate() error {
	return fmt.Errorf("not yet implemented")
}

// Validate checks basic rules for Config's fields
func (conf Config) Validate() error {
	valid := true
	for _, machine := range conf.Machines {
		err := machine.validate()
		if err != nil {
			valid = false
			log.Fatalf(err.Error())
		}
	}
	if valid == false {
		return fmt.Errorf("Configuration file non valid")
	}
	return nil
}
