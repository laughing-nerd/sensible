package kong

import "os"

type RunCommand struct {
	File string `help:"Mandatory file path for the action which you want to run" required:"" type:"path" short:"f"`
}

func (c *RunCommand) Run() error {
	os.Stdout.WriteString("Initializing sensible\n")
	return nil
}
