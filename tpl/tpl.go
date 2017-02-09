package tpl

import (
	"log"

	"github.com/drud/drud-go/utils/prettify"
)

// Config implements the Template Action
type Config struct {
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

// Tpl is the interface that each plugin must implement
type Tpl interface {
	WriteConfig(in *Config) error
}

// TplMap is used to retrieve the correct plugin
var TplMap = map[string]Tpl{
	"drupal":    &DrupalConfig{},
	"wordpress": &WordpressConfig{},
}

// String prints the Task
func (c Config) String() string {
	return prettify.Prettify(c)
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
