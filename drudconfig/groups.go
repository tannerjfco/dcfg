package drudconfig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ConfigGroup models the config group from the drud.yaml
type ConfigGroup struct {
	Name    string             `yaml:"name",json:"name"`
	Env     map[string]string  `yaml:"env",json:"env"`
	User    string             `yaml:"user",json:"user"`
	Workdir string             `yaml:"workdir",json:"workdir"`
	Tasks   []*json.RawMessage `json:"tasks",yaml:"tasks"`
}

// Run does its best to execute the cmd defined by the user
func (g *ConfigGroup) Run() error {
	baseDir, _ := os.Getwd()
	var workDir string
	if g.Workdir != "" {
		workDir, _ = filepath.Abs(g.Workdir)
	}
	maybeChdir(workDir)

	for k, v := range g.Env {
		if strings.HasPrefix(v, "$") && v == strings.ToUpper(v) {
			g.Env[k] = os.Getenv(v[1:])
		}
	}

	fmt.Println("Runnign group:", g.Name)
	//fmt.Println(string(g.Tasks))

	for _, t := range g.Tasks {
		taskString := string([]byte(*t))
		if HasVars(taskString) {
			var doc bytes.Buffer
			templ := template.New("cmd template")
			templ, _ = templ.Parse(taskString)
			templ.Execute(&doc, g.Env)
			taskString = doc.String()
		}

		var cmdType TaskType
		err := json.Unmarshal([]byte(taskString), &cmdType)
		if err != nil {
			fmt.Println(err)
		}

		var action Action
		switch cmdType.Action {
		case "command":
			var cmd Command
			err = json.Unmarshal([]byte(taskString), &cmd)
			if err != nil {
				fmt.Println(err)
			}
			action = cmd
		case "write":
			var w Write
			err = json.Unmarshal([]byte(taskString), &w)
			if err != nil {
				fmt.Println(err)
			}
			action = w
		default:
			log.Fatalf("unknown command type: %q", cmdType.Action)
		}

		action.Pretty()

	}

	os.Chdir(baseDir)

	return nil
}

// GroupSet represents a lsit of ConfigGroups
type GroupSet []ConfigGroup
