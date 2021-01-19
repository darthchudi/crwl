package crawler

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestCrawl(t *testing.T) {
	crawler := NewCrawler("https://example.com", 10, time.Second*20)
	crawler.Fetcher = MockFetcher{}
	crawler.Crawl()

	// We expect that the crawler visits only known links
	expectedCompleted := int64(len(mockFetcherCache))
	if crawler.Stats.Completed() != expectedCompleted {
		t.Errorf("expected crawler to have completed %v tasks, got %v", expectedCompleted, crawler.Stats.Completed())
	}

	// We expect that there will be no pending tasks by the time Crawl() returns
	if crawler.Stats.Pending() != 0 {
		t.Errorf("expected crawler to have 0 pending tasks, got %v", crawler.Stats.Pending())
	}

	// Check that all known links are in the crawler's cache
	for url := range mockFetcherCache {
		if !crawler.Graph.HasNode(url) {
			t.Errorf("expected crawler to have visited url %v", url)
		}
	}
}

func TestCrawlError(t *testing.T) {

	// Provide a URL that is not recognized by the mock fetcher
	// This will cause fetch to fail and the crawler should record an error
	crawler := NewCrawler("https://test.com", 10, time.Second*20)
	crawler.Fetcher = MockFetcher{}

	// Use a mock log writer
	logWriter := bytes.NewBuffer([]byte{})
	crawler.LogWriter = logWriter

	crawler.Crawl()

	expectedFailures := 1
	if crawler.Stats.Failures() != int64(expectedFailures) {
		t.Errorf("expected crawler to have failed %v task(s), got %v", expectedFailures, crawler.Stats.Failures())
	}

	// We expect that there will be no pending tasks by the time Crawl() returns
	if crawler.Stats.Pending() != 0 {
		t.Errorf("expected crawler to have 0 pending tasks, got %v", crawler.Stats.Pending())
	}

	// Check that we can capture crawler errors
	expected := "failed to fetch https://test.com"
	errorMessage := logWriter.String()

	if !strings.Contains(errorMessage, expected) {
		t.Fatalf("expected error message to contain %v, got %v", expected, errorMessage)
	}
}
