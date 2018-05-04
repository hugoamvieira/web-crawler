package main

import (
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseLink(t *testing.T) {
	currentWebsite, _ := url.Parse("https://monzo.com/")
	validLinkWithFragment := "/i-am-a-page-in-monzo-com#weirdfragment"
	fullValidLink := "https://bla.com/i-am-a-totally-legit-link?what=didntexpectthis"
	fullValidLinkWithWWW := "https://www.bla.com/the-www-should-go-away"
	linkWithNoPath := "https://test.com"

	Convey("Provided with a link starting in '/', it should return the full address in the currentWebsite context", t, func() {
		url, err := parseLink(validLinkWithFragment, currentWebsite)

		So(url, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(url.Scheme, ShouldEqual, "https")
		So(url.Host, ShouldEqual, "monzo.com")
		So(url.Path, ShouldEqual, "/i-am-a-page-in-monzo-com/")
		So(url.Fragment, ShouldBeEmpty)
	})

	Convey("Provided with a full link, it should return the url.URL for that link", t, func() {
		url, err := parseLink(fullValidLink, currentWebsite)

		So(url, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(url.Scheme, ShouldEqual, "https")
		So(url.Host, ShouldEqual, "bla.com")
		So(url.Path, ShouldEqual, "/i-am-a-totally-legit-link/")
		So(url.Fragment, ShouldBeEmpty)
		So(url.RawQuery, ShouldBeEmpty)
	})

	Convey("Provided with a full link with `www. at the beginning, it should return the url.URL for that link without `www.`", t, func() {
		url, err := parseLink(linkWithNoPath, currentWebsite)

		So(url, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})

	Convey("Provided with a full link without path, it should return an error.", t, func() {
		url, err := parseLink(fullValidLinkWithWWW, currentWebsite)

		So(url, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(url.Host, ShouldEqual, "bla.com")
		So(url.Scheme, ShouldEqual, "https")
		So(url.Path, ShouldEqual, "/the-www-should-go-away/")
		So(url.Fragment, ShouldBeEmpty)
		So(url.RawQuery, ShouldBeEmpty)
	})
}

func TestGetPathDepth(t *testing.T) {
	// All paths will have a '/' artificially added at the end
	validPathDepth6 := "/hello/hi/this/is/totally/valid/"
	validPathDepth1 := "/hi/"
	invalidPath := ""
	Convey("Provided with a valid path, should return its length and nil error", t, func() {
		depth, err := getPathDepth(validPathDepth6)

		So(err, ShouldBeNil)
		So(*depth, ShouldEqual, 6)

	})

	Convey("Provided with a valid path, should return its length and nil error", t, func() {
		depth, err := getPathDepth(validPathDepth1)

		So(err, ShouldBeNil)
		So(*depth, ShouldEqual, 1)

	})

	Convey("Provided with an invalid path, should nil depth and error", t, func() {
		depth, err := getPathDepth(invalidPath)

		So(err, ShouldNotBeNil)
		So(depth, ShouldBeNil)

	})
}
