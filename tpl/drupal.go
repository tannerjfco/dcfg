package tpl

import (
	"fmt"
	"os"
	"text/template"

	"path"

	"log"

	"io/ioutil"
	"regexp"

	"github.com/Masterminds/sprig"
	"github.com/drud/drud-go/utils/system"
	"gopkg.in/oleiade/reflections.v1"
)

// DrupalConfig encapsulates all the configurations for a Drupal site.
type DrupalConfig struct {
	Core             string
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
	SiteURL          string
}

// DefaultDrupalConfig initializes a DrupalConfig object with defaults
func DefaultDrupalConfig() *DrupalConfig {
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
		HashSalt:         PassTheSalt(),
		IsDrupal8:        false,
	}
}

// WriteAppConfig produces a valid settings.php file from the defined configurations
func (c *DrupalConfig) WriteAppConfig(in *Config) error {
	c = DefaultDrupalConfig()

	if in.Core == "8.x" {
		c.IsDrupal8 = true
	}

	srcFieldList, err := reflections.Items(in)
	if err != nil {
		return err
	}

	for field, val := range srcFieldList {
		if val != "" {
			has, err := reflections.HasField(c, field)
			if err != nil {
				return err
			}
			if has && val != 0 {
				err = reflections.SetField(c, field, val)
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

	filepath := "sites/default/"
	if in.ConfigPath != "" {
		filepath = SlashIt(in.ConfigPath, true)
	}

	file, err := os.Create(filepath + "settings.php")
	if err != nil {
		return err
	}

	err = tmpl.Execute(file, c)
	if err != nil {
		return err
	}

	return nil
}

// PlaceFiles determines where file upload directories should go.
func (c *DrupalConfig) PlaceFiles(in *Config, move bool) error {
	src := os.Getenv("FILE_SRC")
	dest := "sites/default/files"
	if in.PublicFiles != "" {
		dest = in.PublicFiles
	}

	if src == "" || !system.FileExists(src) {
		log.Fatalf("source path for files does not exist")
	}

	if system.FileExists(dest) {
		log.Printf("destination path exists, removing")
		os.RemoveAll(dest)
	}

	// ensure parent of destination is writable
	os.Chmod(path.Dir(src), 0755)

	if move {
		out, err := system.RunCommand(
			"rsync",
			[]string{"-avz", "--recursive", src + "/", dest},
		)
		if err != nil {
			return fmt.Errorf("%s - %s", err.Error(), string(out))
		}
	} else {
		err := os.Symlink(src, dest)
		if err != nil {
			return err
		}
	}

	return nil
}

// WriteWebConfig updates the web server configuration to support the provided app configurations
// @TODO: need to update rules for public/private files, holding off until more firm on approach for this task.
func (c *DrupalConfig) WriteWebConfig(in *Config) error {
	dest := os.Getenv("NGINX_SITE_CONF")
	root := "root /var/www/html"

	if !system.FileExists(dest) {
		log.Fatalf("target %s does not exist", dest)
	}

	conf, err := ioutil.ReadFile(dest)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(root)
	new := re.ReplaceAllString(string(conf), root+"/"+in.DocRoot)

	err = ioutil.WriteFile(dest, []byte(new), 0644)
	if err != nil {
		return err
	}

	return nil
}
