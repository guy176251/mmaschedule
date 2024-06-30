package main

import (
	_ "embed"
	"encoding/json"
	"log"
	"testing"
)

//go:embed testdata/ufc-297-good.html
var UFCEventHTML []byte

//go:embed testdata/ufc-events.html
var UFCEventListHTML []byte

func TestUFCEventList(t *testing.T) {
	urls := ParseUFCEventList(UFCEventListHTML)
    PrintJSON(urls)
}

func TestUFCEvent(t *testing.T) {
	event := ParseUFCEvent(UFCEventHTML)
    PrintJSON(event)
}

func testScrapeUFC(t *testing.T) {
    ScrapeUFC(func(e *Event) {
        PrintJSON(e)
    })
}

func PrintJSON(v any) {
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Println(err)
	} else {
        log.Println(string(out))
    }
}
