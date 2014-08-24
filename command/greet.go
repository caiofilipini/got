package command

import (
	"fmt"
	"regexp"
)

type GreetCommand struct {
	pattern *regexp.Regexp
}

func Greet() GreetCommand {
	return GreetCommand{regexp.MustCompile(`(?i)greet\s+([^\s].*)`)}
}

func (c GreetCommand) Pattern() *regexp.Regexp {
	return c.pattern
}

func (c GreetCommand) Run(query string) []string {
	return []string{fmt.Sprintf("ohai there, %s!", query)}
}
