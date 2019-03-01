package main

import "github.com/weaveworks/footloose/pkg/config"

var defaultConfig = config.Config{
	Cluster: config.Cluster{
		Name:       "cluster",
		PrivateKey: "cluster-key",
	},
	Machines: []config.MachineReplicas{{
		Count: 1,
		Spec: config.Machine{
			Name:  "node%d",
			Image: "quay.io/footloose/centos7:0.1.0",
			PortMappings: []config.PortMapping{{
				ContainerPort: 22,
			}},
		},
	}},
}
