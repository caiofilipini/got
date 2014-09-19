package command

import (
	"fmt"
	"math/rand"
	"regexp"
)

var swearings = []string{
	"a maronn",
	"san giuseppe",
	"san pietro",
	"o patatern 'n croc",
	"tutti i santi",
	"gesu",
	"gesu bambin 'n croc",
	"gesu crist",
}

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
	return []string{fmt.Sprintf("annagg %s", swearings[rand.Intn(len(swearings))])}
}
