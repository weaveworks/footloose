package main

import "github.com/dlespiau/footloose/pkg/config"

var machineSpec = config.Machine{
	Name:  "node%d",
	Image: "quay.io/footloose/centos7",
}

var clusterSpec = config.Cluster{
	Name: "cluster",
	Templates: []config.MachineReplicas{{
		Spec:  machineSpec,
		Count: 1,
	}},
}
