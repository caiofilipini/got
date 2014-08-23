package command

import (
	"fmt"
	"regexp"
)

type Greet struct{}

func (g Greet) Pattern() *regexp.Regexp {
	return regexp.MustCompile(`(?i)greet\s+([^\s].*)`)
}

func (g Greet) Run(query string) []string {
	return []string{fmt.Sprintf("ohai there, %s!", query)}
}
