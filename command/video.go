package command

import (
	"encoding/json"
	"log"
	"math/rand"
	"regexp"
)

const (
	VideoSearchUrl = "http://gdata.youtube.com/feeds/api/videos"
)

type videoResults struct {
	Feed struct {
		Entries []struct {
			Links []struct {
				Rel  string `json:"rel"`
				Type string `json:"type"`
				Href string `json:"href"`
			} `json:"link"`
		} `json:"entry"`
	} `json:"feed"`
}

type VideoCommand struct {
	name    string
	pattern *regexp.Regexp
}

func Video() VideoCommand {
	return VideoCommand{
		"video",
		regexp.MustCompile(`(?i)(video|youtube|yt)\s+([^\s].*)`),
	}
}

func (c VideoCommand) Name() string {
	return c.name
}

func (c VideoCommand) Pattern() *regexp.Regexp {
	return c.pattern
}

func (c VideoCommand) Help() string {
	return c.name + " – video search"
}

func (c VideoCommand) Usage() []string {
	return []string{
		"video|youtube|yt <query>",
	}
}

func (c VideoCommand) Run(query string) []string {
	params := Params{
		"q":           query,
		"orderBy":     "relevance",
		"max-results": "15",
		"alt":         "json",
	}

	if body, err := NewHTTPClient(VideoSearchUrl).With(params).Get(); err == nil {
		var result videoResults
		json.Unmarshal(body, &result)

		if videos := result.Feed.Entries; len(videos) > 0 {
			selected := videos[rand.Intn(len(videos))]
			var link string

			for _, l := range selected.Links {
				if l.Type == "text/html" && l.Rel == "alternate" {
					link = l.Href
				}
			}

			if link != "" {
				return []string{link}
			}
		}
	} else {
		log.Println("ERROR:", err)
	}

	return []string{}
}
