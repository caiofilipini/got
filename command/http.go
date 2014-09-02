package command

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Params map[string]string

type HTTPClient struct {
	baseUrl string
	params  Params
}

func NewHTTPClient(url string) HTTPClient {
	return HTTPClient{url, make(Params)}
}

func (c HTTPClient) With(params Params) HTTPClient {
	for k, v := range params {
		c.params[k] = url.QueryEscape(v)
	}
	return c
}

func (c HTTPClient) Get() ([]byte, error) {
	info(fmt.Sprintf("Requesting %s", c.fullUrl()))

	resp, err := http.Get(c.fullUrl())
	if err != nil {
		return nil, err
	}

	info(fmt.Sprintf("Response code: %d", resp.StatusCode))

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c HTTPClient) fullUrl() string {
	if len(c.params) > 0 {
		return fmt.Sprintf("%s?%s", c.baseUrl, c.queryString())
	} else {
		return c.baseUrl
	}
}

func (c HTTPClient) queryString() string {
	var pairs []string
	for k, v := range c.params {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(pairs, "&")
}

func info(msg string) {
	log.Printf("[HTTPClient] %s\n", msg)
}
