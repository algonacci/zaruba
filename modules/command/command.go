package command

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/state-alchemists/zaruba/modules/config"
	"github.com/state-alchemists/zaruba/modules/logger"
)

// GetCmd get cmd object
func GetCmd(dir, command string, args ...string) (cmd *exec.Cmd, err error) {
	cmd = exec.Command(command, args...)
	cmd.Env = os.Environ()
	cmd.Dir, err = filepath.Abs(dir)
	return cmd, err
}

// GetShellCmd get cmd object for running shell command
func GetShellCmd(dir, script string) (cmd *exec.Cmd, err error) {
	shell, shellArg := config.GetShellAndShellArg()
	args := []string{shellArg, script}
	return GetCmd(dir, shell, args...)
}

// runCmdAndReturn run cmd object
func runCmdAndReturn(cmd *exec.Cmd) (output string, err error) {
	logger.Info("[%s] %s", cmd.Dir, strings.Join(cmd.Args, " "))
	outputB, err := cmd.Output()
	if err != nil {
		return output, err
	}
	output = string(outputB)
	return output, err
}

// RunShellScriptAndReturn run shell script with custom env
func RunShellScriptAndReturn(dir, script string, env []string) (output string, err error) {
	shell, shellArg := config.GetShellAndShellArg()
	args := []string{shellArg, script}
	cmd, err := GetCmd(dir, shell, args...)
	if err != nil {
		return output, err
	}
	cmd.Stderr = os.Stderr
	if len(env) > 0 {
		cmd.Env = env
	}
	return runCmdAndReturn(cmd)
}

// RunAndReturn run command
func RunAndReturn(dir, command string, args ...string) (output string, err error) {
	cmd, err := GetCmd(dir, command, args...)
	if err != nil {
		return output, err
	}
	cmd.Stderr = os.Stderr
	return runCmdAndReturn(cmd)
}

// RunInteractively run command
func RunInteractively(dir, command string, args ...string) (err error) {
	cmd, err := GetCmd(dir, command, args...)
	if err != nil {
		return err
	}
	logger.Info("[%s] %s", cmd.Dir, strings.Join(cmd.Args, " "))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	return err
}
