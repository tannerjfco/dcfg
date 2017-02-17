package tpl

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/drud/drud-go/utils/system"
	"github.com/stretchr/testify/assert"
)

var wpconf = "wp-config.php"
var wapp = "wordpress"

func TestWordPressWriteConfig(t *testing.T) {
	assert := assert.New(t)

	in := new(Config)
	in.App = wapp
	in.IgnoreFiles = true
	wp := new(WordpressConfig)

	err := wp.WriteConfig(in)
	assert.NoError(err)

	content, err := ioutil.ReadFile(wpconf)
	assert.NoError(err)
	assert.Contains(string(content), "define( 'WP_CONTENT_DIR', dirname( __FILE__ ) . '/wp-content' );")
	os.Remove(wpconf)
}

func TestWordPressPlaceFiles(t *testing.T) {
	assert := assert.New(t)

	src := os.TempDir() + "file_src"
	dest := "wp-content"
	os.Setenv("FILE_SRC", src)
	os.Mkdir(dest, 0755)
	os.MkdirAll(src, 0755)
	os.Create(src + "/testfile")

	in := new(Config)
	in.App = wapp
	wp := new(WordpressConfig)

	err := wp.PlaceFiles(in, false)
	assert.NoError(err)
	assert.True(system.FileExists(dest + "/uploads"))
	assert.True(system.FileExists(dest + "/uploads/testfile"))
	link, err := os.Readlink(dest + "/uploads")
	assert.NoError(err)
	assert.Contains(link, "file_src")
	os.Remove(wpconf)
	os.Remove(dest)

	err = wp.PlaceFiles(in, true)
	assert.NoError(err)
	assert.True(system.FileExists(dest + "/uploads"))
	assert.True(system.FileExists(dest + "/uploads/testfile"))
	os.Remove(wpconf)
	os.RemoveAll(src)
	os.RemoveAll(dest)
	os.Unsetenv("FILE_SRC")
}

func TestWordPressWebConfig(t *testing.T) {
	assert := assert.New(t)

	in := new(Config)
	in.App = wapp
	in.IgnoreFiles = true
	in.DocRoot = "potato"
	wp := new(WordpressConfig)
	webConf := "test.conf"
	ioutil.WriteFile(webConf, []byte(testConf), os.FileMode(0644))
	os.Setenv("NGINX_SITE_CONF", webConf)

	err := wp.WebConfig(in)
	assert.NoError(err)
	result, err := ioutil.ReadFile(webConf)
	assert.NoError(err)
	assert.Contains(string(result), "root /var/www/html/potato;")
	os.Remove(webConf)
	os.Remove(wpconf)
}
