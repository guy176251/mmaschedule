package scraper

import (
	_ "embed"
	"testing"
)

//go:embed testdata/one-fight-night-23.html
var ONEEventHTML []byte

//go:embed testdata/one-events.html
var ONEEventListHTML []byte

func TestONEEventList(t *testing.T) {
	PrintJSON(ParseONEEventList(ONEEventListHTML))
}

func TestONEEvent(t *testing.T) {
	PrintJSON(ParseONEEvent("some_url", ONEEventHTML))
}

func testScrapeONE(t *testing.T) {
	ScrapeONE(func(e *Event) {
		PrintJSON(e)
	})
}
