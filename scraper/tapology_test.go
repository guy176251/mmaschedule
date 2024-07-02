package scraper

import (
	_ "embed"
	"testing"
)

//go:embed testdata/tapology-search.html
var search_result []byte

func testTapologyGetter(t *testing.T) {
	get_tapology_for := tapology_getter()
	print_json(get_tapology_for("Dustin Poirier"))
	print_json(get_tapology_for("Jon Jones"))
	print_json(get_tapology_for("Conor Mcgregor"))
}

func TestTapologyParseFragment(t *testing.T) {
    print_json(parse_tapology_results(search_result))
}
