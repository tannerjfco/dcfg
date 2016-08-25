package drudconfig

import (
	"fmt"
	"os"
	"path/filepath"
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
