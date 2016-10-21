package drudconfig

import (
	"fmt"
	"time"

	"github.com/drud/drud-go/utils"
)

type Command struct {
	Task
	Cmd string `yaml:"cmd"`
}

func (c Command) GetPayload() string {
	return "fudge"
}

func (c Command) Pretty() {
	fmt.Println(utils.Prettify(c))
}

func (c *Command) Run() error {

	for i := c.Repeat; i >= 0; i-- {

		if c.Wait != "" {
			lengthOfWait, _ := time.ParseDuration(c.Wait)
			time.Sleep(lengthOfWait)
		}

		taskPayload := c.Cmd
		if taskPayload == "" {
			return fmt.Errorf("No cmd specified")
		}

		err := RunCommand(taskPayload)
		if err != nil {
			if !c.Ignore {
				return err
			}
		}

	}
	return nil
}
