package main

import (
	"mmaschedule-go/events"
	"mmaschedule-go/scraper"
)

func stuff(db events.DbEvent) {
	e := scraper.Event{}
	db.ReadData(&e)
}
