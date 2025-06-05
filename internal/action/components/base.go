package components

type Component interface {
	ValidateParams() error
	Run() error
}

type Base struct {
	Parallel bool `hcl:"parallel"`
}

var ComponentMap = map[string]Component{
	"shell":     &ShellComponent{},
	"installer": &InstallerComponent{},
	"cron":      &CronComponent{},
}

func Execute(c Component) error {
	if err := c.ValidateParams(); err != nil {
		return err
	}

	return c.Run()
}
