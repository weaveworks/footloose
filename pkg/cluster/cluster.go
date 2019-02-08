package cluster

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dlespiau/footloose/pkg/config"
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/kind/pkg/docker"
	"sigs.k8s.io/kind/pkg/exec"
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

func (c *Cluster) ensureSSHKey() error {
	path := c.spec.PrivateKey
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	log.Infof("Creating SSH key: %s ...", path)
	return run(
		"ssh-keygen", "-q",
		"-t", "rsa",
		"-b", "4096",
		"-C", f("%s@footloose.mail", c.spec.Name),
		"-f", path,
		"-N", "",
	)
}

const initScript = `
set -e
rm -f /run/nologin
sshdir=/root/.ssh
mkdir $sshdir; chmod 700 $sshdir
touch $sshdir/authorized_keys; chmod 600 $sshdir/authorized_keys
`

func (c *Cluster) publicKey() ([]byte, error) {
	return ioutil.ReadFile(c.spec.PrivateKey + ".pub")
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

	// Start the container.
	log.Infof("Creating machine: %s ...", name)
	_, err := docker.Run(machine.Image,
		runArgs,
		[]string{"/sbin/init"},
	)
	if err != nil {
		return err
	}

	// Initial provisioning.
	if err := containerRunShell(name, initScript); err != nil {
		return err
	}
	publicKey, err := c.publicKey()
	if err != nil {
		return err
	}
	if err := copy(name, publicKey, "/root/.ssh/authorized_keys"); err != nil {
		return err
	}

	return nil
}

// Create creates the cluster.
func (c *Cluster) Create() error {
	if err := c.ensureSSHKey(); err != nil {
		return err
	}
	return c.forEachMachine(c.createMachine)
}

func (c *Cluster) deleteMachine(machine *config.Machine, i int) error {
	return docker.Kill("KILL", c.containerName(machine, i))
}

// Delete deletes the cluster.
func (c *Cluster) Delete() error {
	return c.forEachMachine(c.deleteMachine)
}

func containerIP(nameOrID string) (string, error) {
	output, err := docker.Inspect(nameOrID, "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}")
	if err != nil {
		for _, line := range output {
			log.Error(line)
		}
		return "", err
	}
	if len(output) != 1 {
		return "", fmt.Errorf("expected 1 IP for %s got %d", nameOrID, len(output))
	}
	return strings.Trim(output[0], "'"), nil
}

// SSH logs into the name machine with SSH.
func (c *Cluster) SSH(name string) error {
	ip, err := containerIP(f("%s-%s", c.spec.Name, name))
	if err != nil {
		return err
	}
	cmd := exec.Command(
		"ssh", "-q",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "StrictHostKeyChecking=no",
		"-i", c.spec.PrivateKey,
		f("%s@%s", "root", ip),
	)
	cmd.SetStdin(os.Stdin)
	exec.InheritOutput(cmd)
	return cmd.Run()
}
