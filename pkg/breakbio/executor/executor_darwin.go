// +build darwin

package executor

import "os"
import "os/exec"

func getCommandExecutor(shell string, commandArgument string, command string) *exec.Cmd {
	cmd := exec.Command(shell, commandArgument, command)
	cmd.Env = os.Environ()
	return cmd
}
