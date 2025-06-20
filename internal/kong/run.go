package kong

import (
	"sensible/internal/action"
)

type RunCommand struct {
	File string `help:"Mandatory file path for the action which you want to run" required:"" type:"path" short:"f"`
}

func (c *RunCommand) Run() error {
	variables, groups, err := action.Sync()
	if err != nil {
		return err
	}

	return action.Do(c.File, variables, groups)
}
