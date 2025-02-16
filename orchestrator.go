package main

import (
	"context"
	"fmt"
	"mmaschedule-go/event"
)

type Scraper func(ClientScraper) (*[]*event.Event, error)

func ScrapeEvents(q *event.Queries, client ClientScraper) {
	events := []*event.Event{}
	scrapers := []Scraper{
		ScrapeONE,
		ScrapeUFC,
	}

	for _, scraper := range scrapers {
		e, err := scraper(client)
		if err != nil {
			fmt.Println("Error scraping events: ", err)
		} else {
			events = append(events, *e...)
		}
	}

	if len(events) > 0 {
		err := q.UpsertEvents(context.Background(), events)
		if err != nil {
			fmt.Println("Error updating events in database: ", err)
		}
	}
}
