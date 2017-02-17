package tpl

import (
	"testing"

	"os"

	"io/ioutil"

	"github.com/drud/drud-go/utils/system"
	"github.com/stretchr/testify/assert"
)

const (
	testConf = `
	server {
		listen 80; ## listen for ipv4; this line is default and implied
		listen [::]:80 default ipv6only=on; ## listen for ipv6
		root /var/www/html;
	}`
)

var conf = "settings.php"
var app = "drupal"

func TestDrupalWriteConfig(t *testing.T) {
	assert := assert.New(t)
	path := os.TempDir()

	in := new(Config)
	in.App = app
	in.IgnoreFiles = true
	in.ConfigPath = path
	drupal := new(DrupalConfig)

	err := drupal.WriteConfig(in)
	assert.NoError(err)

	content, err := ioutil.ReadFile(path + conf)
	assert.NoError(err)
	assert.Contains(string(content), "'database' => \"data\"")
	os.Remove(conf)
}

func TestDrupalPlaceFiles(t *testing.T) {
	assert := assert.New(t)

	src := os.TempDir() + "file_src"
	dest := "sites/default"
	os.Setenv("FILE_SRC", src)
	os.MkdirAll(dest, 0755)
	os.MkdirAll(src, 0755)
	os.Create(src + "/testfile")

	in := new(Config)
	in.App = app
	drupal := new(DrupalConfig)

	err := drupal.PlaceFiles(in, false)
	assert.NoError(err)
	assert.True(system.FileExists(dest + "/files"))
	assert.True(system.FileExists(dest + "/files/testfile"))
	link, err := os.Readlink(dest + "/files")
	assert.NoError(err)
	assert.Contains(link, "file_src")
	os.Remove(conf)
	os.Remove(dest)

	err = drupal.PlaceFiles(in, true)
	assert.NoError(err)
	assert.True(system.FileExists(dest + "/files"))
	assert.True(system.FileExists(dest + "/files/testfile"))
	os.Remove(conf)
	os.RemoveAll(src)
	os.RemoveAll(dest)
	os.Unsetenv("FILE_SRC")
}

func TestDrupalWebConfig(t *testing.T) {
	assert := assert.New(t)

	in := new(Config)
	in.App = app
	in.IgnoreFiles = true
	in.DocRoot = "potato"
	drupal := new(DrupalConfig)
	webConf := "test.conf"
	ioutil.WriteFile(webConf, []byte(testConf), os.FileMode(0644))
	os.Setenv("NGINX_SITE_CONF", webConf)

	err := drupal.WebConfig(in)
	assert.NoError(err)
	result, err := ioutil.ReadFile(webConf)
	assert.NoError(err)
	assert.Contains(string(result), "root /var/www/html/potato;")
	os.Remove(webConf)
	os.Remove(conf)
}
