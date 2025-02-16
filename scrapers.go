package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"mmaschedule-go/event"

	"github.com/PuerkitoBio/goquery"
)

type ClientScraper interface {
	Get(url string, options ...RequestOption) (*goquery.Selection, error)
	AddHeader(key string, value string)
}

func ScrapeONE(client ClientScraper) (*[]*event.Event, error) {
	events := []*event.Event{}
	page, err := client.Get("https://www.onefc.com/events/")

	if err != nil {
		return nil, err
	}

	page.Find("#upcoming-events-section .is-event").Each(func(i int, s *goquery.Selection) {
		_, exists := s.Find("a.smart-link").Attr("href")

		if !exists {
			return
		}

		url, exists := s.Find("a.is-image-zoom[href]").Attr("href")

		if !exists {
			return
		}

		e := event.Event{
			Url:          url,
			Organization: "ONE",
			Slug: strings.Trim(
				strings.Replace(url, "https://www.onefc.com/events/", "", 1),
				"/",
			),
		}

		err := ScrapeONEEvent(client, &e)

		if err != nil {
			log.Println(err)
			return
		}

		if e.Fights == "null" {
			return
		}

		events = append(events, &e)
	})

	return &events, nil
}

type ONEEventLocation struct {
	Name string `json:"name"`
}

type ONEEventData struct {
	Name      string           `json:"name"`
	StartDate string           `json:"startDate"`
	Image     []string         `json:"image"`
	Location  ONEEventLocation `json:"location"`
}

func ScrapeONEEvent(client ClientScraper, e *event.Event) error {
	page, err := client.Get(e.Url)

	if err != nil {
		return err
	}

	data := ONEEventData{}
	_ = json.Unmarshal([]byte(page.Find("#site-main script[type=\"application/ld+json\"]").Text()), &data)

	e.Name = data.Name
	e.Location = data.Location.Name

	if len(data.Image) > 2 {
		e.Image = data.Image[2]
	}

	date, err := time.Parse(time.RFC3339, data.StartDate)
	if err == nil {
		e.Date = date.Unix()
	}

	var fights []*event.Fight
	page.Find(".event-matchup").Each(func(i int, s *goquery.Selection) {
		fights = append(fights, ScrapeONEFight(s))
	})
	e.MarshalFights(fights)

	return nil
}

func ScrapeONEFight(s *goquery.Selection) *event.Fight {
	fight := event.Fight{}

	fight.Weight = TextContent(s, ".title")
	fight.FighterA = ScrapeONEFighter(s, true)
	fight.FighterB = ScrapeONEFighter(s, false)

	return &fight
}

func ScrapeONEFighter(s *goquery.Selection, is_a bool) *event.Fighter {
	fighter := event.Fighter{}
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

func ScrapeUFC(client ClientScraper) (*[]*event.Event, error) {
	events := []*event.Event{}
	page, err := client.Get("https://www.ufc.com/events")

	if err != nil {
		return nil, err
	}

	page.Find("#events-list-upcoming .c-card-event--result__headline a").Each(func(i int, s *goquery.Selection) {
		slug, exists := s.Attr("href")
		if exists {
			e := event.Event{
				Url:          "https://www.ufc.com" + slug,
				Organization: "UFC",
				Slug:         strings.Replace(slug, "/event/", "", 1),
			}
			err := ScrapeUFCEvent(client, &e)
			if err != nil {
				log.Println(err)
			} else if e.Fights != "null" {
				events = append(events, &e)
			}
		}
	})

	return &events, nil
}

func ScrapeUFCEvent(client ClientScraper, e *event.Event) error {
	page, err := client.Get(e.Url)

	if err != nil {
		return err
	}

	e.Name = strings.Join(
		page.Find(".c-hero__headline-prefix, .c-hero__headline").Map(func(i int, s *goquery.Selection) string {
			return CleanupWhitespace(s.Text())
		}),
		": ",
	)
	e.Location = TextContent(page, ".field--name-venue")

	image, exists := page.Find(".c-hero__image img").First().Attr("src")
	if exists {
		e.Image = image
	}

	datestr, exists := page.Find(`

        #early-prelims .c-event-fight-card-broadcaster__time,
        #prelims-card .c-event-fight-card-broadcaster__time,
        .c-hero__headline-suffix

    `).Last().Attr("data-timestamp")

	if exists {
		date, err := strconv.Atoi(datestr)
		if err == nil {
			e.Date = int64(date)
		}
	}

	var fights []*event.Fight
	page.Find(".l-listing__item").Each(func(i int, s *goquery.Selection) {
		fights = append(fights, ScrapeUFCFight(s))
	})
	e.MarshalFights(fights)

	return nil
}

func ScrapeUFCFight(s *goquery.Selection) *event.Fight {
	fight := event.Fight{}

	fight.Weight = strings.Replace(TextContent(s, ".c-listing-fight__class-text"), " Bout", "", -1)
	fight.FighterA = ScrapeUFCFighter(s, "red")
	fight.FighterB = ScrapeUFCFighter(s, "blue")

	return &fight
}

func ScrapeUFCFighter(s *goquery.Selection, color string) *event.Fighter {
	fighter := event.Fighter{}

	fighter.Name = TextContent(s, ".c-listing-fight__corner-name--"+color)
	fighter.Country = TextContent(s, ".c-listing-fight__country--"+color+" .c-listing-fight__country-text")

	image, exists := s.Find(".c-listing-fight__corner-image--" + color + " img").First().Attr("src")
	if exists {
		fighter.Image = image
	}

	return &fighter
}

func SetTapologyCSRF(client ClientScraper) error {
	selection, err := client.Get("https://www.tapology.com")
	if err != nil {
		return err
	}
	token, exists := selection.Find("meta[name=\"csrf-token\"]").First().Attr("content")
	if !exists {
		return fmt.Errorf("Could not parse CSRF token.")
	}
	client.AddHeader("X-CSRF-Token", token)
	return nil
}

var nickname *regexp.Regexp = regexp.MustCompile(`"(\w| )+"`)

type TapologyResult struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func GetTapologyLink(client ClientScraper, name string) (string, error) {
	selection, err := client.Get("https://www.tapology.com/search/nav", func(r *http.Request) {
		query := url.Values{}
		query.Set("ajax", "true")
		query.Set("model", "fighters")
		query.Set("term", name)
		r.URL.RawQuery = query.Encode()
	})

	if err != nil {
		return "", err
	}

	results := []TapologyResult{}

	selection.Find("span.star a[href]").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		name := CleanupWhitespace(nickname.ReplaceAllString(s.Text(), ""))
		results = append(results, TapologyResult{Name: name, Url: url})
	})

	if len(results) == 0 {
		return "", nil
	}

	slices.SortStableFunc(results, func(a, b TapologyResult) int {
		score_a := HammingScore(name, a.Name)
		score_b := HammingScore(name, b.Name)

		if score_a > score_b {
			return -1
		} else if score_a < score_b {
			return 1
		} else {
			return 0
		}
	})

	return results[0].Url, nil
}
