package main

import (
	"bytes"
	"errors"

	"golang.org/x/net/html"
)

/*
These three functions could possibly be flattened to one that traverses the HTML
and finds everything we need, however it would make it larger, more complex and harder to test.
I like readable code, so I decided to leave them as is...

Additionally, (O(N) * 3) still amortizes to O(N), so doing 3 for-loops
instead of one doesn't really affect performance that much.
*/

// GetPageStaticAssets receives some HTML in []byte form and returns a list of present static assets.
func GetPageStaticAssets(bodyBytes []byte) []string {
	staticAssets := make([]string, 0)

	r := bytes.NewReader(bodyBytes)
	tokenizer := html.NewTokenizer(r)
	for {
		tokType := tokenizer.Next()

		switch {
		case tokType == html.ErrorToken:
			return staticAssets
		case tokType == html.StartTagToken:
			tok := tokenizer.Token()

			// Assuming that "static assets" == images
			if tok.Data == "img" {
				for _, img := range tok.Attr {
					if img.Key == "src" && img.Val != "" {
						staticAssets = append(staticAssets, img.Val)
						break
					}
				}
			}
		}
	}
}

// GetPageLinks receives some HTML in []byte form and returns a list of both present internal and external links.
func GetPageLinks(bodyBytes []byte) []string {
	hrefs := make([]string, 0)

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
					if a.Key == "href" && a.Val != "" {
						hrefs = append(hrefs, a.Val)
						break
					}
				}
			}
		}
	}
}

/*
GetPageTitle receives some HTML in []byte form and returns a string pointer with the title.
It may error if it cannot parse the []byte to HTML or if the page doesn't have a title.
*/
func GetPageTitle(bodyBytes []byte) (*string, error) {
	r := bytes.NewReader(bodyBytes)
	n, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	title := traverseHTMLForTitle(n)
	if title != nil {
		return title, nil
	}

	return nil, errors.New("Couldn't find page title")
}

func traverseHTMLForTitle(n *html.Node) *string {
	if n.Type == html.ElementNode && n.Data == "title" {
		return &n.FirstChild.Data
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		title := traverseHTMLForTitle(c)
		if title != nil {
			return title
		}
	}

	return nil
}
