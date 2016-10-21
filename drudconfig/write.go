package drudconfig

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/drud/drud-go/utils"
)

type Write struct {
	Task
	Write string      `yaml:"write"`
	Mode  os.FileMode `yaml:"mode",json:"write"`
}

func (w Write) GetPayload() string {
	return "fudge"
}

func (w Write) Pretty() {
	fmt.Println(utils.Prettify(w))
}

func (w *Write) Run() error {

	if w.Wait != "" {
		lengthOfWait, _ := time.ParseDuration(w.Wait)
		time.Sleep(lengthOfWait)
	}

	if w.Write == "" {
		return fmt.Errorf("Nothing to write for task: %s", w.Name)
	}

	err := ioutil.WriteFile(w.Dest, []byte(w.Write), w.Mode)
	if err != nil {
		log.Fatalln("Could not read config file:", err)
	}

	info, _ := os.Stat(w.Dest)
	if info.Mode() != w.Mode {
		err := os.Chmod(w.Dest, w.Mode)
		if err != nil {
			return err
		}
	}
	return nil
}
