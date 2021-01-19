// Provides a mock fetcher that is only available within crawler tests
package crawler

import (
	"bytes"
	"fmt"
	"html/template"
)

// fetcherCache is a map of mock urls and the path to their equivalent
// mock html file
var mockFetcherCache = map[string]string{
	"https://example.com":             "index.html",
	"https://example.com/loans":       "loans.html",
	"https://example.com/shared-tabs": "shared-tabs.html",
}

// MockFetcher is a mock fetcher that fetches URLs from an
// internal mock cache. This allows us to inject canned
// response bodies within tests
type MockFetcher struct{}

// Fetch fetches a URL from the internal fetcher cache.
// If the url is found in the cache, it returns the body of
// the mock html file the URL points to
func (f MockFetcher) Fetch(url string) ([]byte, error) {
	if path, exists := mockFetcherCache[url]; exists {
		return f.getHTMLPage(path)
	}

	return nil, fmt.Errorf("%v not found in mock cache", url)
}

// getHTMLPage returns a mock HTML page
func (f MockFetcher) getHTMLPage(path string) ([]byte, error) {
	tmpl, err := template.ParseFiles(fmt.Sprintf("./mocks/%v", path))

	if err != nil {
		return nil, err
	}

	buff := bytes.NewBuffer([]byte{})

	err = tmpl.Execute(buff, struct{}{})

	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
