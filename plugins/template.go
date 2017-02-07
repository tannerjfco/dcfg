package plugins

import (
	"fmt"
	"log"

	"github.com/drud/dcfg/apptpl"
	"github.com/drud/drud-go/utils"
)

// Template implements the Template Action
type Template struct {
	TaskDefaults
	App           string `yaml:"app"`
	Core          string `yaml:"core"`
	ConfigPath    string `yaml:"configPath"`
	VHost         string `yaml:"vHost"`
	DBPort        int    `yaml:"dbPort"`
	DBPrefix      string `yaml:"dbPrefix"`
	PublicFiles   string `yaml:"publicFiles"`
	PrivateFiles  string `yaml:"privateFiles"`
	ConfigSyncDir string `yaml:"configSyncDir"`
}

// String prints the Task
func (t Template) String() string {
	return utils.Prettify(t)
}

// Run creates configurations for an app to deploy within the container.
func (t Template) Run() error {
	if t.App == "" {
		return fmt.Errorf("No template specified")
	}

	if t.App == "drupal" {
		log.Printf("Creating configurations for Drupal")
		if t.ConfigPath == "" {
			t.ConfigPath = "sites/default/"
		}
		confFile := apptpl.NewDrupalConfig()
		confFile.DatabasePort = t.DBPort
		confFile.DatabasePrefix = t.DBPrefix
		if t.Core == "8.x" {
			confFile.IsDrupal8 = true
		}
		err := apptpl.WriteDrupalConfig(confFile, t.ConfigPath+"settings.php")
		if err != nil {
			return err
		}
	}

	if t.App == "wordpress" {
		log.Printf("Creating configurations for WordPress")
		confFile := apptpl.NewWordpressConfig()

		err := apptpl.WriteWordpressConfig(confFile, t.ConfigPath+"wp-config.php")
		if err != nil {
			return err
		}
	}

	return nil
}
