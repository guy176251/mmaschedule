package main

import (
	"context"
	"fmt"
	"mmaschedule-go/event"
	"time"
)

type Scraper func(ClientScraper) (*[]*event.Event, error)

func ScrapeEvents(q *event.Queries, client ClientScraper, tapology bool) {
	events := []*event.Event{}
	scrapers := []Scraper{
		ScrapeONE,
		ScrapeUFC,
	}

	for _, scraper := range scrapers {
		e, err := scraper(client)
		if err != nil {
			fmt.Println("Error scraping events:", err)
		} else {
			events = append(events, *e...)
		}
	}

	if tapology {
		UpdateTapology(q, client, &events)
	}

	if len(events) > 0 {
		err := q.UpsertEvents(context.Background(), events)
		if err != nil {
			fmt.Println("Error updating events in database:", err)
		}
	}
}

func UpdateTapology(q *event.Queries, client ClientScraper, events *[]*event.Event) {
	err := SetTapologyCSRF(client)
	if err != nil {
		fmt.Println("Error settings tapology CSRF:", err)
	}

	for _, e := range *events {
		fights := e.UnmarshalFights()
		for _, f := range fights {
			fighters := []*event.Fighter{
				f.FighterA,
				f.FighterB,
			}
			for _, ff := range fighters {
				fmt.Println("Getting tapology link for", ff.Name)
				tapology, err := q.GetTapology(context.Background(), ff.Name)
				if err != nil {
					fmt.Println("Error getting tapology from database:", err)
					link, err := GetTapologyLink(client, ff.Name)
					if err != nil {
						fmt.Println("Error getting tapology from site:", err)
					} else if len(link) > 0 {
						err := q.CreateTapology(context.Background(), event.CreateTapologyParams{Name: ff.Name, Url: link})
						if err != nil {
							fmt.Println("Error creating tapology in database:", err)
						}
						ff.Link = link
					}
					time.Sleep(time.Duration(5_000_000_000))
				} else {
					ff.Link = tapology.Url
				}
			}
		}
		e.MarshalFights(fights)
	}
}

func UpdateAllTapology(q *event.Queries, client ClientScraper) error {
	events_, err := q.ListEvents(context.Background())
	if err != nil {
		return err
	}

	events := []*event.Event{}
	for _, e := range events_ {
		events = append(events, &e)
	}

	UpdateTapology(q, client, &events)

	err = q.UpsertEvents(context.Background(), events)
	if err != nil {
		return err
	}

	return nil
}
