package tpl

import (
	"log"

	"github.com/drud/drud-go/utils/pretty"
)

// Config implements the Template Action
type Config struct {
	App              string `yaml:"app"`
	Core             string `yaml:"core"`
	ConfigPath       string `yaml:"configPath"`
	DocRoot          string `yaml:"docroot"`
	DatabaseName     string `yaml:"dbName"`
	DatabaseUsername string `yaml:"dbUser"`
	DatabasePassword string `yaml:"dbPass"`
	DatabaseHost     string `yaml:"dbHost"`
	DatabaseDriver   string `yaml:"dbDriver"`
	DatabasePort     int    `yaml:"dbPort"`
	DatabasePrefix   string `yaml:"dbPrefix"`
	PublicFiles      string `yaml:"publicFiles"`
	PrivateFiles     string `yaml:"privateFiles"`
	ConfigSyncDir    string `yaml:"configSyncDir"`
	DeployURL        string `yaml:"deployURL"`
}

// Tpl is the interface that each plugin must implement
type Tpl interface {
	WriteConfig(in *Config) error
	PlaceFiles(in *Config) error
}

// TplMap is used to retrieve the correct plugin
var TplMap = map[string]Tpl{
	"drupal":    &DrupalConfig{},
	"wordpress": &WordpressConfig{},
}

// String prints the Task
func (c Config) String() string {
	return pretty.Prettify(c)
}

// Run creates configurations for an application
func (c *Config) Run() error {
	log.Printf("this is a %s app", c.App)

	app := TplMap[c.App]

	err := app.WriteConfig(c)
	if err != nil {
		return err
	}

	return nil
}
