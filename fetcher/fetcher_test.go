package fetcher

import (
	"testing"
	"time"
)

func TestHTTPFetcher(t *testing.T) {
	fetcher := NewHTTPFetcher(time.Second * 10)

	_, err := fetcher.Fetch("https://google.com")

	if err != nil {
		t.Fatalf("http fetcher error: %v", err)
	}
}
