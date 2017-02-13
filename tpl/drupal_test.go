package tpl

import (
	"testing"

	"os"

	"io/ioutil"

	"github.com/drud/drud-go/utils/system"
	"github.com/stretchr/testify/assert"
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

	// // test configuring a site contained in a folder called docroot
	// args[1] = "have_docroot"

	// // test a drupal app definition with all options set
	// args[1] = "all_the_things"

}
