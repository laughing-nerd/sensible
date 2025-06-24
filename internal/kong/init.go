package kong

import (
	"sensible/internal/initialize"
)

type InitCommand struct {
	CommonFlags
	// File string `help:"Optional file path." type:"path" short:"f" default:"./.sensible"`
}

func (c *InitCommand) Run() error {
	return initialize.Start(c.Env)
}
