package ignite

import (
	"github.com/weaveworks/footloose/pkg/exec"
)

// Start starts a vm.
func Start(vmname string) error {
	return exec.Command(execName, "start", vmname).Run()
}
