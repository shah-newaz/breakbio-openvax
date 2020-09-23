package executor

import (
	"bufio"
	"breakbio-openvax/pkg/breakbio/log"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type ExecutorInterface interface {
	Execute(writer io.StringWriter, shell string, command string) (exitCode int)
}

type Executor struct {
}

func NewExecutor() *Executor {
	return &Executor{}
}

func (e *Executor) Execute(writer io.StringWriter, shell string, command string) (exitCode int) {
	const failureExitCode = 1
	const defaultExitCode = 0

	cmd := e.createCommand(shell, command)
	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		log.Redf("%v", err)
		return failureExitCode
	}
	cmdErr, err := cmd.StderrPipe()
	if err != nil {
		log.Redf("%v", err)
		return failureExitCode
	}

	err = cmd.Start()
	if err != nil {
		log.Redf("%v", err)
		_, err := writer.WriteString(err.Error())
		if err != nil {
			log.Redf("Failed to write the command output to log file. Message: %+v", err)
		}
		return failureExitCode
	}

	go func() {
		outScanner := bufio.NewScanner(cmdOut)
		for outScanner.Scan() {
			scan := outScanner.Text()
			fmt.Println(scan)
			if len(scan) > 0 {
				logEntry := scan + "\n"
				_, err := writer.WriteString(logEntry)
				if err != nil {
					log.Redf("Failed to write the command output to the log file. Message: %+v", err)
				}
			}
		}
	}()

	go func() {
		errScanner := bufio.NewScanner(cmdErr)
		for errScanner.Scan() {
			scan := errScanner.Text()
			fmt.Println(scan)
			if len(scan) > 0 {
				logEntry := scan + "\n"
				_, err := writer.WriteString(logEntry)
				if err != nil {
					log.Redf("Failed to write the command output to the log file. Message: %+v", err)
				}
			}
		}
	}()

	err = cmd.Wait()
	if err == nil {
		return defaultExitCode
	}

	if exiterr, ok := err.(*exec.ExitError); ok {
		exitCode := exiterr.ProcessState.ExitCode()
		return exitCode
	} else {
		return -1
	}
}

func (e *Executor) createCommand(shell string, command string) *exec.Cmd {
	commandArgument := ""
	if shell == "cmd.exe" {
		commandArgument = "/c"
	} else if shell == "powershell.exe" {
		commandArgument = "-noexit"
	} else if shell == "/bin/sh" || shell == "/bin/bash" {
		commandArgument = "-c"
	}
	if commandArgument == "" {
		cmd := exec.Command(shell, strings.Split(command, " ")...)
		log.Cyanf("Executing command %s %s", shell, command)
		return cmd
	}
	cmd := getCommandExecutor(shell, commandArgument, command)
	log.Cyanf("Executing command %s %s %s", shell, commandArgument, command)
	return cmd
}
