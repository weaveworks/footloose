package cluster

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dlespiau/footloose/pkg/config"
	"github.com/ghodss/yaml"
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
	spec config.Config
}

// New creates a new cluster. It takes as input the description of the cluster
// and its machines.
func New(conf config.Config) *Cluster {
	return &Cluster{
		spec: conf,
	}
}

// NewFromFile creates a new Cluster from a YAML serialization of its
// configuration.
func NewFromFile(path string) (*Cluster, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	spec := config.Config{}
	err = yaml.Unmarshal(data, &spec)
	return New(spec), err
}

// Save writes the Cluster configure to a file.
func (c *Cluster) Save(path string) error {
	data, err := yaml.Marshal(c.spec)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0666)
}

func f(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

func (c *Cluster) containerName(machine *config.Machine, i int) string {
	format := "%s-" + machine.Name
	return f(format, c.spec.Cluster.Name, i)
}

func (c *Cluster) forEachMachine(do func(*config.Machine, int) error) error {
	for _, template := range c.spec.Machines {
		for i := 0; i < template.Count; i++ {
			if err := do(&template.Spec, i); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Cluster) ensureSSHKey() error {
	path := c.spec.Cluster.PrivateKey
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	log.Infof("Creating SSH key: %s ...", path)
	return run(
		"ssh-keygen", "-q",
		"-t", "rsa",
		"-b", "4096",
		"-C", f("%s@footloose.mail", c.spec.Cluster.Name),
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
	return ioutil.ReadFile(c.spec.Cluster.PrivateKey + ".pub")
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
	name := c.containerName(machine, i)
	log.Infof("Deleting machine: %s ...", name)
	return docker.Kill("KILL", name)
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
	ip, err := containerIP(f("%s-%s", c.spec.Cluster.Name, name))
	if err != nil {
		return err
	}
	cmd := exec.Command(
		"ssh", "-q",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "StrictHostKeyChecking=no",
		"-i", c.spec.Cluster.PrivateKey,
		f("%s@%s", "root", ip),
	)
	cmd.SetStdin(os.Stdin)
	exec.InheritOutput(cmd)
	return cmd.Run()
}
