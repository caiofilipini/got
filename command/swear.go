package command

import "regexp"

type Swear struct{}

func (s Swear) Pattern() *regexp.Regexp {
	return regexp.MustCompile(`(?i)swear\s?(.*)`)
}

func (s Swear) Run(query string) []string {
	return []string{"annagg a maronn"}
}
