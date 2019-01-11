package urltools

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var (
	// ErrStatusCodeNotOK is returned if an HTTP request's status code is < 200 or > 299.
	ErrStatusCodeNotOK = errors.New("Status code was not in the 200 range for this request")
	// ErrNilURLValues is returned when you pass nil values into the GetDomainWebsiteURLs func.
	ErrNilURLValues = errors.New("One or more URLs passed are nil, cannot continue")
)

func GetURLFromStr(urlStr string, check bool) (*url.URL, error) {
	url, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	if check {
		// A valid URL is defined by one that has a host,
		// not empty path (at least '/') and 'https' or 'http' for the scheme
		if url.Host == "" || url.Path == "" || (url.Scheme != "https" && url.Scheme != "http") {
			return nil, errors.New("Invalid URL")
		}
	}

	// Remove fragments and query parameters as we're trying to determine uniqueness
	url.Fragment = ""
	url.RawQuery = ""
	return url, nil
}

// GetDomainWebsiteURLs receives a context, the URL it'll get the content from, the time-out for the HTTP request,
// the current website you're on and it returns any URLs found on that page's content for that domain.
// The last parameter (currentWebsite) is used to build up new links that are in the domain but
// not defined completely (eg: href="/about").
func GetDomainWebsiteURLs(ctx context.Context, u *url.URL, timeout time.Duration, currentWebsite *url.URL) ([]*url.URL, error) {
	if u == nil || currentWebsite == nil {
		// No need to validate timeout - already been done on config load.
		return nil, ErrNilURLValues
	}

	rq, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	// If a request takes longer than what is outlined in `defaultHTTPTimeout`, I think it's safe to assume that
	// either the website is having issues or unreachable - Either way we're not interested anymore.
	// We're also taking into account the main context, so if the user is not interested
	// anymore, the request gets cancelled.
	ctxWithTimeout, cancelTimeoutCtx := context.WithTimeout(ctx, timeout)
	defer cancelTimeoutCtx()

	rq = rq.WithContext(ctxWithTimeout)

	client := http.DefaultClient
	resp, err := client.Do(rq)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, ErrStatusCodeNotOK
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	urls := make([]*url.URL, 0)
	for _, link := range getPageLinks(bodyBytes) {
		url, err := parseLink(link, currentWebsite)
		if err != nil {
			continue
		}
		if !strings.HasSuffix(url.Host, currentWebsite.Host) {
			// We discard any URL that does not belong to the domain
			continue
		}

		urls = append(urls, url)
	}

	return urls, nil
}

func getPageLinks(bodyBytes []byte) []string {
	var hrefs []string

	r := bytes.NewReader(bodyBytes)
	tokenizer := html.NewTokenizer(r)
	for {
		tokType := tokenizer.Next()

		switch {
		case tokType == html.ErrorToken:
			return hrefs
		case tokType == html.StartTagToken:
			tok := tokenizer.Token()
			if tok.Data == "a" {
				for _, a := range tok.Attr {
					if a.Key == "href" {
						hrefs = append(hrefs, a.Val)
						break
					}
				}
			}
		}
	}
}

func parseLink(link string, currentWebsite *url.URL) (*url.URL, error) {
	url, err := GetURLFromStr(link, false)
	if err != nil {
		return nil, err
	}

	if len(url.Path) == 0 {
		return nil, errors.New("This link has no path")
	}

	if strings.HasPrefix(url.Host, "www.") {
		// Trim www. out of the host for standardisation purposes
		url.Host = strings.TrimPrefix(url.Host, "www.")
	}

	firstURLChar := url.String()[0]
	if string(firstURLChar) == "/" {
		// If the whole link starts with a slash, it is very likely that it
		// belongs to the same domain, so we'll append some details to it.
		url.Scheme = currentWebsite.Scheme
		url.Host = currentWebsite.Host
	}

	pathLen := len(url.Path)
	lastPathChar := string(url.Path[pathLen-1])
	if lastPathChar != "/" {
		// If the last character in the Path is not a slash, add it to avoid
		// having duplicate cases such as https://monzo.com/about and https://monzo.com/about/
		url.Path += "/"
	}

	return url, nil
}

func GetVisitedMapKey(url url.URL) string {
	return url.Host + url.Path
}
