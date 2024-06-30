package main

type Event struct {
	Name         string  `json:"name"`
	Location     string  `json:"location"`
	Organization string  `json:"organization"`
	Image        string  `json:"image"`
	Date         int     `json:"date"`
	Fights       []Fight `json:"fights"`
}

type Fight struct {
	Weight   string   `json:"weight"`
	FighterA *Fighter `json:"fighter_a"`
	FighterB *Fighter `json:"fighter_b"`
}

type Fighter struct {
	Name    string `json:"name"`
	Image   string `json:"image"`
	Country string `json:"country"`
}

type EventCallback func(e *Event)
type EventScraper func(c EventCallback)

const USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36"

var Scrapers []EventScraper = []EventScraper{
	ScrapeUFC,
}
