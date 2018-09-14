package main

import (
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseLink(t *testing.T) {
	currentWebsite, _ := url.Parse("https://monzo.com/")

	validLinkWithFragment := "/i-am-a-page-in-monzo-com#weirdfragment"
	Convey("Provided with a link starting in '/', it should return the full address in the currentWebsite context", t, func() {
		url, err := parseLink(validLinkWithFragment, currentWebsite)

		So(url, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(url.Scheme, ShouldEqual, "https")
		So(url.Host, ShouldEqual, "monzo.com")
		So(url.Path, ShouldEqual, "/i-am-a-page-in-monzo-com/")
		So(url.Fragment, ShouldBeEmpty)
	})

	fullValidLink := "https://bla.com/i-am-a-totally-legit-link?what=didntexpectthis"
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

	linkWithNoPath := "https://test.com"
	Convey("Provided with a full link with `www. at the beginning, it should return the url.URL for that link without `www.`", t, func() {
		url, err := parseLink(linkWithNoPath, currentWebsite)

		So(url, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})

	fullValidLinkWithWWW := "https://www.bla.com/the-www-should-go-away"
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
