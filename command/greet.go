package command

import "fmt"

func Greet(query string) []string {
	return []string{fmt.Sprintf("ohai there, %s!", query)}
}
