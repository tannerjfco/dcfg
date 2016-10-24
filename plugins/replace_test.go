package plugins

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testDir         = "testing"
	testFile        = testDir + "/testfile,txt"
	testContent     = "testkey: testvalue"
	testReplace     = "heyheyhey"
	testRegex       = `(testkey:) ([a-z]*)`
	regexReplace    = "$1 newvalue"
	testRegexEndVal = "testkey: newvalue"
)

func TestReplaceString(t *testing.T) {
	setup()
	defer teardown()

	assert := assert.New(t)

	task := Replace{
		Find:    testContent,
		Replace: "heyheyhey",
	}
	task.Name = "test replace"
	task.Dest = testFile

	err := task.Run()
	assert.Nil(err)

	content, err := ioutil.ReadFile(testFile)
	assert.Nil(err)
	assert.Contains(string(content), testReplace)

}

func TestReplaceStringRegex(t *testing.T) {
	setup()
	defer teardown()
	assert := assert.New(t)

	task := Replace{
		Find:    testRegex,
		Replace: regexReplace,
	}
	task.Name = "test replace"
	task.Dest = testFile

	err := task.Run()
	assert.Nil(err)

	content, err := ioutil.ReadFile(testFile)
	assert.Nil(err)
	assert.Contains(string(content), testRegexEndVal)

}

func teardown() {
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		return
	}

	err := os.RemoveAll(testDir)
	if err != nil {
		log.Fatal(err)
	}
}

func setup() {
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		os.Mkdir(testDir, 0777)
	}
	err := ioutil.WriteFile(testFile, []byte(testContent), 0777)
	if err != nil {
		log.Fatal(err)
	}
}
