package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func UFCCollector() *colly.Collector {
	return colly.NewCollector(
		colly.UserAgent(USER_AGENT),
		colly.AllowedDomains("ufc.com", "www.ufc.com"),
	)
}

func ScrapeUFC(callback EventCallback) {
	ufclist := UFCCollector()
	ufcdetail := UFCCollector()

	ufclist.OnResponse(func(r *colly.Response) {
		for _, url := range ParseUFCEventList(r.Body) {
			log.Printf("Visiting %s\n", url)
			err := ufcdetail.Visit(url)
			if err != nil {
				log.Println(err)
			}
		}
	})

	ufcdetail.OnResponse(func(r *colly.Response) {
		callback(ParseUFCEvent(r.Body))
	})

	err := ufclist.Visit("https://www.ufc.com/events")
	if err != nil {
		log.Println(err)
	}
}

func ParseUFCEventList(b []byte) []string {
	urls := []string{}
	doc := DocumentFromBytes(b)

	doc.Find(".c-card-event--result__headline a").Each(func(i int, s *goquery.Selection) {
		slug, ok := s.Attr("href")
		if ok {
			urls = append(urls, "https://ufc.com"+slug)
		}
	})

	return urls
}

func ParseUFCEvent(b []byte) *Event {
	event := Event{Organization: "UFC"}
	doc := DocumentFromBytes(b)

	parts := doc.Find(".c-hero__headline-prefix, .c-hero__headline").Map(func(i int, s *goquery.Selection) string {
		return CleanupString(s.Text())
	})
	event.Name = strings.Join(parts, ": ")
	event.Location = TextContent(doc, ".field--name-venue")

	image, ok := doc.Find(".c-hero__image img").First().Attr("src")
	if ok {
		event.Image = image
	}

	datestr, ok := doc.Find(`

        #early-prelims .c-event-fight-card-broadcaster__time,
        #prelims-card .c-event-fight-card-broadcaster__time,
        .c-hero__headline-suffix

    `).Last().Attr("data-timestamp")
	if ok {
		date, err := strconv.Atoi(datestr)
		if err == nil {
			event.Date = date
		}
	}

	doc.Find(".l-listing__item").Each(func(i int, s *goquery.Selection) {
		event.Fights = append(event.Fights, *ParseUFCFight(s))
	})

	return &event
}

func ParseUFCFight(doc *goquery.Selection) *Fight {
	fight := Fight{}

	fight.Weight = strings.Replace(TextContent(doc, ".c-listing-fight__class-text"), " Bout", "", -1)
	fight.FighterA = ParseUFCFighter(doc, "red")
	fight.FighterB = ParseUFCFighter(doc, "blue")

	return &fight
}

func ParseUFCFighter(doc *goquery.Selection, color string) *Fighter {
	fighter := Fighter{}

	fighter.Name = TextContent(doc, ".c-listing-fight__corner-name--"+color)
	fighter.Country = TextContent(doc, ".c-listing-fight__country--"+color+" .c-listing-fight__country-text")
	image, ok := doc.Find(".c-listing-fight__corner-image--" + color + " img").First().Attr("src")
	if ok {
		fighter.Image = image
	}

	return &fighter
}
