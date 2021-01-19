package fetcher

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Fetcher is an abstraction that allows us to
// configure how we fetch pages
type Fetcher interface {
	// Fetch returns the body of URL
	Fetch(url string) ([]byte, error)
}

// HTTPFetcher fetches pages over HTTP using a custom "net/http" client
type HTTPFetcher struct {
	client http.Client
}

// NewHTTPFetcher initializes a new HTTP Fetcher with a custom
// HTTP client with a timeout
func NewHTTPFetcher(timeout time.Duration) *HTTPFetcher {
	client := http.Client{Timeout: timeout}

	return &HTTPFetcher{client: client}
}

// Fetch makes a HTTP request to fetch a URL and returns the URL page body
func (h *HTTPFetcher) Fetch(url string) ([]byte, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	response, err := h.client.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		message := fmt.Sprintf("request failed with http %v", response.StatusCode)
		return nil, errors.New(message)
	}

	pageBody, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	return pageBody, nil
}
