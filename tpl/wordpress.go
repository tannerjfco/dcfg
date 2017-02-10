package tpl

import (
	"os"
	"text/template"

	"github.com/Masterminds/sprig"
)

// WordpressConfig encapsulates all the configurations for a Wordpress site.
type WordpressConfig struct {
	WPGeneric        bool
	DeployName       string
	DeployProtocol   string
	DeployURL        string
	DatabaseName     string
	DatabaseUsername string
	DatabasePassword string
	DatabaseHost     string
	AuthKey          string
	SecureAuthKey    string
	LoggedInKey      string
	NonceKey         string
	AuthSalt         string
	SecureAuthSalt   string
	LoggedInSalt     string
	NonceSalt        string
	Docroot          string
	TablePrefix      string
}

// NewWordpressConfig produces a WordpressConfig object with defaults.
func NewWordpressConfig() *WordpressConfig {
	return &WordpressConfig{
		WPGeneric:        false,
		DatabaseName:     "data",
		DatabaseUsername: "root",
		DatabasePassword: "root",
		DatabaseHost:     "127.0.0.1",
		Docroot:          "/var/www/html/docroot",
		TablePrefix:      "wp_",
		DeployURL:        os.Getenv("DEPLOY_URL"),
		DeployProtocol:   os.Getenv("DEPLOY_PROTOCOL"),
	}
}

// WriteConfig produces a valid settings.php file from the defined configurations
func (c *WordpressConfig) WriteConfig(in *Config) error {
	conf := NewWordpressConfig()
	conf.TablePrefix = in.DatabasePrefix

	tmpl, err := template.New("conf").Funcs(sprig.TxtFuncMap()).Parse(wordpressTemplate)
	if err != nil {
		return err
	}

	filepath := ""
	if in.ConfigPath != "" {
		filepath = in.ConfigPath
	}

	file, err := os.Create(filepath + "wp-config.php")
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
func (c *WordpressConfig) PlaceFiles(in *Config) error {
	return nil
}
