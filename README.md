# Web Crawler

This software will crawl a domain for its links and display them to stdout.

## Running it

In order to run this program from source code, you should get its dependencies by running `dep ensure`.

After, you can run it directly by executing `go build && ./web-crawler -url <WEBSITE_URL>`

A valid `WEBSITE_URL` is defined by one that has a host, not empty path (at least '/') and 'https' or 'http' for the scheme (eg: `https://monzo.com/`)

## Remarks

By default, theres a depth cap for the URLs (Which means the program will ignore any links that have more than a certain depth - This is calculated based on the amount of '/' in the path). You can override this by passing the flag `-depth` followed by an integer. Beware that it will take longer, though.

This program takes the following steps to deal with some edge cases:

- Every link that has `www.` on it will have that part of it discarded;
- If the link starts with `/`, it will be assumed that it belongs to the current website's domain and details will be added to the URL to reflect that;
- Every link will end with a `/`;
- Every link will have its fragments (`#this`) and query parameters (`?this=that`) removed;
- The hash map storing the seen websites will store a relation of <HOSTNAME + PATH> -> url.URL. This avoids duplicates if the scheme is different (eg: Same page but with `http` and `https`).
