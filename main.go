package main

import (
	"sensible/internal/kong"
)

const version = "v0.0.1"

func main() {

	kctx := kong.ParseCmd(version)
	if kctx == nil {
		return
	}

	kctx.Run()

	// initialize.FetchMetadata()

	// action.Parse("./internal/initialize/templates/action.hcl")
	// return
}
