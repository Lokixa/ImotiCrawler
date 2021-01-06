package imoti

import (
	"io"
	"net/http"
	"regexp"

	"golang.org/x/net/html"
)

// Crawler is the base crawler struct
type Crawler struct {
	Links     chan string
	seenLinks map[string]bool
}

// NewCrawler is a crawler constructor
func NewCrawler() Crawler {
	return Crawler{make(chan string), make(map[string]bool)}
}

// FetchLinks gets url content and extracts links
func (crawler Crawler) FetchLinks(url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	return crawler.extractLinks(res.Body, crawler.Links)
}

// ExtractLinks goes through the html
// and sends out anchors with an href matching the regex
func (crawler Crawler) extractLinks(body io.ReadCloser, out chan<- string) error {
	linkRegex := regexp.MustCompile(`.*\/bg\/obiava/.*`)
	z := html.NewTokenizer(body)
	for {
		token := z.Next()
		switch token {
		case html.StartTagToken, html.EndTagToken:
			tn, _ := z.TagName()
			// if anchor
			if len(tn) == 1 && tn[0] == 'a' {
				for {
					key, val, more := z.TagAttr()
					// has link
					if string(key) == "href" {
						valStr := string(val)
						// link matches regex
						if linkRegex.MatchString(valStr) {
							// and not seen
							if !crawler.seenLinks[valStr] {
								// assumes href value is relative
								out <- "https://imoti.net" + valStr
								crawler.seenLinks[valStr] = true
							}
						}
					} else if !more {
						// No more attributes
						break
					}
				}
			}
		case html.ErrorToken:
			// Exit point
			if z.Err() == io.EOF {
				return nil
			}
			return z.Err()
		}
	}
}
