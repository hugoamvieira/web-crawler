package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type WebCrawler struct {
	RootWebsite *url.URL
	MaxURLDepth int
}

// NewWebCrawler returns a new WebCrawler "object" with the parameters passed to the function.
func NewWebCrawler(r *url.URL, d int) *WebCrawler {
	return &WebCrawler{
		RootWebsite: r,
		MaxURLDepth: d,
	}
}

// Crawl grabs the root website passed in the object and
// will recursively crawl through pages in the same domain
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

			depth, err := getPathDepth(parsedURL.Path)
			if err != nil {
				continue
			}

			if *depth > wc.MaxURLDepth {
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

func parseLink(link string, currentWebsite *url.URL) (*url.URL, error) {
	url, err := getURLFromStr(link, false)
	if err != nil {
		return nil, err
	}

	if len(url.Path) == 0 {
		return nil, errors.New("This link has no path")
	}

	if strings.HasPrefix(url.Host, "www.") {
		// Trim www. out of the host for standardisation purposes
		url.Host = strings.TrimPrefix(url.Host, "www.")
	}

	firstURLChar := url.String()[0]
	if string(firstURLChar) == "/" {
		// If the whole link starts with a slash, it is very likely that it
		// belongs to the same domain, so we'll append some details to it.
		url.Scheme = currentWebsite.Scheme
		url.Host = currentWebsite.Host
	}

	pathLen := len(url.Path)
	lastPathChar := string(url.Path[pathLen-1])
	if lastPathChar != "/" {
		// If the last character in the Path is not a slash, add it to avoid
		// having duplicate cases such as https://monzo.com/about and https://monzo.com/about/
		url.Path += "/"
	}

	return url, nil
}

// Depth is calculated by the amount of slashes in the URL path.
func getPathDepth(path string) (*int, error) {
	paths := strings.Split(path, "/")
	depth := len(paths) - 1

	if depth == 0 {
		return nil, errors.New("Path is empty")
	}

	// Paths will at least have one '/' which was artificially inserted to standardise the URLs
	ret := depth - 1
	return &ret, nil
}
