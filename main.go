package main

import (
	"flag"
	"log"
)

var urlFlag = flag.String("url", "https://monzo.com/", "Full URL to a website (eg: https://monzo.com/)")

func init() {
	flag.Parse()
}

func main() {
	urlStr := *urlFlag

	url, err := getURLFromStr(urlStr, true)
	if err != nil {
		log.Fatalln("The URL you entered is not valid.")
		return
	}

	wc := NewWebCrawler(url)
	wc.Crawl()
}
