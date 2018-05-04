package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

// Validates string and returns valid url.URL by trying to convert the string into Go's URL struct and then checking for some parameters
func getValidURL(urlStr string) (*url.URL, error) {
	url, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	// A valid URL is defined by one that has a host,
	// not empty path (at least '/') and 'https' or 'http' for the scheme
	if url.Host == "" || url.Path == "" || (url.Scheme != "https" && url.Scheme != "http") {
		return nil, errors.New("Invalid URL")
	}

	// Remove fragments as they're not required for this
	url.Fragment = ""
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
	r := bytes.NewReader(*bodyBytes)
	tokenizer := html.NewTokenizer(r)

	var hrefs []string
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
						// Found href inside <a> tag
						hrefs = append(hrefs, a.Val)
						break
					}
				}
			}
		}
	}
}
