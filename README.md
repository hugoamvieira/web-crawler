# Web Crawler

This software will crawl a webpage for its links and display them to stdout.

## Running it

In order to run this program from source code, you should get its dependencies by running `dep ensure`.

After, you can run it directly by executing `go run *.go -url <WEBSITE_URL>`

A valid `WEBSITE_URL` is defined by one that has a host, not empty path (at least '/') and 'https' or 'http' for the scheme (eg: `https://monzo.com/`)