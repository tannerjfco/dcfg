package drudconfig

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

// GetConfigGroups unmarshalls config groups from the config file into structs
func GetConfigGroups(confByte []byte) (GroupSet, error) {
	var groups GroupSet
	err := yaml.Unmarshal(confByte, &groups)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

// RunCommand runs aribtrary commands
func RunCommand(command string) error {
	cmdParts := strings.Split(command, " ")

	proc := exec.Command(
		cmdParts[0],
		cmdParts[1:]...,
	)
	proc.Stdout = os.Stdout
	proc.Stdin = os.Stdin
	proc.Stderr = os.Stderr

	err := proc.Run()

	return err
}

// maybeChdir changes to a directory if there is one
func maybeChdir(d string) {
	if d != "" {
		err := os.Chdir(d)
		if err != nil {
			log.Fatal(err)
		}
	}
}
