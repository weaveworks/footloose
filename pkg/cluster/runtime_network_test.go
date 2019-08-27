package cluster

import (
	"testing"

	"github.com/docker/docker/api/types/network"
	"github.com/stretchr/testify/assert"
)

func TestNewRuntimeNetworks(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		networks := map[string]*network.EndpointSettings{}
		networks["mynetwork"] = &network.EndpointSettings{
			Gateway:     "172.17.0.1",
			IPAddress:   "172.17.0.4",
			IPPrefixLen: 16,
		}
		res := NewRuntimeNetworks(networks)

		expectedRuntimeNetworks := []*RuntimeNetwork{
			&RuntimeNetwork{Name: "mynetwork", Gateway: "172.17.0.1", IP: "172.17.0.4", Mask: "255.255.0.0"}}
		assert.Equal(t, expectedRuntimeNetworks, res)
	})
}
