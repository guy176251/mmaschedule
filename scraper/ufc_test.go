package scraper

import (
	_ "embed"
	"testing"
)

//go:embed testdata/ufc-297-good.html
var UFCEventHTML []byte

//go:embed testdata/ufc-events.html
var UFCEventListHTML []byte

func TestUFCEventList(t *testing.T) {
	PrintJSON(parse_ufc_event_list(UFCEventListHTML))
}

func TestUFCEvent(t *testing.T) {
	PrintJSON(parse_ufc_event("https://ufc.com/event/ufc-297", UFCEventHTML))
}

func TestParseUFCSlug(t *testing.T)  {
    result := parse_ufc_slug("https://ufc.com/event/ufc-fight-night-july-20-2024")
    expected := "ufc-fight-night-july-20-2024"
    if result != expected {
        PrintJSON(result)
        t.FailNow()
    }
}

func testScrapeUFC(t *testing.T) {
	ScrapeUFC(func(e *Event) {
		PrintJSON(e)
	})
}
