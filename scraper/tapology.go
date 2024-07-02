package scraper

import (
	"log"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

const tapology_url string = "https://www.tapology.com"

type TapologyResult struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func tapology_collector() *colly.Collector {
	return colly.NewCollector(
		colly.UserAgent(USER_AGENT),
		colly.AllowedDomains("tapology.com", "www.tapology.com"),
	)
}

func TapologyGetter() func(n string) []TapologyResult {
	csrf_token := ""
	name := ""
	results := []TapologyResult{}

	index := tapology_collector()
	index.OnHTML("meta[name=\"csrf-token\"]", func(h *colly.HTMLElement) {
		csrf_token = h.Attr("content")
		time.Sleep(5 * time.Second)
	})

	search := tapology_collector()
	search.AllowURLRevisit = true
	search.OnRequest(func(r *colly.Request) {
		r.Headers.Add("X-CSRF-Token", csrf_token)

		query := url.Values{}
		query.Set("ajax", "true")
		query.Set("model", "fighters")
		query.Set("term", name)
		r.URL.RawQuery = query.Encode()

		log.Printf("Making tapology request to %s", r.URL)
	})
	search.OnResponse(func(r *colly.Response) {
		results = parse_tapology_results(r.Body)
		time.Sleep(5 * time.Second)
	})

	if err := index.Visit(tapology_url); err != nil {
		log.Println(err)
	}

	return func(n string) []TapologyResult {
		if csrf_token == "" {
			return nil
		}

		name = n
		if err := search.Visit(tapology_url + "/search/nav"); err != nil {
			log.Println(err)
			return nil
		}

		return results
	}
}

func parse_tapology_results(b []byte) []TapologyResult {
	results := []TapologyResult{}

	DocumentFromBytes(b).Find("span.star a[href]").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		name := s.Text()
		results = append(results, TapologyResult{Name: name, Url: url})
	})

	return results
}
