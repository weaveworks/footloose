package config

import (
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func NewConfigFromYAML(data []byte) (*Config, error) {
	spec := Config{}
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, err
	}
	return &spec, nil
}

func NewConfigFromFile(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return NewConfigFromYAML(data)
}

// Save writes the Config to a file.
func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0666)
}

// MachineReplicas are a number of machine following the same specification.
type MachineReplicas struct {
	Spec  Machine `json:"spec" yaml:"spec,omitempty"`
	Count int     `json:"count" yaml:"count,omitempty"`
}

// Cluster is a set of Machines.
type Cluster struct {
	// Name is the cluster name. Defaults to "cluster".
	Name string `json:"name" yaml:"name,omitempty"`

	// PrivateKey is the path to the private SSH key used to login into the cluster
	// machines. Can be expanded to user homedir if ~ is found. Ex. ~/.ssh/id_rsa
	PrivateKey string `json:"privateKey" yaml:"privateKey,omitempty"`
}

// Config is the top level config object.
type Config struct {
	// Cluster describes cluster-wide configuration.
	Cluster Cluster `json:"cluster" yaml:"cluster,omitempty"`
	// Machines describe the machines we want created for this cluster.
	Machines []MachineReplicas `json:"machines" yaml:"machines,omitempty"`
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
