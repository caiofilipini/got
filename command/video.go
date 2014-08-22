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

type VideoLink struct {
	Rel  string `json:"rel"`
	Type string `json:"type"`
	Href string `json:"href"`
}

type VideoResult struct {
	Links []VideoLink `json:"link"`
}

type FeedData struct {
	Entries []VideoResult `json:"entry"`
}

type VideoResults struct {
	Feed FeedData `json:"feed"`
}

type Video struct{}

func (v Video) Pattern() *regexp.Regexp {
	return regexp.MustCompile(`(?i)video|youtube|yt\s+([^\s]+)`)
}

func (v Video) Run(query string) []string {
	params := map[string]string{
		"q":           query,
		"orderBy":     "relevance",
		"max-results": "15",
		"alt":         "json",
	}

	if body, err := NewHTTPClient(VideoSearchUrl).With(params).Get(); err == nil {
		var result VideoResults
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
