package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLinkFetcher(t *testing.T) {
	validHTML := []byte(`<html><body><p><a href="hithere">Hi</a></p></body></html>`)
	validHTMLNoHrefs := []byte(`<html><body><p><a f="hithere">Hi</a></p></body></html>`)
	malformedHTML := []byte(`<htm><body<p><a href="hithere">Hi</a><p></body>/html>`)
	malformedHTMLNoHrefs := []byte(`<hml><body<p><a hf="hithere">Hi</a><p></body>/html>`)

	Convey("Given an empty byte array it should return an empty slice", t, func() {
		links := getPageLinks(&[]byte{})

		So(len(links), ShouldEqual, 0)
	})

	Convey("Given valid HTML with hrefs, it should return a slice with the found hrefs", t, func() {
		links := getPageLinks(&validHTML)

		So(len(links), ShouldEqual, 1)
		So(links[0], ShouldEqual, "hithere")
	})

	Convey("Given valid HTML with no hrefs, it should return an empty slice", t, func() {
		links := getPageLinks(&validHTMLNoHrefs)

		So(len(links), ShouldEqual, 0)
	})

	Convey("Given malformed HTML (with hrefs), it should return slice with the found hrefs", t, func() {
		links := getPageLinks(&malformedHTML)

		So(len(links), ShouldEqual, 1)
		So(links[0], ShouldEqual, "hithere")
	})

	Convey("Given malformed HTML (with no hrefs), it should return an empty slice", t, func() {
		links := getPageLinks(&malformedHTMLNoHrefs)

		So(len(links), ShouldEqual, 0)
	})
}
