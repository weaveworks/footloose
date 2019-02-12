package cluster

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"time"

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

func (c *Cluster) machine(spec *config.Machine, i int) *Machine {
	return &Machine{
		spec:     spec,
		name:     c.containerName(spec, i),
		hostname: f(spec.Name, i),
	}

}

func (c *Cluster) forEachMachine(do func(*Machine, int) error) error {
	for _, template := range c.spec.Machines {
		for i := 0; i < template.Count; i++ {
			machine := c.machine(&template.Spec, i)
			if err := do(machine, i); err != nil {
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

func (c *Cluster) createMachine(machine *Machine, i int) error {
	name := machine.ContainerName()
	runArgs := []string{
		"-it", "-d", "--rm",
		"--name", name,
		"--hostname", machine.Hostname(),
		"--tmpfs", "/run",
		"--tmpfs", "/tmp",
		"-v", "/sys/fs/cgroup:/sys/fs/cgroup:ro",
	}

	for _, volume := range machine.spec.Volumes {
		mount := f("type=%s", volume.Type)
		if volume.Source != "" {
			mount += f(",src=%s", volume.Source)
		}
		mount += f(",dst=%s", volume.Destination)
		if volume.ReadOnly {
			mount += ",readonly"
		}
		runArgs = append(runArgs, "--mount", mount)
	}

	for _, mapping := range machine.spec.PortMappings {
		publish := ""
		if mapping.Address != "" {
			publish += f("%s:", mapping.Address)
		}
		publish += f("%d", mapping.ContainerPort)
		if mapping.HostPort != 0 {
			publish += f(":%d", int(mapping.HostPort)+i)
		}
		if mapping.Protocol != "" {
			publish += f("/%s", mapping.Protocol)
		}
		runArgs = append(runArgs, "-p", publish)
	}

	if machine.spec.Privileged {
		runArgs = append(runArgs, "--privileged")
	}

	// Start the container.
	log.Infof("Creating machine: %s ...", name)
	_, err := docker.Run(machine.spec.Image,
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
	for _, template := range c.spec.Machines {
		if _, err := docker.PullIfNotPresent(template.Spec.Image, 2); err != nil {
			return err
		}
	}
	return c.forEachMachine(c.createMachine)
}

func (c *Cluster) deleteMachine(machine *Machine, i int) error {
	name := machine.ContainerName()
	log.Infof("Deleting machine: %s ...", name)
	return docker.Kill("KILL", name)
}

// Delete deletes the cluster.
func (c *Cluster) Delete() error {
	return c.forEachMachine(c.deleteMachine)
}

// io.Writer filter that writes that it receives to writer. Keeps track if it
// has seen a write matching regexp.
type matchFilter struct {
	writer       io.Writer
	writeMatched bool // whether the filter should write the matched value or not.

	regexp  *regexp.Regexp
	matched bool
}

func (f *matchFilter) Write(p []byte) (n int, err error) {
	// Assume the relevant log line is flushed in one write.
	if match := f.regexp.Match(p); match {
		f.matched = true
		if !f.writeMatched {
			return len(p), nil
		}
	}
	return f.writer.Write(p)
}

// Matches:
//   ssh: connect to host 172.17.0.2 port 22: Connection refused
var connectRefused = regexp.MustCompile("ssh: connect to host .* port [0-9]{1,5}: Connection refused")

// Matches:
//   Warning:Permanently added '172.17.0.2' (ECDSA) to the list of known hosts
var knownHosts = regexp.MustCompile("Warning: Permanently added .* to the list of known hosts.")

// ssh returns true if the command should be tried again.
func ssh(args []string) (bool, error) {
	cmd := exec.Command("ssh", args...)

	refusedFilter := &matchFilter{
		writer:       os.Stderr,
		writeMatched: false,
		regexp:       connectRefused,
	}

	errFilter := &matchFilter{
		writer:       refusedFilter,
		writeMatched: false,
		regexp:       knownHosts,
	}

	cmd.SetStdin(os.Stdin)
	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(errFilter)

	err := cmd.Run()
	if err != nil && refusedFilter.matched {
		return true, err
	}
	return false, err
}

func (c *Cluster) machineFromHostname(hostname string) (*Machine, error) {
	for _, template := range c.spec.Machines {
		for i := 0; i < template.Count; i++ {
			if hostname == f(template.Spec.Name, i) {
				return c.machine(&template.Spec, i), nil
			}
		}
	}
	return nil, fmt.Errorf("%s: invalid machine hostname", hostname)
}

func mappingFromPort(spec *config.Machine, containerPort int) (*config.PortMapping, error) {
	for i := range spec.PortMappings {
		if int(spec.PortMappings[i].ContainerPort) == containerPort {
			return &spec.PortMappings[i], nil
		}
	}
	return nil, fmt.Errorf("unknown containerPort %d", containerPort)
}

// SSH logs into the name machine with SSH.
func (c *Cluster) SSH(name string, remoteArgs ...string) error {
	machine, err := c.machineFromHostname(name)
	if err != nil {
		return err
	}
	hostPort, err := machine.HostPort(22)
	if err != nil {
		return err
	}
	mapping, err := mappingFromPort(machine.spec, 22)
	if err != nil {
		return err
	}
	remote := "localhost"
	if mapping.Address != "" {
		remote = mapping.Address
	}
	args := []string{
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "StrictHostKeyChecking=no",
		"-i", c.spec.Cluster.PrivateKey,
		"-p", f("%d", hostPort),
		f("%s@%s", "root", remote),
	}
	args = append(args, remoteArgs...)
	// If we ssh in a bit too quickly after the container creation, ssh errors out
	// with:
	//   ssh: connect to host 172.17.0.2 port 22: Connection refused
	// Let's loop a few times if we receive this message.
	retries := 3
	var retry bool
	for retries > 0 {
		retry, err = ssh(args)
		if !retry {
			break
		}
		retries--
		time.Sleep(200 * time.Millisecond)
	}
	return err
}
