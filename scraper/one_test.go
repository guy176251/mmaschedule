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
	print_json(parse_one_event_list(ONEEventListHTML))
}

func TestONEEvent(t *testing.T) {
	print_json(parse_one_event("some_url", ONEEventHTML))
}

func testScrapeONE(t *testing.T) {
	scrape_one(func(e *Event) {
		print_json(e)
	})
}
