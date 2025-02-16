package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

//go:embed testdata/ufc-event-list.html
var UFCEventList []byte

//go:embed testdata/ufc-312.html
var UFC312 []byte

//go:embed testdata/ufc-fn-feb-15.html
var UFCFNFeb15 []byte

//go:embed testdata/ufc-fn-feb-22.html
var UFCFNFeb22 []byte

//go:embed testdata/ufc-fn-mar-1.html
var UFCFNMar1 []byte

//go:embed testdata/ufc-313.html
var UFC313 []byte

//go:embed testdata/ufc-fn-mar-15.html
var UFCFNMar15 []byte

//go:embed testdata/ufc-fn-mar-22.html
var UFCFNMar22 []byte

//go:embed testdata/ufc-fn-mar-29.html
var UFCFNMar29 []byte

//go:embed testdata/one-event-list.html
var ONEEventList []byte

//go:embed testdata/one-friday-fights-97.html
var ONEFridayFights97 []byte

//go:embed testdata/one-171.html
var ONE171 []byte

//go:embed testdata/one-fn-29.html
var ONEFN29 []byte

//go:embed testdata/one-172.html
var ONE172 []byte

var HTMLContent map[string][]byte = map[string][]byte{
	"https://www.ufc.com/events":                                 UFCEventList,
	"https://www.ufc.com/event/ufc-312":                          UFC312,
	"https://www.ufc.com/event/ufc-fight-night-february-15-2025": UFCFNFeb15,
	"https://www.ufc.com/event/ufc-fight-night-february-22-2025": UFCFNFeb22,
	"https://www.ufc.com/event/ufc-fight-night-march-01-2025":    UFCFNMar1,
	"https://www.ufc.com/event/ufc-313":                          UFC313,
	"https://www.ufc.com/event/ufc-fight-night-march-15-2025":    UFCFNMar15,
	"https://www.ufc.com/event/ufc-fight-night-march-22-2025":    UFCFNMar22,
	"https://www.ufc.com/event/ufc-fight-night-march-29-2025":    UFCFNMar29,
	"https://www.onefc.com/events/":                              ONEEventList,
	"https://www.onefc.com/events/one-friday-fights-97/":         ONEFridayFights97,
	"https://www.onefc.com/events/one171/":                       ONE171,
	"https://www.onefc.com/events/onefightnight29/":              ONEFN29,
	"https://www.onefc.com/events/one172/":                       ONE172,
}

type DummyClient struct{}

func (c *DummyClient) Get(url string, options ...RequestOption) (*goquery.Selection, error) {
	htmlstring, exists := HTMLContent[url]
	if !exists {
		return nil, fmt.Errorf("Invalid URL: %s", url)
	}
	document, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlstring))
	if err != nil {
		return nil, err
	}
	return document.Find("html").First(), nil
}

func (c *DummyClient) AddHeader(key string, value string) {}

func TestUFC(t *testing.T) {
	client := DummyClient{}
	events, err := ScrapeUFC(&client)
	if err != nil {
		t.Errorf("Scraping UFC failed: %s", err)
	}
	if len(*events) == 0 {
		t.Fatal("Scraping UFC failed: Events are empty")
	}
	if len(*events) < 7 {
		t.Fatal("Scraping UFC failed: Did not get all events")
	}
	//PrintJson(events)
}

func TestONE(t *testing.T) {
	client := DummyClient{}
	events, err := ScrapeONE(&client)
	if err != nil {
		t.Errorf("Scraping ONE failed: %s", err)
	}
	PrintJson(events)
}

func TestJSONString(t *testing.T) {
	type User struct {
		Name string
	}

	type Thing struct {
		Users []User `json:"user,string"`
		Num   int    `json:"num,string"`
	}

	u := Thing{
		Users: []User{
			{Name: "Bob"},
			{Name: "John"},
		},
	}

	PrintJson(u)
}
