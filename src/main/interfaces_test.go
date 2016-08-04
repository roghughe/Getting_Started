package main

import (
	"testing"
)

var infoOutput = []byte(`{ "Environment": "production" }`)

var statusOutput = []byte(`{ "Status": "up" }`)

type stubFetcher struct{}

func (fetcher stubFetcher) Fetch(url string) ([]byte, error) {
	if strings.Contains(url, "/info") {
		return infoOutput, nil
	}

	if strings.Contains(url, "/status") {
		return statusOutput, nil
	}

	return nil, errors.New("Don't recognize URL: " + url)
}

var info *Info
var stub stubFetcher

func TestXYZ(t *testing.T) {
	// We would make some assertions around this:
	populateInfo(stubFetcher, info)
}
