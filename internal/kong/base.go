package kong

import (
	"os"

	"github.com/alecthomas/kong"
)

type CommonFlags struct {
	Env string `help:"Used to refer to the env. Default is base" default:"base" short:"e"`
}

var BaseCommand struct {
	Init   InitCommand   `cmd:"" help:"Initialize sensible by creating a .sensible directory with required files"`
	Run    RunCommand    `cmd:"" help:"Run an action present in .sensible/<env>/actions/ directory"`
	Secret SecretCommand `cmd:"" help:"Manage secrets in the .sensible/<env>/secrets/ directory"`
}

func ParseCmd(version string) *kong.Context {
	if len(os.Args) < 2 {
		_, _ = os.Stdout.WriteString("sensible " + version + "\nRun sensible --help to see available commands.\n")
		return nil
	}

	return kong.Parse(&BaseCommand)
}
