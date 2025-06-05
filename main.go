package main

import (
	"fmt"
	"os"
	"sensible/internal"
	"sensible/internal/action"
	"sensible/internal/initialize"

	"github.com/alecthomas/kingpin"
)

const version = "v0.0.1"

var (
	verbose = kingpin.Flag("verbose", "Verbose mode.").Short('v').String()
	name    = kingpin.Arg("action", "In", ).String()
)

func main() {

	kingpin.Parse()
	fmt.Printf("%v, %s\n", *verbose, *name)
	return

	// initialize.FetchMetadata()

	action.Parse("./internal/initialize/templates/action.hcl")
	return

	if os.Args == nil || len(os.Args) < 2 {
		os.Stdout.WriteString("sensible " + version + "\nRun sensible --help to see available commands.\n")
		return
	}

	command := os.Args[1]
	switch command {
	case "init":
		initialize.Start()
	case "sync":
		internal.Sync()
	}
}
