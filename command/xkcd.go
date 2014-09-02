package command

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

const (
	XKCDLatestUrl = "http://xkcd.com/info.0.json"
	XKCDNumberUrl = "http://xkcd.com/%d/info.0.json"
)

var numRegexp = regexp.MustCompile(`(^\d+$)`)

type XKCDResult struct {
	Num   int    `json:"num"`
	Img   string `json:"img"`
	Title string `json:"title"`
	Alt   string `json:"alt"`
}

type XKCDCommand struct {
	pattern *regexp.Regexp
}

func XKCD() XKCDCommand {
	return XKCDCommand{regexp.MustCompile(`(?i)xkcd\s*(.*)`)}
}

func (c XKCDCommand) Pattern() *regexp.Regexp {
	return c.pattern
}

func (c XKCDCommand) Run(query string) []string {
	q := strings.Trim(query, " ")

	if match := numRegexp.FindStringSubmatch(q); len(match) > 1 {
		if number, err := strconv.Atoi(match[1]); err == nil {
			if number == 404 {
				return []string{"smart ass :)"}
			}

			return render(loadComic(fmt.Sprintf(XKCDNumberUrl, number)))
		}
	}

	current := loadComic(XKCDLatestUrl)
	comic := current

	if q == "random" {
		comic = loadComic(fmt.Sprintf(XKCDNumberUrl, rand.Intn(current.Num)))
	}

	return render(comic)
}

func loadComic(url string) *XKCDResult {
	if body, err := NewHTTPClient(url).Get(); err == nil {
		var result XKCDResult
		json.Unmarshal(body, &result)

		return &result
	} else {
		log.Println("ERROR:", err)
	}
	return nil
}

func render(comic *XKCDResult) []string {
	return []string{
		comic.Img,
		comic.Title,
		comic.Alt,
	}
}
