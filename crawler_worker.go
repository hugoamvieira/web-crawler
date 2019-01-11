package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/hugoamvieira/web-crawler/config"

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

		if _, ok := visited.Load(url); ok {
			continue
		}

		fmt.Printf("Got website: %v\n", url.String())

		// TODO
		// Connect to website
		// Pull all links
		// Analyse them
		// Put them in the queue if they haven't been visited
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
