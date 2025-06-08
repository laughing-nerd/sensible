package kong

import (
	"os"
	"sensible/internal/initialize"
)

type InitCommand struct {
	File string `help:"Optional file path." type:"path" short:"f" default:"./.sensible"`
}

func (c *InitCommand) Run() error {
	os.Stdout.WriteString("Initializing sensible\n")
	return initialize.Start(c.File)
}
