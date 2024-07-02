package scraper

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type ONEEventLocation struct {
	Name string `json:"name"`
}

type ONEEventData struct {
	Name      string           `json:"name"`
	StartDate string           `json:"startDate"`
	Image     []string         `json:"image"`
	Location  ONEEventLocation `json:"location"`
}

func ONECollector() *colly.Collector {
	return colly.NewCollector(
		colly.UserAgent(USER_AGENT),
		colly.AllowedDomains("onefc.com", "www.onefc.com"),
	)
}

func ScrapeONE(callback EventCallback) {
	list := ONECollector()
	detail := ONECollector()
	url := ""

	list.OnResponse(func(r *colly.Response) {
		for _, u := range ParseONEEventList(r.Body) {
			url = u
			log.Printf("Visiting %s\n", u)
			if err := detail.Visit(u); err != nil {
				log.Println(err)
			}
		}
	})

	detail.OnResponse(func(r *colly.Response) {
		callback(ParseONEEvent(url, r.Body))
	})

	if err := list.Visit("https://www.onefc.com/events/"); err != nil {
		log.Println(err)
	}
}

func ParseONEEventList(b []byte) []string {
	return DocumentFromBytes(b).Find(".is-event a.is-image-zoom[href]").Map(func(i int, s *goquery.Selection) string {
		url, _ := s.Attr("href")
		return url
	})
}

func ParseONEEvent(url string, b []byte) *Event {
	event := Event{Organization: "ONE", Url: url}
	doc := DocumentFromBytes(b)

	data := ONEEventData{}
	_ = json.Unmarshal([]byte(doc.Find("#site-main script[type=\"application/ld+json\"]").Text()), &data)

	event.Name = data.Name
	event.Image = data.Image[2]
	event.Location = data.Location.Name

	date, err := time.Parse(time.RFC3339, data.StartDate)
	if err == nil {
		event.Date = int(date.Unix())
	}

	doc.Find(".event-matchup").Each(func(i int, s *goquery.Selection) {
		event.Fights = append(event.Fights, *ParseONEFight(s))
	})

	return &event
}

func ParseONEFight(s *goquery.Selection) *Fight {
	fight := Fight{}

	fight.Weight = TextContent(s, ".title")
	fight.FighterA = ParseONEFighter(s, true)
	fight.FighterB = ParseONEFighter(s, false)

	return &fight
}

func ParseONEFighter(s *goquery.Selection, is_a bool) *Fighter {
	fighter := Fighter{}
	nth_child := "1"
	face := "1"

	if is_a {
		nth_child = "3"
		face = "2"
	}

	fighter.Name = TextContent(s, fmt.Sprintf("tr.vs :nth-child(%s) a", nth_child))
	fighter.Country = TextContent(s, fmt.Sprintf("tr.vs + tr :nth-child(%s)", nth_child))

	image, _ := s.Find(fmt.Sprintf(".face%s img[src]", face)).Attr("src")
	fighter.Image = image

	return &fighter
}
