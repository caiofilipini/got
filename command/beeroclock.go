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

type BeerOClockCommand struct {
	name    string
	pattern *regexp.Regexp
}

func BeerOClock() BeerOClockCommand {
	return BeerOClockCommand{
		"beer",
		regexp.MustCompile(`(?i)beer\s*(.*)`),
	}
}

func (c BeerOClockCommand) Name() string {
	return c.name
}

func (c BeerOClockCommand) Pattern() *regexp.Regexp {
	return c.pattern
}

func (c BeerOClockCommand) Help() string {
	return c.name + " – is it beer o'clock yet?"
}

func (c BeerOClockCommand) Usage() []string {
	return []string{
		c.name + " – tells if it's beer o'clock",
		c.name + " ETA – tells how long until beer o'clock",
	}
}

func (c BeerOClockCommand) Run(query string) []string {
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
