package command

import "regexp"

type LucaCommand struct {
	name    string
	pattern *regexp.Regexp
}

func Luca() LucaCommand {
	return LucaCommand{
		"luca",
		regexp.MustCompile(`(?i)(luca)\s?(.*)`),
	}
}

func (c LucaCommand) Name() string {
	return c.name
}

func (c LucaCommand) Pattern() *regexp.Regexp {
	return c.pattern
}

func (c LucaCommand) Help() string {
	return c.name + " – a tribute to Luca Pette"
}

func (c LucaCommand) Usage() []string {
	return []string{c.name}
}

func (c LucaCommand) Run(query string) []string {
	return findImages("grumpy cat", Params{})
}
