package ignite

import (
	"encoding/json"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/footloose/pkg/exec"
)

type Metadata struct {
	Name    string
	UID     string
	Created string
}

type Port struct {
	HostPort uint16
	VMPort   uint16
	Protocol string
}

type Network struct {
	Ports []Port
}

type Spec struct {
	Network  Network
	Cpus     uint
	Memory   string
	DiskSize string
}

type Status struct {
	Running     bool
	StartTime   string
	IpAddresses []string
}

type VM struct {
	Metadata Metadata
	Spec     Spec
	Status   Status
}

func PopulateMachineDetails(name string) (*VM, error) {
	cmd := exec.Command(execName, "inspect", "vm", name)
	lines, err := exec.CombinedOutputLines(cmd)
	if err != nil {
		log.Errorf("Ignite.IsStarted error:%v\n", err)
		return nil, err
	}

	var sb strings.Builder
	for _, s := range lines {
		sb.WriteString(s)
	}
	return toVM([]byte(sb.String()))
}

func toVM(data []byte) (*VM, error) {
	obj := &VM{}
	err := json.Unmarshal(data, obj)
	if err != nil {
		log.Errorf("Ignite.toVM error:%v\n", err)
		return nil, err
	}
	return obj, nil
}
