package main

import (
	"flag"
	"log"
)

var urlFlag = flag.String("url", "https://monzo.com/", "Full URL to a website (eg: https://monzo.com/)")
var maxPagesFlag = flag.Int("maxPages", 100, "The amount of pages to visit in total")

func init() {
	flag.Parse()
}

func main() {
	urlStr := *urlFlag
	maxPages := *maxPagesFlag

	url, err := getURLFromStr(urlStr, true)
	if err != nil {
		log.Fatalln("The URL you entered is not valid.")
		return
	}

	wc := NewWebCrawler(url, maxPages)
	wc.Crawl()
}
