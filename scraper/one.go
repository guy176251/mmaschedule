package scraper

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type event_location struct {
	Name string `json:"name"`
}

type event_data struct {
	Name      string         `json:"name"`
	StartDate string         `json:"startDate"`
	Image     []string       `json:"image"`
	Location  event_location `json:"location"`
}

func scrape_one(callback EventCallback) {
	list := one_collector()
	detail := one_collector()
	url := ""

	list.OnResponse(func(r *colly.Response) {
		for _, u := range parse_one_event_list(r.Body) {
			url = u
			log.Printf("Visiting %s\n", u)
			if err := detail.Visit(u); err != nil {
				log.Println(err)
			}
		}
	})

	detail.OnResponse(func(r *colly.Response) {
		callback(parse_one_event(url, r.Body))
	})

	if err := list.Visit("https://www.onefc.com/events/"); err != nil {
		log.Println(err)
	}
}

func parse_one_event_list(b []byte) []string {
	return document_from_bytes(b).Find(".is-event a.is-image-zoom[href]").Map(func(i int, s *goquery.Selection) string {
		url, _ := s.Attr("href")
		return url
	})
}

func parse_one_event(url string, b []byte) *Event {
	event := Event{Organization: "ONE", Url: url}
	doc := document_from_bytes(b)

	data := event_data{}
	_ = json.Unmarshal([]byte(doc.Find("#site-main script[type=\"application/ld+json\"]").Text()), &data)

	event.Name = data.Name
	event.Location = data.Location.Name

	if len(data.Image) > 2 {
		event.Image = data.Image[2]
	}

	date, err := time.Parse(time.RFC3339, data.StartDate)
	if err == nil {
		event.Date = int(date.Unix())
	}

	doc.Find(".event-matchup").Each(func(i int, s *goquery.Selection) {
		event.Fights = append(event.Fights, *parse_one_fight(s))
	})

	return &event
}

func parse_one_fight(s *goquery.Selection) *Fight {
	fight := Fight{}

	fight.Weight = text_content(s, ".title")
	fight.FighterA = parse_one_fighter(s, true)
	fight.FighterB = parse_one_fighter(s, false)

	return &fight
}

func parse_one_fighter(s *goquery.Selection, is_a bool) *Fighter {
	fighter := Fighter{}
	nth_child := "1"
	face := "1"

	if is_a {
		nth_child = "3"
		face = "2"
	}

	fighter.Name = text_content(s, fmt.Sprintf("tr.vs :nth-child(%s) a", nth_child))
	fighter.Country = text_content(s, fmt.Sprintf("tr.vs + tr :nth-child(%s)", nth_child))

	image, _ := s.Find(fmt.Sprintf(".face%s img[src]", face)).Attr("src")
	fighter.Image = image

	return &fighter
}

func one_collector() *colly.Collector {
	return colly.NewCollector(
		colly.UserAgent(USER_AGENT),
		colly.AllowedDomains("onefc.com", "www.onefc.com"),
	)
}
