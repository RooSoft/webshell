package shell

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
)

func ExecuteCommand(command, options, arguments string, w io.Writer) error {
	cmd := exec.Command(command, options, arguments)

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var errStdout, errStderr error

	err := cmd.Start()

	if err != nil {
		message := fmt.Sprintf("cmd.Start() failed with '%s'\n", err)
		return errors.New(message)
	}

	go func() {
		_, errStdout = io.Copy(w, stdoutIn)
	}()

	go func() {
		_, errStderr = io.Copy(w, stderrIn)
	}()

	err = cmd.Wait()
	if err != nil {
		message := fmt.Sprintf("cmd.Run() failed with %s\n", err)
		return errors.New(message)
	}

	if errStdout != nil || errStderr != nil {
		message := fmt.Sprintf("failed to capture stdout or stderr\n")
		return errors.New(message)
	}

	return nil
}
