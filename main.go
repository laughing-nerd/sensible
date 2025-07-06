package main

import (
	"sensible/internal/kong"
	"sensible/pkg/logger"
)

var version = "v0.0.0"

func main() {
	kctx := kong.ParseCmd(version)
	if kctx == nil {
		return
	}

	if err := kctx.Run(); err != nil {
		logger.Error(err.Error())
	}
}
