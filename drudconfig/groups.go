package drudconfig

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
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
		Name    string      `yaml:"name"`
		Cmd     string      `yaml:"cmd"`
		Write   string      `yaml:"write"`
		Dest    string      `yaml:"dest"`
		Mode    os.FileMode `yaml:"mode"`
		Workdir string      `yaml:"workdir"`
		Wait    string      `yaml:"wait"`
		Repeat  int         `yaml:"repeat"`
		Ignore  bool        `yaml:"ignore"`
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

			taskPayload := ""
			if t.Cmd != "" {
				taskPayload = t.Cmd
			} else if t.Write != "" {
				taskPayload = t.Write
			}

			if HasVars(taskPayload) {
				var doc bytes.Buffer
				templ := template.New("cmd template")
				templ, _ = templ.Parse(taskPayload)
				templ.Execute(&doc, g.Env)
				taskPayload = doc.String()
			}

			if t.Cmd != "" {
				err := RunCommand(taskPayload)
				if err != nil {
					if !t.Ignore {
						return err
					}
				}
			} else if t.Write != "" {

				err := ioutil.WriteFile(t.Dest, []byte(taskPayload), t.Mode)
				if err != nil {
					log.Fatalln("Could not read config file:", err)
				}

				info, _ := os.Stat(t.Dest)
				if info.Mode() != t.Mode {
					err := os.Chmod(t.Dest, t.Mode)
					if err != nil {
						log.Fatalln(err)
					}
				}
			}

		}

		maybeChdir(workDir)
		fmt.Println("---")
	}

	os.Chdir(baseDir)

	return nil
}
