package page

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	netUrl "net/url"
	"strings"
)

type RawPage struct {
	// URL is the page url
	URL string

	// Raw HTML of the page
	Body []byte
}

type Page struct {
	// ParentURL is the URL of the web crawler instance
	// that fetched the page
	ParentURL string

	// URL is the page url
	URL string

	// Document is a goquery representation of the page HTML document
	Document *goquery.Document

	// AllURLs are all the links found in the page
	AllURLs []string

	// InternalURLs are links found in the page that belong to the
	// same domain as the web crawler url
	InternalURLs []string
}

// NewPage creates a new page and populates it's links from its HTML
// document
func NewPage(parentURL, URL string, document *goquery.Document) Page {
	page := Page{
		ParentURL: parentURL, URL: URL, Document: document, AllURLs: []string{}, InternalURLs: []string{},
	}

	page.fetchLinks()

	return page
}

// normalizeURL removes leading and trailing slashes in a URL
func (p *Page) normalizeURL(url string) string {
	// if the url is a relative url, normalize it
	if strings.HasPrefix(url, "/") {
		url = p.ParentURL + url
	}

	// if the url has a trailing slash, normalize it
	if strings.HasSuffix(url, "/") {
		url = url[:len(url)-1]
	}

	return url
}

// isURLSameDomain checks if a url belongs to the same domain
// as the page's parent
func (p *Page) isURLSameDomainAsParent(url string) (bool, error) {
	parsedParentURL, err := netUrl.Parse(p.ParentURL)

	if err != nil {
		return false, err
	}

	parsedURL, err := netUrl.Parse(url)

	if err != nil {
		return false, err
	}

	if parsedParentURL.Host != parsedURL.Host {
		return false, nil
	}

	return true, nil
}

// fetchLinks gets all URLs in the page and finds internal (local) URLs
func (p *Page) fetchLinks() {
	allURLs := []string{}
	internalURLs := []string{}

	// Used to deduplicate stored URLs
	allURLsCache := NewSet()
	internalURLsCache := NewSet()

	p.Document.Find("a").Each(func(i int, s *goquery.Selection) {
		url, exists := s.Attr("href")

		if !exists {
			return
		}

		url = p.normalizeURL(url)

		// Only add this URL to the all URLs array if we haven't seen it before
		if !allURLsCache.Has(url) {
			allURLs = append(allURLs, url)
			allURLsCache.Add(url)
		}

		isInternalURL, err := p.isURLSameDomainAsParent(url)

		if err != nil {
			log.Printf("🤠 %v Domain validation error \n", err)
			return
		}

		if !isInternalURL {
			return
		}

		// Only add this URL to the all URLs array only if we haven't added it before
		if !internalURLsCache.Has(url) {
			internalURLs = append(internalURLs, url)
			internalURLsCache.Add(url)
		}

		return
	})

	p.AllURLs = allURLs
	p.InternalURLs = internalURLs
}

// Print prints out all links in a page to a writer
func (p *Page) Print(w io.Writer) {
	message := fmt.Sprintf("✨ Extracted links in URL: %v \n", p.URL)

	for _, url := range p.AllURLs {
		message += fmt.Sprintf("\t ✨ %v \n", url)
	}

	message += "\n"

	fmt.Fprint(w, message)
}
