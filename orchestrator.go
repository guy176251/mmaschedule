package main

import (
	"fmt"
	"mmaschedule-go/event"
)

type Scraper func(ClientScraper) (*[]*event.Event, error)

func ScrapeEvents() {
	client := NewScraperClient()
	events := []*event.Event{}
	scrapers := []Scraper{
		ScrapeONE,
		ScrapeUFC,
	}

	for _, scraper := range scrapers {
		e, err := scraper(&client)
		if err != nil {
			fmt.Println(err)
		} else {
			events = append(events, *e...)
		}
	}
}
