package drudconfig

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// ConfigGroup models the config group from the drud.yaml
type ConfigGroup struct {
	Name    string            `yaml:"name"`
	Env     map[string]string `yaml:"env"`
	User    string            `yaml:"user"`
	Workdir string            `yaml:"workdir"`
	Tasks   []struct {
		Name    string `yaml:"name"`
		Cmd     string `yaml:"cmd"`
		Workdir string `yaml:"workdir"`
		Wait    string `yaml:"wait"`
		Repeat  int    `yaml:"repeat"`
		Ignore  bool   `yaml:"ignore"`
	} `json:"tasks"`
}

// GroupSet represents a lsit of ConfigGroups
type GroupSet []ConfigGroup

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

	fmt.Println("-", g.Name)
	for i, t := range g.Tasks {

		if t.Name != "" {
			fmt.Println(i, t.Name)
		}

		maybeChdir(t.Workdir)
		for c := t.Repeat; c >= 0; c-- {

			if t.Wait != "" {
				lengthOfWait, _ := time.ParseDuration(t.Wait)
				time.Sleep(lengthOfWait)
			}

			if HasVars(t.Cmd) {
				var doc bytes.Buffer
				templ := template.New("cmd template")
				templ, _ = templ.Parse(t.Cmd)
				templ.Execute(&doc, g.Env)
				t.Cmd = doc.String()
				fmt.Println(t.Cmd)
			}

			err := RunCommand(t.Cmd)
			if err != nil {
				if !t.Ignore {
					return err
				}
			}
		}

		maybeChdir(workDir)
		fmt.Println("---")
	}

	os.Chdir(baseDir)

	return nil
}
