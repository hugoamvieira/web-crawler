package main

import (
	"flag"
	"fmt"
	"log"
)

var urlFlag = flag.String("url", "", "Full URL to a website (eg: https://monzo.com/)")

func init() {
	flag.Parse()
}

func main() {
	urlStr := *urlFlag

	url, err := getValidURL(urlStr)
	if err != nil {
		log.Fatalln("The URL you entered is not valid")
	}

	bodyBytes, err := getBodyBytes(url.String())
	if err != nil {
		log.Fatalln("Couldn't get website body")
	}

	links := getPageLinks(bodyBytes)
	var domainLinks []string

	for _, link := range links {
		if len(link) == 0 {
			continue
		}

		if string(link[0]) == "/" {
			// If the link starts with a slash, it is very likely that it belongs to the same domain, so we'll append that
			link = url.Scheme + "://" + url.Host + link
			domainLinks = append(domainLinks, link)
		} else {
			validatedURL, err := getValidURL(link)
			if err != nil {
				continue
			}
			if validatedURL.Host == url.Host {
				domainLinks = append(domainLinks, validatedURL.String())
			}
		}
	}

	for _, link := range domainLinks {
		fmt.Println(link)
	}
}

func main2() {
	urlStr := *urlFlag

	url, err := getValidURL(urlStr)
	if err != nil {
		log.Fatalln("The URL you entered is not valid")
		return
	}

	getDomainMap(url)
}
