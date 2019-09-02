package cluster

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchFilter(t *testing.T) {
	const refused = "ssh: connect to host 172.17.0.2 port 22: Connection refused"

	filter := matchFilter{
		writer: ioutil.Discard,
		regexp: connectRefused,
	}

	_, err := filter.Write([]byte("foo\n"))
	assert.NoError(t, err)
	assert.Equal(t, false, filter.matched)

	_, err = filter.Write([]byte(refused))
	assert.NoError(t, err)
	assert.Equal(t, false, filter.matched)
}

func TestNewClusterWithHostPort(t *testing.T) {
	cluster, err := NewFromYAML([]byte(`cluster:
  name: cluster
  privateKey: cluster-key
machines:
- count: 2
  spec:
    image: quay.io/footloose/centos7
    name: node%d
    portMappings:
    - containerPort: 22
      hostPort: 2222
`))
	assert.NoError(t, err)
	assert.NotNil(t, cluster)
	assert.Equal(t, 1, len(cluster.spec.Machines))
	template := cluster.spec.Machines[0]
	assert.Equal(t, 2, template.Count)
	assert.Equal(t, 1, len(template.Spec.PortMappings))
	portMapping := template.Spec.PortMappings[0]
	assert.Equal(t, uint16(22), portMapping.ContainerPort)
	assert.Equal(t, uint16(2222), portMapping.HostPort)

	machine0 := cluster.machine(&template.Spec, 0)
	args0 := cluster.createMachineRunArgs(machine0, machine0.ContainerName(), 0)
	i := indexOf("-p", args0)
	assert.NotEqual(t, -1, i)
	assert.Equal(t, "2222:22", args0[i+1])

	machine1 := cluster.machine(&template.Spec, 1)
	args1 := cluster.createMachineRunArgs(machine1, machine1.ContainerName(), 1)
	i = indexOf("-p", args1)
	assert.NotEqual(t, -1, i)
	assert.Equal(t, "2223:22", args1[i+1])
}

func indexOf(element string, array []string) int {
	for k, v := range array {
		if element == v {
			return k
		}
	}
	return -1 // element not found.
}
