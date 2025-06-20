package main

import (
	"sensible/internal/kong"
	"sensible/pkg/logger"
)

const version = "v0.0.1"

func main() {
	kctx := kong.ParseCmd(version)
	if kctx == nil {
		return
	}

	if err := kctx.Run(); err != nil {
		logger.Error(err.Error())
	}
}
