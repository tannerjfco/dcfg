package plugins

import (
	"fmt"
	"log"

	"github.com/drud/dcfg/apptpl"
	"github.com/drud/dcfg/apptpl/drupal"
	"github.com/drud/drud-go/utils/prettify"
)

// // AppMap is used to retrieve the correct plugin
// var AppMap = map[string]apptpl.App{
// 	"drupal":    &drupal.DrupalConfig{},
// 	"wordpress": &wordpress.WordpressConfig{},
// }

// String prints the Task
func (t apptpl.Template) String() string {
	return prettify.Prettify(t)
}

// Run creates configurations for an app to deploy within the container.
func (t apptpl.Template) Run() error {
	if t.App == "" {
		return fmt.Errorf("No app specified")
	}

	newapp := apptpl.AppMap[t.App]

	if t.App == "drupal" {
		log.Printf("Creating configurations for Drupal")
		if t.ConfigPath == "" {
			t.ConfigPath = "sites/default/"
		}
		confFile := drupal.NewDrupalConfig()
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
