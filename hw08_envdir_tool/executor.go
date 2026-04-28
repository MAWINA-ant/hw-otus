package main

import (
	"fmt"
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
		if v.NeedRemove {
			_, exist := os.LookupEnv("k")
			if exist {
				os.Unsetenv(k)
			}
		} else {
			envString := fmt.Sprintf("%s=%s", k, v.Value)
			command.Env = append(command.Env, envString)
		}
	}
	command.Run()
	return command.ProcessState.ExitCode()
}
