package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/hugoamvieira/web-crawler/config"
	"github.com/hugoamvieira/web-crawler/urltools"

	"github.com/hugoamvieira/web-crawler/datastructures"
)

type webCrawlerWorker struct {
	id     int
	q      *datastructures.Queue
	config *config.WebCrawlerConfig
	done   bool
}

// Work starts the worker and it keeps it alive until it gets its context cancelled, or the queue is empty
func (wcw *webCrawlerWorker) Work(ctx context.Context, wg *sync.WaitGroup, visited *sync.Map) {
	log.Printf("Worker ID %v reporting for duty 👨‍🏭\n", wcw.id)
	go wcw.listenForCtxDone(ctx)

	for {
		if wcw.done {
			wg.Done()
			return
		}

		url, err := wcw.q.Dequeue()
		if err == datastructures.ErrEmptyQueue {
			log.Printf("No more websites to look at, worker %v says bye bye 👋\n", wcw.id)
			wcw.done = true // Avoid the ctx goroutine from logging the "Worker spinning down" message at the end.
			wg.Done()
			return
		}
		if url == nil {
			continue
		}

		// Last version of this program did HTTP requests and all the parsing
		// even if the website had been visited - This version fixes that.
		if _, ok := visited.Load(urltools.GetVisitedMapKey(*url)); ok {
			continue
		}

		visited.Store(urltools.GetVisitedMapKey(*url), true)

		fmt.Println("Got website:", url.String())
		domainURLs, err := urltools.GetDomainWebsiteURLs(ctx, url, wcw.config.HTTPTimeout, url)
		if err == urltools.ErrStatusCodeNotOK {
			continue
		}
		if err != nil {
			continue
		}

		for _, u := range domainURLs {
			if _, ok := visited.Load(urltools.GetVisitedMapKey(*u)); !ok {
				wcw.q.Enqueue(u)
			}
		}
	}
}

func (wcw *webCrawlerWorker) listenForCtxDone(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			if wcw.done {
				// Avoid the "No more websites" message being printed if
				// it reaches the end of the queue before the context is ever cancelled.
				return
			}

			log.Printf("Worker %v gracefully spinning down 💃\n", wcw.id)
			wcw.done = true
			return
		}
	}
}
