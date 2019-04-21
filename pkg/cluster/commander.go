package cluster

import (
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/kind/pkg/docker"
	"sigs.k8s.io/kind/pkg/exec"
)

// GetCommanderInstance returns CommanderSingleton
func GetCommanderInstance() CommanderSingleton {
	if instance == nil {
		instance = new(commanderSingleton)
	}
	return instance
}

// CommanderSingleton provide way to manage command line call verbosity
type CommanderSingleton interface {
	SetVerbosity(vFlag bool)
	RunCommand(cmd string, args ...string) error
	RunCommandContainer(container, cmd string, args ...string) error
	DockerInspect(container, format string) ([]string, error)
	DockerStop(container string) error
	DockerStart(container string) error
	DockerRun(image string, runArgs []string, containerArgs []string) (string, error)
	DockerRm(container string) error
	DockerPullIfNotPresent(image string, retries int) (bool, error)
	DockerKill(container, signal string) error
}

var instance *commanderSingleton

type commanderSingleton struct {
	isVerbose bool
}

// SetVerbosity stores if running commands have to be logged or not
func (v *commanderSingleton) SetVerbosity(verboseFlag bool) {
	v.isVerbose = verboseFlag
}

// RunCommand
func (v *commanderSingleton) RunCommand(name string, args ...string) error {
	if v.isVerbose {
		cmdMessage := make([]interface{}, len(args)+2)
		cmdMessage[0] = "Running command :"
		cmdMessage[1] = name
		for i, v := range args {
			cmdMessage[i+2] = v
		}
		log.Println(cmdMessage...)
	}

	cmd := exec.Command(name, args...)
	output, err := exec.CombinedOutputLines(cmd)
	if err != nil {
		// log error output if there was any
		for _, line := range output {
			log.Error(line)
		}
	}
	return err
}

func (v *commanderSingleton) RunCommandContainer(container, name string, args ...string) error {
	if v.isVerbose {
		cmdMessage := make([]interface{}, len(args)+3)
		cmdMessage[0] = "Running container command :"
		cmdMessage[1] = container
		cmdMessage[2] = name
		for i, v := range args {
			cmdMessage[i+3] = v
		}
		log.Println(cmdMessage...)
	}

	exe := docker.ContainerCmder(container)
	cmd := exe.Command(name, args...)
	output, err := exec.CombinedOutputLines(cmd)
	if err != nil {
		// log error output if there was any
		for _, line := range output {
			log.WithField("machine", container).Error(line)
		}
	}
	return err
}

func (v *commanderSingleton) DockerRun(image string, runArgs []string, containerArgs []string) (string, error) {
	if v.isVerbose {
		args := []string{"docker", "run"}
		args = append(args, runArgs...)
		args = append(args, image)
		args = append(args, containerArgs...)
		cmdMessage := make([]interface{}, len(args)+1)
		cmdMessage[0] = "Running docker command :"
		for i, v := range args {
			cmdMessage[i+1] = v
		}
		log.Println(cmdMessage...)
	}
	return docker.Run(image,
		runArgs,
		containerArgs,
	)
}

func (v *commanderSingleton) DockerInspect(container, format string) ([]string, error) {
	if v.isVerbose {
		log.Println("Inspect docker command : docker inspect -f", format, container)
	}
	return docker.Inspect(container, format)
}

func (v *commanderSingleton) DockerStop(container string) error {
	if v.isVerbose {
		log.Println("Stop docker command: docker stop", container)
	}
	return run("docker", "stop", container)
}

func (v *commanderSingleton) DockerRm(container string) error {
	if v.isVerbose {
		log.Println("Rm docker command: docker rm", container)
	}
	return run("docker", "rm", container)
}

func (v *commanderSingleton) DockerStart(container string) error {
	if v.isVerbose {
		log.Println("Start docker command: docker start", container)
	}
	return run("docker", "start", container)
}

func (v *commanderSingleton) DockerPullIfNotPresent(image string, retries int) (bool, error) {
	if v.isVerbose {
		log.Println("PullIfNotPresent docker command: docker inspect --type=image", image)
	}
	pulled, err := docker.PullIfNotPresent(image, retries)
	if pulled == true && v.isVerbose == true {
		log.Println("Pull docker command: docker pull ", image)
	}
	return pulled, err
}

func (v *commanderSingleton) DockerKill(container, signal string) error {
	if v.isVerbose {
		log.Println("Kill docker command: docker kill -s", signal, container)
	}
	return docker.Kill("KILL", container)
}
