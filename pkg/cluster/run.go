package cluster

import (
	"bytes"
	"fmt"
)

// run runs a command. It will output the combined stdout/error on failure.
func run(name string, args ...string) error {
	return GetCommanderInstance().RunCommand(name, args...)
}

// Run a command in a container. It will output the combined stdout/error on failure.
func containerRun(nameOrID string, name string, args ...string) error {
	return GetCommanderInstance().RunCommandContainer(nameOrID, name, args...)
}

func containerRunShell(nameOrID string, script string) error {
	return containerRun(nameOrID, "/bin/bash", "-c", script)
}

func copy(nameOrID string, content []byte, path string) error {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("cat <<__EOF | tee -a %s\n", path))
	buf.Write(content)
	buf.WriteString("__EOF")
	return containerRunShell(nameOrID, buf.String())
}
