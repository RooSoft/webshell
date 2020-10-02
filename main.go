package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
)

var pass *string

type command struct {
	Pass string
	Cmd  string
	Opt  string
	Args string
}

func decodeRequest(r *http.Request) (command, error) {
	var req command

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		return req, errors.New("error parsing json payload")
	}

	return req, nil
}

func verifyPassword(givenPassword string, requiredPassword string) error {
	if givenPassword != requiredPassword {
		return errors.New("wrong password")
	}

	return nil
}

func executeCommand(command, options, arguments string, w io.Writer) error {
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

func handler(w http.ResponseWriter, r *http.Request) {
	command, err := decodeRequest(r)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = verifyPassword(command.Pass, *pass)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = executeCommand(command.Cmd, command.Opt, command.Args, w)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func main() {
	pass = flag.String("pass", "mypass", "passwd")
	addr := flag.String("addr", ":9090", "bind addr and port")

	flag.Parse()

	http.HandleFunc("/", handler)

	log.Printf("start http server\n")
	log.Fatal(http.ListenAndServe(*addr, nil))
}

// curl -X POST -d '{"pass": "mypass", "cmd": "bash", "opt": "-c", "args": "ls -l ~; echo hello"}' http://localhost:9090
