package main

import "github.com/weaveworks/footloose/pkg/config"

// imageTag computes the docker image tag given the footloose version.
func imageTag(v string) string {
	if v == "git" {
		return "latest"
	}
	return v
}

// defaultKeyStore is the path where to store the public keys.
const defaultKeyStorePath = "keys"

var defaultConfig = config.Config{
	Cluster: config.Cluster{
		Name:       "cluster",
		PrivateKey: "cluster-key",
	},
	Machines: []config.MachineReplicas{{
		Count: 1,
		Spec: config.Machine{
			Name:  "node%d",
			Image: "quay.io/footloose/centos7:" + imageTag(version),
			PortMappings: []config.PortMapping{{
				ContainerPort: 22,
			}},
			Backend: "docker",
		},
	}},
}
