package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newConfigForTest(name, key, image, machineName string, privileged bool, count int) *Config {
	return &Config{
		Cluster: Cluster{Name: name, PrivateKey: key},
		Machines: []MachineReplicas{
			MachineReplicas{
				Count: count,
				Spec: Machine{
					Image:      image,
					Name:       machineName,
					Privileged: privileged,
				},
			},
		},
	}
}

func TestSetValueToConfig(t *testing.T) {

	tests := []struct {
		name           string
		stringPath     string
		newValue       interface{}
		config         *Config
		expectedOutput interface{}
		expectedErr    bool
	}{
		{
			"simple path set string",
			"cluster.name",
			"new-clu",
			newConfigForTest("cluster", "key", "some-image", "node%d", true, 2),
			newConfigForTest("new-clu", "key", "some-image", "node%d", true, 2),
			false,
		},
		{
			"array path set int",
			"machines[0].count",
			3,
			newConfigForTest("cluster", "key", "some-image", "node%d", true, 2),
			newConfigForTest("cluster", "key", "some-image", "node%d", true, 3),
			false,
		},
		{
			"array path set bool",
			"machines[0].spec.privileged",
			false,
			newConfigForTest("cluster", "key", "some-image", "node%d", true, 2),
			newConfigForTest("cluster", "key", "some-image", "node%d", false, 2),
			false,
		},
		{
			"array path set bool to non bool var",
			"machines[0].spec",
			false,
			newConfigForTest("cluster", "key", "some-image", "node%d", true, 2),
			newConfigForTest("cluster", "key", "some-image", "node%d", true, 2),
			true,
		},
		{
			"array path set int to non int var",
			"cluster.name",
			1,
			newConfigForTest("cluster", "key", "some-image", "node%d", true, 2),
			newConfigForTest("cluster", "key", "some-image", "node%d", true, 2),
			true,
		},
		{
			"array path set string to non string var",
			"machines[0].count",
			"value",
			newConfigForTest("cluster", "key", "some-image", "node%d", true, 2),
			newConfigForTest("cluster", "key", "some-image", "node%d", true, 2),
			true,
		},
		{
			"array path set int to out of bound of array",
			"machines[2].count",
			1,
			newConfigForTest("cluster", "key", "some-image", "node%d", true, 2),
			newConfigForTest("cluster", "key", "some-image", "node%d", true, 2),
			true,
		},
	}

	for _, utest := range tests {
		t.Run(utest.name, func(t *testing.T) {
			err := SetValueToConfig(utest.stringPath, utest.config, utest.newValue)
			if utest.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, utest.expectedOutput, utest.config)
		})
	}
}

func TestIsSetValueValid(t *testing.T) {
	tests := []struct {
		name          string
		stringPath    string
		value         string
		expectedError bool
	}{
		{
			"machine name invalid",
			"machines[0].spec.name",
			"myMachineName",
			true,
		},
		{
			"machine name valid",
			"machines[0].spec.name",
			"myMachine%dName",
			false,
		},
		{
			"machine name invalid for uppercase path",
			"Machines[0].Spec.Name",
			"myMachineName",
			true,
		},
		{
			"port too high invalid containerPort",
			"machines[0].spec.portMappings[0].containerPort",
			"65536",
			true,
		},
		{
			"port too high invalid containerPort for uppercase path",
			"Machines[0].Spec.PortMappings[0].ContainerPort",
			"65536",
			true,
		},
		{
			"port too high invalid hostPort",
			"machines[0].spec.portMappings[0].hostPort",
			"65536",
			true,
		},
		{
			"port too high invalid hostPort for uppercase path",
			"Machines[0].Spec.PortMappings[0].HostPort",
			"65536",
			true,
		},
		{
			"port lower than 1 invalid containerPort",
			"machines[0].spec.portMappings[0].containerPort",
			"0",
			true,
		},
		{
			"port lower than 1 invalid hostPort",
			"machines[0].spec.portMappings[0].hostPort",
			"0",
			true,
		},
		{
			"port valid containerPort",
			"machines[0].spec.portMappings[0].containerPort",
			"22",
			false,
		},
		{
			"port valid hostPort",
			"machines[0].spec.portMappings[0].hostPort",
			"22",
			false,
		},
	}

	for _, utest := range tests {
		t.Run(utest.name, func(t *testing.T) {
			err := IsSetValueValid(utest.stringPath, utest.value)
			if utest.expectedError == true {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
