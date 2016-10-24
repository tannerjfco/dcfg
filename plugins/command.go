package plugins

import (
	"fmt"
	"strings"
	"time"

	"github.com/drud/drud-go/utils"
)

type Command struct {
	Task
	Cmd string `yaml:"cmd"`
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

		parts := strings.Split(taskPayload, " ")

		err := utils.RunCommandPipe(parts[0], parts[1:])
		if err != nil {
			if !c.Ignore {
				return err
			}
		}

	}
	return nil
}
