package command

import (
	"fmt"
	"regexp"
	"time"
)

const (
	Beer         = "\U0001f37a"
	Clock        = "\U000023f0"
	StartingHour = 18
)

var queryRegexp *regexp.Regexp

func init() {
	queryRegexp = regexp.MustCompile(`(?i)time|long|til|left|remaining|eta`)
}

func BeerOclock(query string) []string {
	now := time.Now()
	hour := now.Hour()

	var result string

	if queryRegexp.MatchString(query) {
		hourDiff := StartingHour - hour - 1
		minuteDiff := 60 - now.Minute()
		secondDiff := 60 - now.Second()

		result = fmt.Sprintf(
			"%s%s in %d hour(s), %d minute(s) and %d second(s)",
			Beer,
			Clock,
			hourDiff,
			minuteDiff,
			secondDiff)
	} else if hour >= StartingHour {
		result = fmt.Sprintf("YES! %s%s%s", Beer, Beer, Beer)
	} else {
		result = "Not yet. :("
	}
	return []string{result}
}
