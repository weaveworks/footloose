package cluster

import (
	"fmt"

	"github.com/dlespiau/footloose/pkg/config"
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/kind/pkg/docker"
)

// Container represents a running machine.
type Container struct {
	ID string
}

// Cluster is a running cluster.
type Cluster struct {
	spec config.Cluster
}

// New creates a new cluster. It takes as input the description of the cluster
// and its machines.
func New(cluster config.Cluster) *Cluster {
	return &Cluster{
		spec: cluster,
	}
}

func f(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

func (c *Cluster) containerName(machine *config.Machine, i int) string {
	format := "%s-" + machine.Name
	return f(format, c.spec.Name, i)
}

func (c *Cluster) forEachMachine(do func(*config.Machine, int) error) error {
	for _, template := range c.spec.Templates {
		for i := 0; i < template.Count; i++ {
			if err := do(&template.Spec, i); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Cluster) createMachine(machine *config.Machine, i int) error {
	name := c.containerName(machine, i)
	runArgs := []string{
		"-it", "-d", "--rm",
		"--name", name,
		"--tmpfs", "/run",
		"--tmpfs", "/tmp",
		"-v", "/sys/fs/cgroup:/sys/fs/cgroup:ro",
	}

	if machine.Privileged {
		runArgs = append(runArgs, "--privileged")
	}

	log.Infof("Creating machine: %s ...", name)
	_, err := docker.Run(machine.Image,
		runArgs,
		[]string{"/sbin/init"},
	)
	return err

}

// Create creates the cluster.
func (c *Cluster) Create() error {
	return c.forEachMachine(c.createMachine)
}

func (c *Cluster) deleteMachine(machine *config.Machine, i int) error {
	return docker.Kill("KILL", c.containerName(machine, i))
}

// Delete deletes the cluster.
func (c *Cluster) Delete() error {
	return c.forEachMachine(c.deleteMachine)
}
