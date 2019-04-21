package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetValueFromConfig(t *testing.T) {
	config := Config{
		Cluster: Cluster{Name: "clustername", PrivateKey: "privatekey"},
		Machines: []MachineReplicas{
			MachineReplicas{
				Count: 3,
				Spec: Machine{
					Image:      "myImage",
					Name:       "myName",
					Privileged: true,
				},
			},
		},
	}

	tests := []struct {
		name           string
		stringPath     string
		config         Config
		expectedOutput interface{}
	}{
		{
			"simple path select string",
			"cluster.name",
			Config{
				Cluster:  Cluster{Name: "clustername", PrivateKey: "privatekey"},
				Machines: []MachineReplicas{MachineReplicas{Count: 3, Spec: Machine{}}},
			},
			"clustername",
		},
		{
			"array path select global",
			"machines[0].spec",
			config,
			Machine{
				Image:      "myImage",
				Name:       "myName",
				Privileged: true,
			},
		},
		{
			"array path select bool",
			"machines[0].spec.Privileged",
			config,
			true,
		},
	}

	for _, utest := range tests {
		t.Run(utest.name, func(t *testing.T) {
			res, err := GetValueFromConfig(utest.stringPath, utest.config)
			assert.Nil(t, err)
			assert.Equal(t, utest.expectedOutput, res)
		})
	}
}
