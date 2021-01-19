package crawler

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/darthchudi/crwl/fetcher"
	"github.com/darthchudi/crwl/graph"
	"github.com/darthchudi/crwl/page"
	"github.com/darthchudi/crwl/stats"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type Crawler struct {
	// Starting URL to crawl
	URL string

	// Fetcher fetches pages a URL and returns the page body
	Fetcher fetcher.Fetcher

	// Workers is numbers of workers to be created in the worker queue
	Workers int

	// LogWriter sets the destination to which the crawler's log data will be written
	// Set to `os.Stdout` by default
	LogWriter io.Writer

	// Cache is a cache of URLs that have been fetched, implemented as a graph
	Graph *graph.Graph

	// Stats provides statistical data about crawler operations
	Stats *stats.Stats

	// errors is a channel through which we receive errors on crawler operations
	errors chan error

	// wg is used to sync goroutines
	wg *sync.WaitGroup
}

// NewCrawler initializes a new crawler with a given number of worker instances.
// The number of worker instances determines the number of pages that can be fetched
// at the same time.
// A request timeout specifies the timeout for HTTP requests to fetch pages
func NewCrawler(url string, workers int, timeout time.Duration) *Crawler {
	// Remove trailing slash in the url
	if strings.HasSuffix(url, "/") {
		url = url[:len(url)-1]
	}

	httpFetcher := fetcher.NewHTTPFetcher(timeout)

	return &Crawler{
		URL:       url,
		Fetcher:   httpFetcher,
		LogWriter: os.Stdout,
		Workers:   workers,
		Graph:     graph.NewGraph(),
		Stats:     stats.NewStats(),
		errors:    make(chan error),
		wg:        new(sync.WaitGroup),
	}
}

// listenForURLs creates a new channel for sending urls.
// It starts n Goroutines (according to the concurrency limit)
// that listen for URLs on the URL channel
// When a new url is received on the channel, the HTML page which the link
// points to is fetched and sent to the parser channel.
//
// A send-only channel for sending URLs is returned (URL Channel) and a
// receive-only channel for getting pages to be parsed (parserChannel)
func (c *Crawler) listenForURLs() (chan<- string, <-chan page.RawPage) {
	urlChannel := make(chan string)
	parserChannel := make(chan page.RawPage)

	for i := 1; i <= c.Workers; i++ {
		go func() {
			for url := range urlChannel {
				rawHTMlBody, err := c.Fetcher.Fetch(url)

				if err != nil {
					httpError := fmt.Errorf("failed to fetch %v: %v", url, err)
					c.errors <- httpError
					continue
				}

				rawPage := page.RawPage{
					URL:  url,
					Body: rawHTMlBody,
				}

				// Send the raw page body to the parser channel
				parserChannel <- rawPage
			}
		}()
	}

	return urlChannel, parserChannel
}

// startParser begins goroutines that read raw pages (in bytes),
// parses and evaluates them and returns a receive only channel for
// getting page results
func (c *Crawler) startParser(parserChannel <-chan page.RawPage) chan page.Page {
	pageChannel := make(chan page.Page)

	go func() {
		for rawPage := range parserChannel {
			go func(rawPage page.RawPage) {
				document, err := goquery.NewDocumentFromReader(bytes.NewReader(rawPage.Body))

				if err != nil {
					c.errors <- err
					return
				}

				newPage := page.NewPage(c.URL, rawPage.URL, document)

				// Send processed page to the page channel
				pageChannel <- newPage
			}(rawPage)
		}
	}()

	return pageChannel
}

// listenForPages gets fetched pages and queues urls that haven't been visited in the page to be fetched
func (c *Crawler) listenForPages(pageChannel <-chan page.Page, urlChannel chan<- string) {
	// Start a single goroutine that acts as the coordinator by processing pages (results)
	// and dispatching new URLs to be fetched from the pages.
	go func() {
		for p := range pageChannel {
			for _, url := range p.InternalURLs {

				visited := c.Graph.HasNode(url)

				if visited {
					c.Graph.AddEdge(p.URL, url)
					continue
				}

				c.Graph.AddNode(url)
				c.Graph.AddEdge(p.URL, url)
				c.wg.Add(1)
				c.Stats.RecordNewOperation()
				go func(u string) { urlChannel <- u }(url) // Send url to workers in a new goroutine to prevent blocking if all workers are busy
			}

			p.Print(c.LogWriter)
			c.Stats.RecordOperationCompletion()
			c.wg.Done()
		}
	}()
}

// listenForErrors listens for errors and decrements the wait group when an error occurs
func (c *Crawler) listenForErrors() {
	go func() {
		for err := range c.errors {
			fmt.Fprintf(c.LogWriter, "🥞 %v\n", err)

			c.Stats.RecordOperationFailure()
			c.wg.Done()
		}
	}()
}

// Crawl begins crawling the crawler's URL
func (c *Crawler) Crawl() {
	urlChannel, parserChannel := c.listenForURLs()

	pageChannel := c.startParser(parserChannel)

	c.listenForPages(pageChannel, urlChannel)

	c.listenForErrors()

	// Start crawling by sending the crawler URL to the URL channel
	c.Graph.AddNode(c.URL)
	c.wg.Add(1)
	c.Stats.RecordNewOperation()
	urlChannel <- c.URL

	c.Stats.RecordStartTime()
	c.wg.Wait()
	c.Stats.RecordTotalDuration()
}
