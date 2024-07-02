package scraper

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func ufc_collector() *colly.Collector {
	return colly.NewCollector(
		colly.UserAgent(USER_AGENT),
		colly.AllowedDomains("ufc.com", "www.ufc.com"),
	)
}

func ScrapeUFC(callback EventCallback) {
	list := ufc_collector()
	detail := ufc_collector()
	url := ""

	list.OnResponse(func(r *colly.Response) {
		for _, u := range parse_ufc_event_list(r.Body) {
			url = u
			log.Printf("Visiting %s\n", u)
			if err := detail.Visit(u); err != nil {
				log.Println(err)
			}
		}
	})

	detail.OnResponse(func(r *colly.Response) {
		callback(parse_ufc_event(url, r.Body))
	})

	if err := list.Visit("https://www.ufc.com/events"); err != nil {
		log.Println(err)
	}
}

func parse_ufc_event_list(b []byte) []string {
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

func parse_ufc_event(url string, b []byte) *Event {
	event := Event{Organization: "UFC", Url: url}
	doc := DocumentFromBytes(b)

	parts := doc.Find(".c-hero__headline-prefix, .c-hero__headline").Map(func(i int, s *goquery.Selection) string {
		return CleanupString(s.Text())
	})
	event.Name = strings.Join(parts, ": ")
	event.Location = TextContent(doc, ".field--name-venue")
	event.Slug = parse_ufc_slug(url)

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
		event.Fights = append(event.Fights, *parse_ufc_fight(s))
	})

	return &event
}

func parse_ufc_fight(s *goquery.Selection) *Fight {
	fight := Fight{}

	fight.Weight = strings.Replace(TextContent(s, ".c-listing-fight__class-text"), " Bout", "", -1)
	fight.FighterA = parse_ufc_fighter(s, "red")
	fight.FighterB = parse_ufc_fighter(s, "blue")

	return &fight
}

func parse_ufc_fighter(s *goquery.Selection, color string) *Fighter {
	fighter := Fighter{}

	fighter.Name = TextContent(s, ".c-listing-fight__corner-name--"+color)
	fighter.Country = TextContent(s, ".c-listing-fight__country--"+color+" .c-listing-fight__country-text")
	image, ok := s.Find(".c-listing-fight__corner-image--" + color + " img").First().Attr("src")
	if ok {
		fighter.Image = image
	}

	return &fighter
}

var ufc_event_url *regexp.Regexp = regexp.MustCompile(`https://ufc.com/event/(?P<Slug>(\w|-)+)`)

func parse_ufc_slug(url string) string {
	matches := ufc_event_url.FindStringSubmatch(url)
	index := ufc_event_url.SubexpIndex("Slug")
	if len(matches) > index {
		return matches[index]
	}
	return ""
}
