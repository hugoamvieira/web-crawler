package main

import (
	"context"
	"log"
	"net/url"
	"sync"

	"github.com/hugoamvieira/web-crawler/urltools"

	"github.com/hugoamvieira/web-crawler/config"
	"github.com/hugoamvieira/web-crawler/datastructures"
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
	q           *datastructures.Queue
	workers     []*webCrawlerWorker
	config      *config.WebCrawlerConfig
	RootWebsite *url.URL

	wg      sync.WaitGroup
	visited sync.Map // Map of URL host+path (string) to bool
}

// NewWebCrawlerV2 creates a new `WebCrawlerV2`, creates the workers
// and returns nil if you set `crawlerWorkerCount` to zero.
func NewWebCrawlerV2(r *url.URL) (*WebCrawlerV2, error) {
	config, err := config.LoadJSONConfig("config/config.json")
	if err != nil {
		return nil, err
	}

	q := datastructures.NewQueue()

	workers := make([]*webCrawlerWorker, config.WorkerCount)
	for i := 0; i < config.WorkerCount; i++ {
		workers[i] = &webCrawlerWorker{
			id:     i,
			q:      q,
			config: config,
		}
	}

	return &WebCrawlerV2{
		q:           q,
		workers:     workers,
		config:      config,
		RootWebsite: r,
	}, nil
}

// Crawl boostraps the web crawling process and starts the workers.
func (wc *WebCrawlerV2) Crawl(ctx context.Context) {
	// Bootstrap the channel with root URL links.
	// This reduces the chances of two workers colliding on the same website.
	// They'd (probably) diverge at some point, but ignoring this step would duplicate
	// ("N-licate" actually) work across workers.
	wc.bootstrap(ctx)

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

func (wc *WebCrawlerV2) bootstrap(ctx context.Context) {
	domainURLs, err := urltools.GetDomainWebsiteURLs(ctx, wc.RootWebsite, wc.config.HTTPTimeout, wc.RootWebsite)
	if err != nil {
		log.Fatalln("Failed to get root website's URLs, cannot continue | Error:", err)
	}

	wc.visited.Store(urltools.GetVisitedMapKey(*wc.RootWebsite), true)

	for _, url := range domainURLs {
		if url == nil {
			continue
		}

		k := urltools.GetVisitedMapKey(*url)
		if _, ok := wc.visited.Load(k); !ok {
			wc.q.Enqueue(url)
		}
	}
}
