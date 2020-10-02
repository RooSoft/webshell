package shell

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
)

func ExecuteCommand(command, options, arguments string, w io.Writer) error {
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := exec.Command(command, options, arguments)
	// cmd := exec.Command("sh", "-c", "echo stdout; echo 1>&2 stderr")

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var errStdout, errStderr error
	stdout := io.MultiWriter(w, &stdoutBuf)
	stderr := io.MultiWriter(w, &stderrBuf)
	err := cmd.Start()

	if err != nil {
		message := fmt.Sprintf("cmd.Start() failed with '%s'\n", err)
		return errors.New(message)
	}

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()

	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
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
	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	log.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)

	return nil
}
