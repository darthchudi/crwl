# crwl ✨🕷

Crwl is a concurrent web crawler written in Go.


## Usage

Crwl allows you to specify:
- The URL to be fetched via the `--url` flag
 - The number of URLs that can be fetched in parallel via the `--workers` flag (default: 20)
 - The request timeout duration for fetching each URL via the `--timeout` flag (default: 30 seconds)

````
go run main.go --url=https://example.com --workers=10 --timeout=30s
````

## Testing

````
go test -v ./...
````

## How it works

![](https://res.cloudinary.com/chudi/image/upload/v1611075709/A4_-_1.png)

Crwl uses a worker queue implemented as a set of worker goroutines listening on a channel, where URLs to be fetched are sent. When a worker goroutine successfully fetches a URL page, it sends it to a separate set of goroutines where parsing and link extraction is done. This allows us to utilize our worker goroutines exclusively for HTTP requests i.e we fetch pages as quickly as possible and dispense them to the next stage in the crawling pipeline. 

After a URL page has been parsed and links have been extracted, it is sent to the event loop/coordinator goroutine via a `Page Channel`. When we receive a parsed page in the event loop goroutine, we iteratively send all internal URLs (i.e URLs within the crawler's URL subdomain) on the page that haven't been visited to the worker queue to be fetched.

When the worker queue is empty and all pending tasks have been completed, the Crawler stats are printed and it exits.
