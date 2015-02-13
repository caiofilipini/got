package command

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Params is an alias for map[string]string.
// Its intended to facilitate working with request
// parameters.
type Params map[string]string

// HTTPClient is the basic abstraction for performing
// HTTP requests.
type HTTPClient struct {
	// The base URL to be requested.
	baseUrl string

	// The request parameters.
	params Params
}

// NewHTTPClient returns a new client for the given URL.
func NewHTTPClient(url string) HTTPClient {
	return HTTPClient{url, make(Params)}
}

// With configures the request parameters. It escapes the
// parameter values.
func (c HTTPClient) With(params Params) HTTPClient {
	for k, v := range params {
		c.params[k] = url.QueryEscape(v)
	}
	return c
}

// Get performs a GET request with the pre-configured
// base URL and parameters. It returns the bytes received
// as response, or an error.
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

// fullUrl returns the URL including the query string,
// generated from the configured parameters, if any.
func (c HTTPClient) fullUrl() string {
	if len(c.params) > 0 {
		return fmt.Sprintf("%s?%s", c.baseUrl, c.queryString())
	} else {
		return c.baseUrl
	}
}

// queryString generates the query string.
func (c HTTPClient) queryString() string {
	var pairs []string
	for k, v := range c.params {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(pairs, "&")
}

// info writes the message into the log.
func info(msg string) {
	log.Printf("[HTTPClient] %s\n", msg)
}
