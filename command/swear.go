package command

import "regexp"

type SwearCommand struct {
	pattern *regexp.Regexp
}

func Swear() SwearCommand {
	return SwearCommand{regexp.MustCompile(`(?i)swear\s?(.*)`)}
}

func (c SwearCommand) Pattern() *regexp.Regexp {
	return c.pattern
}

func (c SwearCommand) Run(query string) []string {
	return []string{"annagg a maronn"}
}
