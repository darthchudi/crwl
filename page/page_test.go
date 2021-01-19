package page

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"html/template"
	"strings"
	"testing"
)

// LoadMockHTMLPage returns a buffer representing a HTML page
func LoadMockHTMLPage(path string) (*bytes.Buffer, error) {
	tmpl, err := template.ParseFiles(fmt.Sprintf("./mocks/%v", path))

	if err != nil {
		return nil, err
	}

	buff := bytes.NewBuffer([]byte{})

	err = tmpl.Execute(buff, struct{}{})

	if err != nil {
		return nil, err
	}

	return buff, nil
}

func TestNewPage(t *testing.T) {
	tests := []struct {
		htmlPath             string // htmlPath is the path to the mock HTML page
		URL                  string
		ParentURL            string
		expectedInternalURLs int
		expectedAllURLs      int
	}{
		{htmlPath: "page.html", ParentURL: "https://example.com", URL: "https://example.com/privacy", expectedInternalURLs: 3, expectedAllURLs: 5},
	}

	for _, tc := range tests {
		body, err := LoadMockHTMLPage(tc.htmlPath)

		if err != nil {
			t.Fatalf("failed to load mock html: %v", err)
		}

		document, err := goquery.NewDocumentFromReader(body)

		if err != nil {
			t.Fatalf("failed to create document: %v", document)
		}

		page := NewPage(tc.ParentURL, tc.URL, document)

		if len(page.InternalURLs) != tc.expectedInternalURLs {
			t.Fatalf("expected page to have %v internal URLs, found %v", tc.expectedInternalURLs, len(page.InternalURLs))
		}

		if len(page.AllURLs) != tc.expectedAllURLs {
			t.Fatalf("expected page to have %v internal URLs, found %v", tc.expectedAllURLs, len(page.AllURLs))
		}

		// Validate all internal URLS
		for _, url := range page.InternalURLs {
			if !strings.HasPrefix(url, tc.ParentURL) {
				t.Fatalf("expected internal URL to belong to the domain %v, got %v", tc.ParentURL, url)
			}
		}
	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: "/cards", want: "https://example.com/cards"},
		{input: "https://example.com/help/", want: "https://example.com/help"},
	}

	page := Page{ParentURL: "https://example.com", URL: "https://example.com/loans"}

	for _, tc := range tests {
		result := page.normalizeURL(tc.input)

		if result != tc.want {
			t.Fatalf("expected page to normalize url to %v, got %v", tc.want, result)
		}
	}
}
