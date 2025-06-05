package components

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sensible/constants"

	"github.com/robfig/cron/v3"
)

type CronComponent struct {
	Base
	Name       string `hcl:"cron,label"`
	User       string `hcl:"user,optional"`
	Type       string `hcl:"type"` // "add", "remove"
	Expression string `hcl:"expression,optional"`
	Job        string `hcl:"job"`
}

func (c *CronComponent) ValidateParams() error {
	// type field must be either "add" or "remove"
	if c.Type != constants.CronTypeAdd && c.Type != constants.CronTypeRemove {
		return errors.New("type field must be 'add' or 'remove'")
	}

	// check if the crontab command is available
	if _, err := exec.LookPath("crontab"); err != nil {
		return errors.New("crontab command not found, please install it")
	}

	// additional validation for 'add' type
	if c.Type == constants.CronTypeAdd {
		if c.Expression == "" {
			return errors.New("expression field is required when type is 'add'")
		}
		_, err := cron.ParseStandard(c.Expression)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CronComponent) Run() error {

	// Add -u before the user if specified. Required for cron expression
	userFlag := ""
	if c.User != "" {
		userFlag = "-u " + c.User
	}

	var cmd *exec.Cmd
	switch c.Type {
	case constants.CronTypeAdd:
		cmd = exec.Command("sh", "-c", fmt.Sprintf(constants.AddExpr, userFlag, c.Expression, c.Job, userFlag))
	case constants.CronTypeRemove:
		cmd = exec.Command("sh", "-c", fmt.Sprintf(constants.RemoveExpr, userFlag, c.Job, userFlag))
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to %s cron job: %w", c.Type, err)
	}

	return nil
}
