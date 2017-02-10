package tpl

import (
	"os"
	"text/template"

	"github.com/Masterminds/sprig"
	"gopkg.in/oleiade/reflections.v1"
)

// DrupalConfig encapsulates all the configurations for a Drupal site.
type DrupalConfig struct {
	Core             string
	DeployURL        string
	ConfigSyncDir    string
	DatabaseName     string
	DatabaseUsername string
	DatabasePassword string
	DatabaseHost     string
	DatabaseDriver   string
	DatabasePort     int
	DatabasePrefix   string
	HashSalt         string
	Hostname         string
	IsDrupal8        bool
}

// NewDrupalConfig initializes a DrupalConfig object with defaults
func NewDrupalConfig() *DrupalConfig {
	return &DrupalConfig{
		Core:             "7.x",
		ConfigSyncDir:    "/var/www/html/sync",
		DatabaseName:     "data",
		DatabaseUsername: "root",
		DatabasePassword: "root",
		DatabaseHost:     "127.0.0.1",
		DatabaseDriver:   "mysql",
		DatabasePort:     3306,
		DatabasePrefix:   "",
		IsDrupal8:        false,
	}
}

// WriteConfig produces a valid settings.php file from the defined configurations
func (c *DrupalConfig) WriteConfig(in *Config) error {
	conf := NewDrupalConfig()

	if in.Core == "8.x" {
		conf.IsDrupal8 = true
	}

	srcFieldList, err := reflections.Items(in)
	if err != nil {
		return err
	}

	for field, val := range srcFieldList {
		if val != "" {
			has, err := reflections.HasField(conf, field)
			if err != nil {
				return err
			}
			if has && val != 0 {
				err = reflections.SetField(conf, field, val)
				if err != nil {
					return err
				}
			}
		}
	}

	tmpl, err := template.New("conf").Funcs(sprig.TxtFuncMap()).Parse(drupalTemplate)
	if err != nil {
		return err
	}

	filepath := ""
	if in.ConfigPath != "" {
		filepath = in.ConfigPath
	}

	file, err := os.Create(filepath + "settings.php")
	if err != nil {
		return err
	}

	err = tmpl.Execute(file, conf)
	if err != nil {
		return err
	}

	return nil
}

// PlaceFiles determines where file upload directories should go.
func (c *DrupalConfig) PlaceFiles(in *Config) error {
	return nil
}
