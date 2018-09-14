package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLinkFetcher(t *testing.T) {
	validHTML := []byte(`<html><body><p><a href="hithere">Hi</a></p></body></html>`)
	Convey("Given valid HTML with hrefs, it should return a slice with the found hrefs", t, func() {
		links := GetPageLinks(validHTML)

		So(links, ShouldNotBeNil)
		So(links, ShouldHaveLength, 1)
		So(links[0], ShouldEqual, "hithere")
	})

	validHTMLNoHrefs := []byte(`<html><body><p><a f="hithere">Hi</a></p></body></html>`)
	Convey("Given valid HTML with no hrefs, it should return an empty slice", t, func() {
		links := GetPageLinks(validHTMLNoHrefs)

		So(links, ShouldNotBeNil)
		So(links, ShouldHaveLength, 0)
	})

	malformedHTML := []byte(`<htm><body<p><a href="hithere">Hi</a><p></body>/html>`)
	Convey("Given malformed HTML (with hrefs), it should return slice with the found hrefs", t, func() {
		links := GetPageLinks(malformedHTML)

		So(links, ShouldNotBeNil)
		So(links, ShouldHaveLength, 1)
		So(links[0], ShouldEqual, "hithere")
	})

	malformedHTMLNoHrefs := []byte(`<hml><body<p><a hf="hithere">Hi</a><p></body>/html>`)
	Convey("Given malformed HTML (with no hrefs), it should return an empty slice", t, func() {
		links := GetPageLinks(malformedHTMLNoHrefs)

		So(links, ShouldNotBeNil)
		So(links, ShouldHaveLength, 0)
	})

	Convey("Given an empty byte array it should return an empty slice", t, func() {
		links := GetPageLinks([]byte{})

		So(links, ShouldNotBeNil)
		So(links, ShouldHaveLength, 0)
	})

	Convey("Given a nil byte slice, it should return an empty slice", t, func() {
		links := GetPageLinks(nil)

		So(links, ShouldNotBeNil)
		So(links, ShouldHaveLength, 0)
	})
}

func TestTitleFetcher(t *testing.T) {
	validHTMLWithTitle := []byte(`<html><head><title>Hi</title></head><body><p><a href="hithere">Hi</a></p></body></html>`)
	Convey("Given valid HTML with title, it should return that title and no error", t, func() {
		title, err := GetPageTitle(validHTMLWithTitle)

		So(err, ShouldBeNil)
		So(title, ShouldNotBeNil)
	})

	validHTMLNoTitle := []byte(`<html><body><p><a href="hithere">Hi</a></p></body></html>`)
	Convey("Given valid HTML with no title, it should return an error and title should be nil", t, func() {
		title, err := GetPageTitle(validHTMLNoTitle)

		So(err, ShouldNotBeNil)
		So(title, ShouldBeNil)
	})

	malformedHTML := []byte(`<htm><body<p><a href="hithere">Hi</a><p></body>/html>`)
	Convey("Given invalid HTML, it should return an error and title should be nil", t, func() {
		title, err := GetPageTitle(malformedHTML)

		So(err, ShouldNotBeNil)
		So(title, ShouldBeNil)
	})

	Convey("Given a nil byte slice, it should return an error", t, func() {
		title, err := GetPageTitle(nil)

		So(err, ShouldNotBeNil)
		So(title, ShouldBeNil)
	})
}

func TestStaticsFetcher(t *testing.T) {
	validHTMLWithImgTag := []byte(`<html><body><img src="hi.jpg"></img></body></html>`)
	Convey("Given valid HTML with an `img` tag, it should return it's `src`", t, func() {
		statics := GetPageStaticAssets(validHTMLWithImgTag)

		So(statics, ShouldNotBeNil)
		So(statics, ShouldHaveLength, 1)
		So(statics[0], ShouldEqual, "hi.jpg")
	})

	validHTMLWithImgTagNoSrc := []byte(`<html><body><img src=""></img></body></html>`)
	Convey("Given valid html with an `img` tag with an empty `src`, it should return an empty slice", t, func() {
		statics := GetPageStaticAssets(validHTMLWithImgTagNoSrc)

		So(statics, ShouldNotBeNil)
		So(statics, ShouldHaveLength, 0)
	})

	validHTMLWithEmptyImgTag := []byte(`<html><body><img></img></body></html>`)
	Convey("Given valid html with an empty `img` tag, it should return an empty slice", t, func() {
		statics := GetPageStaticAssets(validHTMLWithEmptyImgTag)

		So(statics, ShouldNotBeNil)
		So(statics, ShouldHaveLength, 0)
	})

	validHTMLWithNoImgTag := []byte(`<html><body><p>Hi</p></body></html>`)
	Convey("Given valid HTML with no `img` tag, it should an empty slice", t, func() {
		statics := GetPageStaticAssets(validHTMLWithNoImgTag)

		So(statics, ShouldNotBeNil)
		So(statics, ShouldHaveLength, 0)
	})

	malformedHTML := []byte(`<htm><body<p><a href="hithere">Hi</a><p></body>/html>`)
	Convey("Given invalid HTML, it should return an empty slice", t, func() {
		statics := GetPageStaticAssets(malformedHTML)

		So(statics, ShouldNotBeNil)
		So(statics, ShouldHaveLength, 0)
	})

	Convey("Given a nil byte slice, it should return an empty slice", t, func() {
		statics := GetPageStaticAssets(nil)

		So(statics, ShouldNotBeNil)
		So(statics, ShouldHaveLength, 0)
	})
}
