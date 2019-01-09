package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func getURLFromStr(urlStr string, check bool) (*url.URL, error) {
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

func getBodyBytes(url string) (*[]byte, error) {
	var cl http.Client
	resp, err := cl.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Couldn't get content. Website returned status " + resp.Status)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &bodyBytes, nil
}

func getPageLinks(bodyBytes *[]byte) []string {
	var hrefs []string

	r := bytes.NewReader(*bodyBytes)
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
	url, err := getURLFromStr(link, false)
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
