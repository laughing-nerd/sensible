package kong

import (
	"path/filepath"
	"sensible/internal/action"
	"sensible/internal/constants"
	"strings"
)

type RunCommand struct {
	File string `help:"Mandatory action file name inside .sensible/action/ which you want to run" required:"" short:"f"`
}

func (c *RunCommand) Run() error {
	index := strings.Index(c.File, ".hcl")
	if index == -1 {
		c.File += ".hcl"
	}

	f, err := filepath.Abs(filepath.Join(constants.ActionsDir, c.File))
	if err != nil {
		return err
	}

	variables, groups, err := action.Sync()
	if err != nil {
		return err
	}

	return action.Do(f, variables, groups)
}
