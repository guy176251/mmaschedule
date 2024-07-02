package scraper

import (
	_ "embed"
	"testing"
)

//go:embed testdata/tapology-search.html
var search_result []byte

func testTapologyGetter(t *testing.T) {
	get_tapology_for := TapologyGetter()
	PrintJSON(get_tapology_for("Dustin Poirier"))
	PrintJSON(get_tapology_for("Jon Jones"))
	PrintJSON(get_tapology_for("Conor Mcgregor"))
}

func TestTapologyParseFragment(t *testing.T) {
    PrintJSON(parse_tapology_results(search_result))
}
