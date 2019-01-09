package main

import (
	"fmt"
	"net/url"
	"strings"
)

type WebCrawler struct {
	RootWebsite *url.URL
}

// NewWebCrawler returns a new WebCrawler "object" with the parameters passed to the function.
func NewWebCrawler(r *url.URL) *WebCrawler {
	return &WebCrawler{
		RootWebsite: r,
	}
}

// Crawl grabs the root website passed in the object and
// will crawl through pages in the same domain
func (wc *WebCrawler) Crawl() {
	var q []*url.URL
	visited := make(map[string]*url.URL)

	visited[wc.RootWebsite.Host+wc.RootWebsite.Path] = wc.RootWebsite
	q = append(q, wc.RootWebsite)
	domain := wc.RootWebsite.Host

	for len(q) != 0 {
		currentWebsite := q[0]
		q = q[1:] // "Dequeue" first element

		bodyBytes, err := getBodyBytes(currentWebsite.String())
		if err != nil {
			continue
		}

		for _, linkStr := range getPageLinks(bodyBytes) {
			parsedURL, err := parseLink(linkStr, currentWebsite)
			if err != nil {
				continue
			}

			if !strings.HasSuffix(parsedURL.Host, domain) {
				// We discard any URL that does not belong to the domain
				continue
			}

			//The hash map storing the seen websites will store a relation
			//of <HOSTNAME + PATH> -> url.URL. This avoids duplicates if the
			// scheme is different (eg: Same page but with `http` and `https`).
			k := parsedURL.Host + parsedURL.Path
			if _, ok := visited[k]; !ok {
				visited[k] = parsedURL
				q = append(q, parsedURL)
				fmt.Println(parsedURL)
			}
		}
	}
}
