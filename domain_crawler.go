package main

import (
	"net/url"
)

func getDomainMap(rootWebsite *url.URL) {
	// Should move away from slices for a Queue data structure to avoid reallocation
	var websiteQueue []*url.URL
	visitedWebsites := make(map[string]*url.URL)

	visitedWebsites[rootWebsite.String()] = rootWebsite
	websiteQueue = append(websiteQueue, rootWebsite)
	searchableDomain := rootWebsite.Host

	for len(websiteQueue) != 0 {
		currentWebsite := websiteQueue[0]
		websiteQueue = websiteQueue[1:]

		bodyBytes, err := getBodyBytes(currentWebsite.String())
		if err != nil {
			continue
		}

		pageLinks := getPageLinks(bodyBytes)
		for _, link := range pageLinks {
			if len(link) == 0 {
				continue
			}
			if string(link[0]) == "/" {
				// If the link starts with a slash, it is very likely that it belongs to the same domain, so we'll append that
				link = currentWebsite.Scheme + "://" + currentWebsite.Host + link
			}

			pageURL, err := getValidURL(link)
			if err != nil {
				continue
			}
			if pageURL.Host != searchableDomain {
				continue
			}

			pageURLString := pageURL.String()
			if _, ok := visitedWebsites[pageURLString]; !ok {
				visitedWebsites[pageURLString] = pageURL
				websiteQueue = append(websiteQueue, pageURL)
			}
		}
	}
}
