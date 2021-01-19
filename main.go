package main

import (
	"flag"
	"fmt"
	"github.com/darthchudi/crwl/crawler"
	"time"
)

func main() {
	crawlURL := flag.String("url", "https://example.com", "URL to Crawl")
	workers := flag.Int("workers", 20, "Workers defines the maximum number of concurrent connections to the provided domain")
	requestTimeout := flag.Duration("timeout", 30*time.Second, "How long should a request to fetch a page take")

	flag.Parse()

	c := crawler.NewCrawler(*crawlURL, *workers, *requestTimeout)
	c.Crawl()

	c.Stats.Print()
	fmt.Printf("Finished crawling %v URLs in in %v", c.Stats.Total(), c.Stats.Duration())
}
