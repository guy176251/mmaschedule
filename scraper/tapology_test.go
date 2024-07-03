package scraper

import (
	_ "embed"
	"testing"
)

//go:embed testdata/tapology-search.html
var search_result []byte

func testTapologyGetter(t *testing.T) {
	get_tapology_for := tapology_getter()
	names := []string{
		"Dustin Poirier",
		"Jon Jones",
		"Conor Mcgregor",
        "Justin Gaethje",
        "Benoit Saint Denis",
	}

	for _, name := range names {
		print_json(get_tapology_for(name))
	}
}

func TestParseTapologyResults(t *testing.T) {
	print_json(parse_tapology_results("Jon Jones", search_result))
}
