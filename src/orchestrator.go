package main

import (
	"context"
	"log/slog"
	"mmaschedule-go/event"
	"time"
)

type Scraper func(ClientScraper) (*[]*event.Event, error)

func ScrapeEventsLoop(db *event.Queries, client ClientScraper, tapology bool) {
	for {
		time.Sleep(time.Duration(1) * time.Hour)
		slog.Debug("Running hourly scraper")
		ScrapeEvents(db, client, tapology)
	}
}

func ScrapeEvents(db *event.Queries, client ClientScraper, tapology bool) {
	allEvents := []*event.Event{}
	scrapers := []Scraper{
		ScrapeONE,
		ScrapeUFC,
	}

	for _, scraper := range scrapers {
		events, err := scraper(client)

		if err != nil {
			slog.Error("Error scraping events", "error", err)
			continue
		}

		for _, e := range *events {
			if !e.HasEmptyFights() {
				allEvents = append(allEvents, e)
			}
		}
	}

	if tapology {
		UpdateTapology(db, client, &allEvents)
	}

	if len(allEvents) > 0 {
		err := db.UpsertEvents(context.Background(), allEvents)
		if err != nil {
			slog.Error("Failed updating events in database", "error", err)
		}
	}
}

func UpdateTapology(db *event.Queries, client ClientScraper, events *[]*event.Event) {
	err := SetTapologyCSRF(client)
	if err != nil {
		slog.Error("Failed settings tapology CSRF", "error", err)
	}

	for _, e := range *events {
		fights := e.UnmarshalFights()

		for _, f := range fights {
			fighters := []*event.Fighter{
				f.FighterA,
				f.FighterB,
			}

			for _, ff := range fighters {
				slog.Debug("Getting tapology link from database", "name", ff.Name)
				tapology, err := db.GetTapology(context.Background(), ff.Name)

				if err == nil {
					ff.Link = tapology.Url
					continue
				}

				slog.Error("Failed getting tapology from database", "name", ff.Name, "error", err)
				link, err := GetTapologyLink(client, ff.Name)
				time.Sleep(5 * time.Second)

				if err != nil {
					slog.Error("Failed getting tapology from site", "name", ff.Name, "error", err)
				}

				ff.Link = link
				slog.Debug("Got new tapology link", "name", ff.Name, "link", link)
				err = db.CreateTapology(context.Background(), event.CreateTapologyParams{Name: ff.Name, Url: link})

				if err != nil {
					slog.Error("Failed creating tapology in database", "name", ff.Name, "error", err)
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
