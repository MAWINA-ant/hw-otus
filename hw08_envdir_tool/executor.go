package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	for k, v := range env {
		if _, exist := os.LookupEnv(k); exist {
			if v.NeedRemove {
				os.Unsetenv(k)
				continue
			}
			os.Setenv(k, v.Value)
		} else {
			if !v.NeedRemove {
				os.Setenv(k, v.Value)
			}
		}
	}
	command.Env = os.Environ()
	command.Run()
	return command.ProcessState.ExitCode()
}
