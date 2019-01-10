package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hugoamvieira/web-crawler/urltools"
)

var urlFlag = flag.String("url", "https://monzo.com/", "Full URL to a website (eg: https://monzo.com/)")

func init() {
	flag.Parse()
}

func main() {
	urlStr := *urlFlag

	url, err := urltools.GetURLFromStr(urlStr, true)
	if err != nil {
		log.Fatalln("The URL you entered is not valid.")
		return
	}

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	go func() {
		var osC = make(chan os.Signal)
		signal.Notify(osC, syscall.SIGTERM)
		signal.Notify(osC, syscall.SIGINT)

		// Wait for signal
		sig := <-osC
		log.Printf("Received %v OS Signal, starting cleanup", sig.String())
		cancelCtx()
	}()

	wc := NewWebCrawlerV2(url)
	if wc == nil {
		log.Fatalln("Web Crawler has not been created, cannot continue (Check your config)")
	}
	wc.Crawl(ctx)

	log.Println("Ta-da! 🌟")
}
