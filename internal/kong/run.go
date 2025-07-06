package kong

import (
	"fmt"
	"path/filepath"
	"sensible/internal/action"
	"sensible/internal/constants"
)

type RunCommand struct {
	CommonFlags

	File string `help:"Mandatory action file name inside .sensible/action/ which you want to run" required:"" short:"f"`
}

func (c *RunCommand) Run() error {
	if filepath.Ext(c.File) != ".hcl" {
		c.File += ".hcl"
	}

	actionsDir := fmt.Sprintf(constants.ActionsDir, c.Env)
	actionFile, err := filepath.Abs(filepath.Join(actionsDir, c.File))
	if err != nil {
		return err
	}

	// we will need the variables
	values, err := action.GetVariables(actionFile, c.Env)
	if err != nil {
		return err
	}

	return action.Do(actionFile, values, c.Env)
}
