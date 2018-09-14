package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// WebCrawler is the object that stores all information to start crawling a particular domain
type WebCrawler struct {
	RootWebsite     *url.URL
	MaxPagesToVisit int
}

// NewWebCrawler returns a new WebCrawler "object" with the parameters passed to the function.
func NewWebCrawler(r *url.URL, maxPages int) *WebCrawler {
	return &WebCrawler{
		RootWebsite:     r,
		MaxPagesToVisit: maxPages,
	}
}

// Crawl grabs the root website passed in the object and
// will recursively crawl through pages in the same domain
func (wc *WebCrawler) Crawl() {
	var q []*url.URL
	visited := make(map[string]*url.URL)
	pagesVisited := 0

	visited[wc.RootWebsite.Host+wc.RootWebsite.Path] = wc.RootWebsite

	q = append(q, wc.RootWebsite)
	domain := wc.RootWebsite.Host

	for len(q) != 0 && (pagesVisited < wc.MaxPagesToVisit) {
		pagesVisited++
		currentWebsite := q[0]
		q = q[1:] // "Dequeue" first element

		fmt.Printf("Analysing %v\n", currentWebsite.String())
		bodyBytes, err := getBodyBytes(currentWebsite.String())
		if err != nil {
			continue
		}

		title, err := GetPageTitle(bodyBytes)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("Page Title: %v\n", *title)

		fmt.Println("Static Assets:")
		for _, staticAssetSrc := range GetPageStaticAssets(bodyBytes) {
			fmt.Printf("\t%v\n", staticAssetSrc)
		}

		fmt.Println("Links:")
		for _, linkStr := range GetPageLinks(bodyBytes) {
			parsedURL, err := parseLink(linkStr, currentWebsite)
			if err != nil {
				continue
			}

			fmt.Printf("\t%v\n", parsedURL)

			// We discard any URL that does not belong to the domain from the upcoming logic
			if !strings.HasSuffix(parsedURL.Host, domain) {
				continue
			}

			//The hash map storing the seen websites will store a relation
			//of <HOSTNAME + PATH> -> url.URL. This avoids duplicates if the
			// scheme is different (eg: Same page but with `http` and `https`).
			k := parsedURL.Host + parsedURL.Path
			if _, ok := visited[k]; !ok {
				visited[k] = parsedURL
				q = append(q, parsedURL)
			}
		}

		fmt.Println() // This is just a new line for better CLI readability
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
