package command

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
)

const (
	WeatherSeachUrl = "http://api.openweathermap.org/data/2.5/weather"
)

type WeatherCommand struct {
	name    string
	pattern *regexp.Regexp
}

func Weather() WeatherCommand {
	return WeatherCommand{
		"weather",
		regexp.MustCompile(`(?i)weather\s+([^\s].*)`),
	}
}

func (c WeatherCommand) Name() string {
	return c.name
}

func (c WeatherCommand) Pattern() *regexp.Regexp {
	return c.pattern
}

func (c WeatherCommand) Help() string {
	return c.name + " – shows weather conditions"
}

func (c WeatherCommand) Usage() []string {
	return []string{
		c.name + " <city> – shows current weather conditions for the given city",
	}
}

func (c WeatherCommand) Run(query string) []string {
	if body, err := NewHTTPClient(WeatherSeachUrl).With(Params{"q": query}).Get(); err == nil {
		var result weatherResults
		json.Unmarshal(body, &result)

		celsius := int(math.Floor(result.Main.Kelvin-273.15) + 0.5)

		return []string{
			fmt.Sprintf("%s, %s: %d°C, %s",
				result.Name,
				result.Sys.Country,
				celsius,
				result.Weather[0].Description),
		}
	}
	return []string{}
}

type weatherResults struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
	Sys struct {
		Country string `json:"country"`
	} `json:"sys"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}
