package main

import "github.com/dlespiau/footloose/pkg/config"

var defaultMachineSpec = config.Machine{
	Name:  "node%d",
	Image: "quay.io/footloose/centos7",
}

var defaultClusterSpec = config.Cluster{
	Name:       "cluster",
	PrivateKey: "cluster-key",
	Templates: []config.MachineReplicas{{
		Spec:  defaultMachineSpec,
		Count: 1,
	}},
}
