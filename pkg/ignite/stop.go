package ignite

import (
	"github.com/weaveworks/footloose/pkg/exec"
)

// Stop stops a vm.
func Stop(vmname string) error {
	return exec.CommandWithLogging(execName, "stop", vmname)
}
