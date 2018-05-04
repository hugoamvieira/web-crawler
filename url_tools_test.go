package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestURLTools(t *testing.T) {
	validStrURL := "https://monzo.com/"
	validStrURL2 := "https://monzo.com/bla#thisisafragment"
	invalidStrURL := "http:////bla"

	Convey("Given a valid url String, it should return a valid url.URL with data", t, func() {
		url, err := getURLFromStr(validStrURL, true)

		So(url, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(url.Host, ShouldEqual, "monzo.com")
	})

	Convey("Given a valid url string with fragments, it should return a valid url.URL with data and no fragment", t, func() {
		url, err := getURLFromStr(validStrURL2, true)

		So(url, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(url.Host, ShouldEqual, "monzo.com")
		So(url.Path, ShouldEqual, "/bla")
		So(url.Fragment, ShouldBeBlank)
	})

	Convey("Given an invalid url String, it should return nil data and error (due to checks)", t, func() {
		url, err := getURLFromStr(invalidStrURL, true)

		So(url, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})
}
