package command

import (
	"fmt"
	"time"
)

const (
	Beer = "\U0001f37a"
)

func BeerOclock(query string) []string {
	hour := time.Now().Hour()

	var result string
	if hour >= 18 {
		result = fmt.Sprintf("YES! %s", Beer)
	} else {
		result = "Not yet. :("
	}
	return []string{result}
}
