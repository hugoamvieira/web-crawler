package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
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

func getBodyBytes(url string) ([]byte, error) {
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

	return bodyBytes, nil
}
