package kong

import (
	"os"

	"github.com/alecthomas/kong"
)

var BaseCommand struct {
	Init InitCommand `cmd:"" help:"Initialize sensible by creating a .sensible directory with required files"`
	Run  RunCommand  `cmd:"" help:"Run an action present in .sensible/actions/ directory"`
}

func ParseCmd(version string) *kong.Context {
	if len(os.Args) < 2 {
		os.Stdout.WriteString("sensible " + version + "\nRun sensible --help to see available commands.\n")
		return nil
	}

	return kong.Parse(&BaseCommand)
}
