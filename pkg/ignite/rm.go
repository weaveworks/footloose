package ignite

import "github.com/weaveworks/footloose/pkg/exec"

func Remove(name string) error {
	runArgs := []string{
		"rm",
		"-f"
		name,
	}
	return exec.CommandWithLogging(execName, runArgs...)
}
