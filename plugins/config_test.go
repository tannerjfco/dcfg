package plugins

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigCreate(t *testing.T) {

	assert := assert.New(t)

	task := Config{
		Delim: " = ",
		Items: map[string]string{
			"name":  "configuration",
			"count": "10",
			"debug": "true",
		},
	}
	task.Name = "test config"
	task.Dest = "newconfig.conf"

	err := task.Run()
	assert.Nil(err)

	content, err := ioutil.ReadFile("newconfig.conf")
	assert.Nil(err)
	assert.Contains(string(content), "name = configuration")

	err = os.Remove(task.Dest)
	assert.Nil(err)

}

func TestConfigReplace(t *testing.T) {

	assert := assert.New(t)

	initialCfg := "count = v\n"
	err := ioutil.WriteFile("newconfig2.conf", []byte(initialCfg), os.FileMode(0774))
	assert.Nil(err)

	task := Config{
		Delim: " = ",
		Items: map[string]string{
			"name":  "configuration",
			"count": "10",
			"debug": "true",
		},
	}
	task.Dest = "newconfig2.conf"

	err = task.Run()
	assert.Nil(err)

	content, err := ioutil.ReadFile("newconfig2.conf")
	assert.Nil(err)
	assert.Contains(string(content), "name = configuration")
	assert.Contains(string(content), "count = 10")
	assert.Contains(string(content), "debug = true")

	err = os.Remove(task.Dest)
	assert.Nil(err)

}
