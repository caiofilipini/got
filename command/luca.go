package command

import "regexp"

type LucaCommand struct {
	pattern *regexp.Regexp
}

func Luca() LucaCommand {
	return LucaCommand{regexp.MustCompile(`(?i)luca\s?(.*)`)}
}

func (c LucaCommand) Pattern() *regexp.Regexp {
	return c.pattern
}

func (c LucaCommand) Run(query string) []string {
	return findImages("grumpy cat", Params{})
}
