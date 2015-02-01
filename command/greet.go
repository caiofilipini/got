package command

import (
	"fmt"
	"regexp"
)

type GreetCommand struct {
	name    string
	pattern *regexp.Regexp
}

func Greet() GreetCommand {
	return GreetCommand{
		"greet",
		regexp.MustCompile(`(?i)greet\s+([^\s].*)`),
	}
}

func (c GreetCommand) Name() string {
	return c.name
}

func (c GreetCommand) Pattern() *regexp.Regexp {
	return c.pattern
}

func (c GreetCommand) Help() string {
	return c.name + " – shows greetings"
}

func (c GreetCommand) Usage() []string {
	return []string{
		c.name + " <nickname>",
	}
}

func (c GreetCommand) Run(query string) []string {
	return []string{fmt.Sprintf("ohai there, %s!", query)}
}
