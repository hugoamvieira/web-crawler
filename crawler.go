package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/hugoamvieira/web-crawler/urltools"
)

const (
	// DefaultHTTPTimeout outlines the timeout for HTTP requests done by the workers.
	DefaultHTTPTimeout = 10 * time.Second
)

type WebCrawlerV2 struct {
	c           chan *url.URL
	workers     []*webCrawlerWorker
	wg          sync.WaitGroup
	visited     sync.Map
	RootWebsite *url.URL
}

type webCrawlerWorker struct {
	id int
	c  chan *url.URL
}

func (wcw *webCrawlerWorker) Work(ctx context.Context, wg *sync.WaitGroup, visited *sync.Map) {
	log.Printf("Worker ID %v reporting for duty 👨‍🏭\n", wcw.id)

	for {
		select {
		case url := <-wcw.c:
			if _, ok := visited.Load(url.String()); ok {
				continue
			}

			// Connect to website
			// Pull all links
			// Analyse them
			// Put them in the channel if they haven't been visited

		case <-ctx.Done():
			log.Printf("Context cancelled, worker %v spinning down\n", wcw.id)
			wg.Done()
			return
		}
	}
}

func NewWebCrawlerV2(r *url.URL, workerCount int) *WebCrawlerV2 {
	if workerCount < 1 {
		return nil
	}

	c := make(chan *url.URL)

	workers := make([]*webCrawlerWorker, workerCount)
	for i := 0; i < workerCount; i++ {
		workers[i] = &webCrawlerWorker{
			id: i,
			c:  c,
		}
	}

	return &WebCrawlerV2{
		c:           c,
		workers:     workers,
		RootWebsite: r,
	}
}

func (wc *WebCrawlerV2) Crawl(ctx context.Context) {
	// Bootstrap the channel with root URL links. This ensures that workers go "their separate ways".
	// They'd (probably) diverge at some point, but ignoring this step would duplicate
	// effort across workers.
	rq, err := http.NewRequest("GET", wc.RootWebsite.String(), nil)
	if err != nil {
		log.Fatalln("Failed to create request, cannot continue | Error:", err)
	}

	// If a request takes longer than what is outlined in DefaultHTTPTimeout, I think it's safe to assume that
	// either the website is having issues or unreachable - Either way we're not interested anymore.
	// We're also injecting this onto the main context, so if the user is not interested
	// anymore, the request gets cancelled
	ctxWithTimeout, cancelCtx := context.WithTimeout(ctx, DefaultHTTPTimeout)
	defer cancelCtx()

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

	wc.visited.Store(wc.RootWebsite.String(), true)

	for _, pageLink := range urltools.GetPageLinks(bodyBytes) {
		url, err := urltools.ParseLink(pageLink, wc.RootWebsite)
		if err != nil {
			log.Println("Failed to parse link from root website | Error:", err)
			continue
		}
		if !strings.HasSuffix(url.Host, wc.RootWebsite.Host) {
			// We discard any URL that does not belong to the domain
			log.Println("URL is not in domain, ignoring")
			continue
		}

		if _, ok := wc.visited.Load(url.String()); !ok {
			log.Printf("Adding %v to channel", url.String())
			go func() {
				wc.c <- url
			}()
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
