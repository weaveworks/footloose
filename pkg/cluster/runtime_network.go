package cluster

import (
	"net"

	"github.com/docker/docker/api/types/network"
)

const (
	ipv4Length = 32
)

func NewRuntimeNetworks(networks map[string]*network.EndpointSettings) []*RuntimeNetwork {
	rnList := make([]*RuntimeNetwork, 0, len(networks))
	for key, value := range networks {
		mask := net.CIDRMask(value.IPPrefixLen, ipv4Length)
		maskIP := net.IP(mask).String()
		rnNetwork := &RuntimeNetwork{
			Name:    key,
			IP:      value.IPAddress,
			Mask:    maskIP,
			Gateway: value.Gateway,
		}
		rnList = append(rnList, rnNetwork)
	}
	return rnList
}

type RuntimeNetwork struct {
	// Name of the network
	Name string `json:"name,omitempty"`
	// IP of the container
	IP string `json:"ip,omitempty"`
	// Mask of the network
	Mask string `json:"mask,omitempty"`
	// Gateway of the network
	Gateway string `json:"gateway,omitempty"`
}
