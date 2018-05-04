package main

import (
	"bytes"

	"golang.org/x/net/html"
)

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
