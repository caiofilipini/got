package command

import (
	"encoding/json"
	"log"
	"math/rand"
	"regexp"
)

const (
	ImageSearchUrl = "http://ajax.googleapis.com/ajax/services/search/images"
)

type ImageCommand struct {
	pattern *regexp.Regexp
}

func Image() ImageCommand {
	return ImageCommand{regexp.MustCompile(`(?i)(image|img)\s+([^\s].*)`)}
}

func (c ImageCommand) Pattern() *regexp.Regexp {
	return c.pattern
}

func (c ImageCommand) Run(query string) []string {
	return findImages(query, Params{})
}

type GIFCommand struct {
	pattern *regexp.Regexp
}

func GIF() GIFCommand {
	return GIFCommand{regexp.MustCompile(`(?i)(gif|animate)\s+([^\s].*)`)}
}

func (c GIFCommand) Pattern() *regexp.Regexp {
	return c.pattern
}

func (c GIFCommand) Run(query string) []string {
	return findImages(query, Params{"imgtype": "animated"})
}

type Params map[string]string

var defaultParams = Params{
	"v":    "1.0",
	"safe": "active",
	"rsz":  "8",
}

type imageResults struct {
	Data struct {
		Images []struct {
			UnescapedUrl string `json:"unescapedUrl"`
		} `json:"results"`
	} `json:"responseData"`
}

func findImages(query string, params Params) []string {
	for k, v := range defaultParams {
		params[k] = v
	}
	params["q"] = query

	var imgUrl string

	if body, err := NewHTTPClient(ImageSearchUrl).With(params).Get(); err == nil {
		var result imageResults
		json.Unmarshal(body, &result)

		if images := result.Data.Images; len(images) > 0 {
			selected := images[rand.Intn(len(images))]

			imgUrl = selected.UnescapedUrl
		}
	} else {
		log.Println("ERROR:", err)
	}

	var urls []string
	if imgUrl != "" {
		urls = append(urls, imgUrl)
	}
	return urls
}
