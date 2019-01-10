package urltools

import (
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetURLFromStr(t *testing.T) {
	validStrURL := "https://monzo.com/"
	validStrURL2 := "https://monzo.com/bla#thisisafragment"
	invalidStrURL := "http:////bla"

	Convey("Given a valid url String, it should return a valid url.URL with data", t, func() {
		url, err := GetURLFromStr(validStrURL, true)

		So(url, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(url.Host, ShouldEqual, "monzo.com")
	})

	Convey("Given a valid url string with fragments, it should return a valid url.URL with data and no fragment", t, func() {
		url, err := GetURLFromStr(validStrURL2, true)

		So(url, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(url.Host, ShouldEqual, "monzo.com")
		So(url.Path, ShouldEqual, "/bla")
		So(url.Fragment, ShouldBeBlank)
	})

	Convey("Given an invalid url String, it should return nil data and error (due to checks)", t, func() {
		url, err := GetURLFromStr(invalidStrURL, true)

		So(url, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})
}

func TestGetPageLinks(t *testing.T) {
	validHTML := []byte(`<html><body><p><a href="hithere">Hi</a></p></body></html>`)
	validHTMLNoHrefs := []byte(`<html><body><p><a f="hithere">Hi</a></p></body></html>`)
	malformedHTML := []byte(`<htm><body<p><a href="hithere">Hi</a><p></body>/html>`)
	malformedHTMLNoHrefs := []byte(`<hml><body<p><a hf="hithere">Hi</a><p></body>/html>`)

	Convey("Given an empty byte array it should return an empty slice", t, func() {
		links := GetPageLinks([]byte{})

		So(len(links), ShouldEqual, 0)
	})

	Convey("Given valid HTML with hrefs, it should return a slice with the found hrefs", t, func() {
		links := GetPageLinks(validHTML)

		So(len(links), ShouldEqual, 1)
		So(links[0], ShouldEqual, "hithere")
	})

	Convey("Given valid HTML with no hrefs, it should return an empty slice", t, func() {
		links := GetPageLinks(validHTMLNoHrefs)

		So(len(links), ShouldEqual, 0)
	})

	Convey("Given malformed HTML (with hrefs), it should return slice with the found hrefs", t, func() {
		links := GetPageLinks(malformedHTML)

		So(len(links), ShouldEqual, 1)
		So(links[0], ShouldEqual, "hithere")
	})

	Convey("Given malformed HTML (with no hrefs), it should return an empty slice", t, func() {
		links := GetPageLinks(malformedHTMLNoHrefs)

		So(len(links), ShouldEqual, 0)
	})
}

func TestParseLink(t *testing.T) {
	currentWebsite, _ := url.Parse("https://monzo.com/")
	validLinkWithFragment := "/i-am-a-page-in-monzo-com#weirdfragment"
	fullValidLink := "https://bla.com/i-am-a-totally-legit-link?what=didntexpectthis"
	fullValidLinkWithWWW := "https://www.bla.com/the-www-should-go-away"
	linkWithNoPath := "https://test.com"

	Convey("Provided with a link starting in '/', it should return the full address in the currentWebsite context", t, func() {
		url, err := ParseLink(validLinkWithFragment, currentWebsite)

		So(url, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(url.Scheme, ShouldEqual, "https")
		So(url.Host, ShouldEqual, "monzo.com")
		So(url.Path, ShouldEqual, "/i-am-a-page-in-monzo-com/")
		So(url.Fragment, ShouldBeEmpty)
	})

	Convey("Provided with a full link, it should return the url.URL for that link", t, func() {
		url, err := ParseLink(fullValidLink, currentWebsite)

		So(url, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(url.Scheme, ShouldEqual, "https")
		So(url.Host, ShouldEqual, "bla.com")
		So(url.Path, ShouldEqual, "/i-am-a-totally-legit-link/")
		So(url.Fragment, ShouldBeEmpty)
		So(url.RawQuery, ShouldBeEmpty)
	})

	Convey("Provided with a full link with `www. at the beginning, it should return the url.URL for that link without `www.`", t, func() {
		url, err := ParseLink(linkWithNoPath, currentWebsite)

		So(url, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})

	Convey("Provided with a full link without path, it should return an error.", t, func() {
		url, err := ParseLink(fullValidLinkWithWWW, currentWebsite)

		So(url, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(url.Host, ShouldEqual, "bla.com")
		So(url.Scheme, ShouldEqual, "https")
		So(url.Path, ShouldEqual, "/the-www-should-go-away/")
		So(url.Fragment, ShouldBeEmpty)
		So(url.RawQuery, ShouldBeEmpty)
	})
}
