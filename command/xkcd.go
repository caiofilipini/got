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
	Num        int    `json:"num"`
	Day        string `json:"day"`
	Month      string `json:"month"`
	Year       string `json:"year"`
	Img        string `json:"img"`
	Link       string `json:"link"`
	SafeTitle  string `json:"safe_title"`
	Title      string `json:"title"`
	Alt        string `json:"alt"`
	Transcript string `json:"transcript"`
	News       string `json:"news"`
}

func XKCD(query string) []string {
	q := strings.Trim(query, " ")
	current := loadComic(XKCDLatestUrl)
	comic := current

	if q == "random" {
		comic = loadComic(fmt.Sprintf(XKCDNumberUrl, rand.Intn(current.Num)))
	} else if match := numRegexp.FindStringSubmatch(q); len(match) > 1 {
		if number, err := strconv.Atoi(match[1]); err == nil {
			if number == 404 {
				return []string{"smart ass :)"}
			}

			comic = loadComic(fmt.Sprintf(XKCDNumberUrl, number))
		}
	}

	return []string{
		comic.Img,
		comic.Title,
		comic.Alt,
	}
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
