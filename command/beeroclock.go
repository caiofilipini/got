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

type BeerOclock struct{}

func (b BeerOclock) Pattern() *regexp.Regexp {
	return regexp.MustCompile(`(?i)beer\s*(.*)`)
}

func (b BeerOclock) Run(query string) []string {
	fmt.Println(query)
	now := time.Now()
	hour := now.Hour()
	beerOclock := hour >= StartingHour
	beerOclockEmoji := fmt.Sprintf("%s%s", Beer, Clock)

	var result string

	if queryRegexp.MatchString(query) {
		hourDiff := StartingHour - hour - 1
		minuteDiff := 60 - now.Minute()
		secondDiff := 60 - now.Second()

		if beerOclock {
			result = fmt.Sprintf("It's already %s! Enjoy!", beerOclockEmoji)
		} else {
			result = fmt.Sprintf(
				"%s in %d hour(s), %d minute(s) and %d second(s)",
				beerOclockEmoji,
				hourDiff,
				minuteDiff,
				secondDiff)
		}
	} else if beerOclock {
		result = fmt.Sprintf("YES! %s%s%s", Beer, Beer, Beer)
	} else {
		result = "Not yet. :("
	}
	return []string{result}
}
