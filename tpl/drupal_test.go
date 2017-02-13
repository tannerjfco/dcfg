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

func TestDrupalWriteConfig(t *testing.T) {
	assert := assert.New(t)
	bin := "dcfg"
	confFile := "settings.php"

	// test creating a drupal config w/ no parameters set in yaml
	args := []string{"run", "default", "--config", "testing/drupal.yaml"}
	out, err := system.RunCommand(bin, args)
	assert.NoError(err)
	result, err := ioutil.ReadFile(confFile)
	assert.NoError(err)
	assert.Contains(string(out), "this is a drupal app")
	assert.Contains(string(result), "'database' => \"data\"")
	err = os.Remove(confFile)
	assert.Nil(err)

	// test creating a config for drupal 8
	args[1] = "drupal8"
	out, err = system.RunCommand(bin, args)
	assert.NoError(err)
	assert.Contains(string(out), "Core: \"8.x\",")
	result, err = ioutil.ReadFile(confFile)
	assert.NoError(err)
	assert.Contains(string(result), "$settings['hash_salt'] =")
	err = os.Remove(confFile)
	assert.Nil(err)

	// test creating a config that uses env vars that ddev would set
	args[1] = "ddev_configured"
	os.Setenv("DEPLOY_URL", "http://www.test.site")
	os.Setenv("DB_NAME", "db")
	_, err = system.RunCommand(bin, args)
	assert.NoError(err)
	result, err = ioutil.ReadFile(confFile)
	assert.NoError(err)
	assert.Contains(string(result), "$base_url = 'http://www.test.site';")
	assert.Contains(string(result), "'database' => \"db\"")
	os.Unsetenv("DEPLOY_URL")
	os.Unsetenv("DB_NAME")
	err = os.Remove(confFile)
	assert.Nil(err)

	// test creating a config that handles a strange files directory
	args[1] = "weird_file_dir"
	src := "testing/file_src"
	dest := "potato"
	os.Setenv("FILE_SRC", src)
	os.MkdirAll(src, 0755)
	os.Create(src + "/testfile")
	_, err = system.RunCommand(bin, args)
	assert.NoError(err)
	assert.True(system.FileExists(dest))
	assert.True(system.FileExists(dest + "/testfile"))
	// validate the symlink
	link, err := os.Readlink(dest)
	assert.NoError(err)
	assert.Contains(link, "testing/file_src")
	// reset to test as dir we rsync to
	os.Remove(confFile)
	os.Remove(dest)
	os.Setenv("DEPLOY_NAME", "local")
	_, err = system.RunCommand(bin, args)
	assert.NoError(err)
	assert.True(system.FileExists(dest))
	assert.True(system.FileExists(dest + "/testfile"))
	// cleanup
	os.Remove(confFile)
	os.Unsetenv("FILE_SRC")
	os.Unsetenv("DEPLOY_NAME")
	os.RemoveAll(src)
	os.RemoveAll(dest)

	// test configuring a site contained in a folder called docroot
	args[1] = "have_docroot"
	webConf := "test.conf"
	ioutil.WriteFile(webConf, []byte(testConf), os.FileMode(0644))
	os.Setenv("NGINX_SITE_CONF", webConf)
	_, err = system.RunCommand(bin, args)
	assert.NoError(err)
	result, err = ioutil.ReadFile(webConf)
	assert.NoError(err)
	assert.Contains(string(result), "root /var/www/html/docroot;")
	os.Remove(webConf)

	// test a drupal app definition with all options set
	args[1] = "all_the_things"
	os.Setenv("FILE_SRC", src)
	os.MkdirAll(src, 0755)
	os.Create(src + "/testfile")
	ioutil.WriteFile(webConf, []byte(testConf), os.FileMode(0644))
	_, err = system.RunCommand(bin, args)
	assert.NoError(err)
	// check web config
	webResult, err := ioutil.ReadFile(webConf)
	assert.NoError(err)
	assert.Contains(string(webResult), "root /var/www/html/potato;")
	// check file dir
	fileDest := "potato_pub"
	assert.True(system.FileExists(fileDest))
	assert.True(system.FileExists(fileDest + "/testfile"))
	// check settings.php contents
	result, err = ioutil.ReadFile(confFile)
	assert.NoError(err)
	assert.Contains(string(result), "CONFIG_SYNC_DIRECTORY => '/var/www/html/potato_conf',")
	assert.Contains(string(result), "'database' => \"potato\"")
	assert.Contains(string(result), "'username' => \"spud\"")
	assert.Contains(string(result), "'password' => \"spudtato\"")
	assert.Contains(string(result), "'host' => \"potatodb.com\"")
	assert.Contains(string(result), "'driver' => \"mysql\"")
	assert.Contains(string(result), "'port' => 1234")
	assert.Contains(string(result), "'prefix' => \"spud_\"")
	// cleanup
	os.Remove(webConf)
	os.Remove(confFile)
	os.RemoveAll(src)
	os.RemoveAll(fileDest)
	os.Unsetenv("FILE_SRC")
}
