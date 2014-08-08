package command

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/caiofilipini/got/irc"
)

const (
	ImageSearchUrl = "http://ajax.googleapis.com/ajax/services/search/images"
)

func init() {
	rand.Seed(time.Now().UnixNano())
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

func Image(bot irc.Bot, query string) {
	params := map[string]string{
		"q":    query,
		"v":    "1.0",
		"safe": "active",
		"rsz":  "8",
	}

	if body, err := NewHTTPClient(ImageSearchUrl).With(params).Get(); err == nil {
		var result ImageResults
		json.Unmarshal(body, &result)

		if images := result.Data.Images; len(images) > 0 {
			selected := images[rand.Intn(len(images))]

			bot.Send(selected.UnescapedUrl)
		}
	} else {
		log.Println("ERROR:", err)
	}
}
