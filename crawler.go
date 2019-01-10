package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/hugoamvieira/web-crawler/datastructures"
	"github.com/hugoamvieira/web-crawler/urltools"
)

const (
	// defaultHTTPTimeout outlines the timeout for HTTP requests done by the workers.
	defaultHTTPTimeout = 10 * time.Second
	// crawlerWorkerCount outlines the amount of workers that will be spun up.
	crawlerWorkerCount = 2
)

type WebCrawlerV2 struct {
	// Using a queue instead of a channel because while the queue implementation
	// is more expensive when slice expands and slice copies have to happen, it
	// gives us some really good properties:
	// - There is a possibility, when using an unbuffered channel (which blocks for the
	// producer when there's no available receiver for the message), that if N-1 workers block,
	// the program will essentially stop working. It might be low probability,
	// but if that happens, I do not see any other way of recovering from this other
	// than restarting the service, because all workers get in a deadlock state. This is the
	// main reason I'm using a queue, as I don't want to write a program that has
	// this behaviour at its core. When using a queue, if N-1 workers block,
	// the program still continues working with one worker, and it'll "self-heal"
	// when the workers become unblocked.

	// Using an unbuffered channel instead of a queue because if I were to implement a queue
	// with slices (which is the way I'd do it), there'd be a lot of re-assignments,
	// appends, slice expansions, copies, which are expensive.
	q           *datastructures.Queue
	workers     []*webCrawlerWorker
	wg          sync.WaitGroup
	visited     sync.Map
	RootWebsite *url.URL
}

type webCrawlerWorker struct {
	id         int
	q          *datastructures.Queue
	isDraining bool
}

func (wcw *webCrawlerWorker) Work(ctx context.Context, wg *sync.WaitGroup, visited *sync.Map) {
	log.Printf("Worker ID %v reporting for duty 👨‍🏭\n", wcw.id)

	go func() {
		for {
			select {
			case <-ctx.Done():
				if wcw.isDraining {
					// Avoid the "No more websites" message being printed if
					// it reaches the end of the queue before the context is ever cancelled.
					return
				}

				log.Printf("Worker %v gracefully spinning down 💃\n", wcw.id)
				wcw.isDraining = true
				return
			}
		}
	}()

	for {
		url, err := wcw.q.Dequeue()
		if err == datastructures.ErrEmptyQueue {
			log.Printf("No more websites to look at, worker %v says bye bye 👋\n", wcw.id)
			wcw.isDraining = true // Avoid the ctx goroutine from logging the "Worker spinning down" message at the end.
			wg.Done()
			return
		}

		if _, ok := visited.Load(url.String()); ok {
			continue
		}

		fmt.Printf("Visited Website: %v\n", url.String())
		if !wcw.isDraining {
			// Connect to website
			// Pull all links
			// Analyse them
			// Put them in the channel if they haven't been visited
		}
	}
}

func NewWebCrawlerV2(r *url.URL) *WebCrawlerV2 {
	if crawlerWorkerCount < 1 {
		return nil
	}

	q := datastructures.NewQueue()

	workers := make([]*webCrawlerWorker, crawlerWorkerCount)
	for i := 0; i < crawlerWorkerCount; i++ {
		workers[i] = &webCrawlerWorker{
			id: i,
			q:  q,
		}
	}

	return &WebCrawlerV2{
		q:           q,
		workers:     workers,
		RootWebsite: r,
	}
}

func (wc *WebCrawlerV2) Crawl(ctx context.Context) {
	// Bootstrap the channel with root URL links.
	// This reduces the chances of two workers colliding on the same website.
	// They'd (probably) diverge at some point, but ignoring this step would duplicate
	// ("N-licate" actually) work across workers.
	rq, err := http.NewRequest("GET", wc.RootWebsite.String(), nil)
	if err != nil {
		log.Fatalln("Failed to create request, cannot continue | Error:", err)
	}

	// If a request takes longer than what is outlined in `defaultHTTPTimeout`, I think it's safe to assume that
	// either the website is having issues or unreachable - Either way we're not interested anymore.
	// We're also taking into account the main context, so if the user is not interested
	// anymore, the request gets cancelled.
	ctxWithTimeout, cancelTimeoutCtx := context.WithTimeout(ctx, defaultHTTPTimeout)

	rq = rq.WithContext(ctxWithTimeout)

	client := http.DefaultClient
	resp, err := client.Do(rq)
	if err != nil {
		log.Fatalln("Failed to get response from root website, cannot continue | Error:", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.Fatalln("Failed to get successful response from root website, cannot continue | Error:", err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Failed to read body bytes from root website, cannot continue | Error:", err)
	}

	resp.Body.Close()

	// Context w/ the timeout is not necessary anymore, so we cancel it to avoid a leak.
	cancelTimeoutCtx()

	wc.visited.Store(wc.RootWebsite.String(), true)

	for _, pageLink := range urltools.GetPageLinks(bodyBytes) {
		url, err := urltools.ParseLink(pageLink, wc.RootWebsite)
		if err != nil {
			continue
		}
		if !strings.HasSuffix(url.Host, wc.RootWebsite.Host) {
			// We discard any URL that does not belong to the domain
			continue
		}

		if _, ok := wc.visited.Load(url.String()); !ok {
			log.Printf("Adding %v to channel", url.String())
			wc.q.Enqueue(url)
		}
	}

	// Start workers
	for _, worker := range wc.workers {
		if worker == nil {
			log.Printf("No worker, ignoring")
			continue
		}

		// This should be done outside the goroutine because there's a chance
		// that the Wait() call (outside the loop) will be done before
		// any of the goroutines actually start. In that situation we're waiting
		// on a WaitGroup of 0, which means it returns immediately.
		wc.wg.Add(1)

		go func(c context.Context, w *webCrawlerWorker) {
			w.Work(c, &wc.wg, &wc.visited)
		}(ctx, worker)
	}

	wc.wg.Wait()
	log.Println("All workers spun down")
}
