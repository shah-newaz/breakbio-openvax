// +build windows

package executor

import (
	"os"
	"os/exec"
	"syscall"
)

func getCommandExecutor(shell string, commandArgument string, command string) *exec.Cmd {
	cmd := exec.Command(shell)
	cmd.Env = os.Environ()
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CmdLine: commandArgument + " " + command,
	}
	return cmd
}
