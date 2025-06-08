package constants

const (
	CronTypeAdd    = "add"
	CronTypeRemove = "remove"
	AddExpr        = "(crontab %s -l 2>/dev/null; echo '%s %s') | crontab %s -"
	RemoveExpr     = "crontab %s -l 2>/dev/null | grep -v '%s' | crontab %s -"
)
