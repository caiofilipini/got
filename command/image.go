package command

import (
	"encoding/json"
	"log"
	"math/rand"
)

const (
	ImageSearchUrl = "http://ajax.googleapis.com/ajax/services/search/images"
)

type Params map[string]string

var defaultParams = Params{
	"v":    "1.0",
	"safe": "active",
	"rsz":  "8",
}

type ImageResult struct {
	UnescapedUrl string `json:"unescapedUrl"`
}

type ResponseData struct {
	Images []ImageResult `json:"results"`
}

type ImageResults struct {
	Data ResponseData `json:"responseData"`
}

func Image(query string) []string {
	return findImages(query, Params{})
}

func GIF(query string) []string {
	return findImages(query, Params{"imgtype": "animated"})
}

func findImages(query string, params Params) []string {
	for k, v := range defaultParams {
		params[k] = v
	}
	params["q"] = query

	var imgUrl string

	if body, err := NewHTTPClient(ImageSearchUrl).With(params).Get(); err == nil {
		var result ImageResults
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
