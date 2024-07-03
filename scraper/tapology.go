package scraper

import (
	"log"
	"net/url"
	"regexp"
	"slices"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

const tapology_url string = "https://www.tapology.com"

type TapologyResult struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type TapologyGetter func(n string) string

func tapology_getter() TapologyGetter {
	csrf_token := ""
	name := ""
	result := ""

	index := tapology_collector()
	index.OnHTML("meta[name=\"csrf-token\"]", func(h *colly.HTMLElement) {
		csrf_token = h.Attr("content")
		time.Sleep(5 * time.Second)
	})

	if err := index.Visit(tapology_url); err != nil {
		log.Println(err)
	}

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
		result = parse_tapology_results(name, r.Body)
		time.Sleep(5 * time.Second)
	})

	return func(n string) string {
		if csrf_token == "" {
			return ""
		}

		name = n
		if err := search.Visit(tapology_url + "/search/nav"); err != nil {
			log.Println(err)
			return ""
		}

		return result
	}
}

var nickname *regexp.Regexp = regexp.MustCompile(`"(\w| )+"`)

func parse_tapology_results(name string, b []byte) string {
	results := []TapologyResult{}

	document_from_bytes(b).Find("span.star a[href]").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		name := cleanup_whitespace(nickname.ReplaceAllString(s.Text(), ""))
		results = append(results, TapologyResult{Name: name, Url: url})
	})

	slices.SortStableFunc(results, func(a, b TapologyResult) int {
		score_a := hamming_score(name, a.Name)
		score_b := hamming_score(name, b.Name)

		if score_a > score_b {
			return -1
		} else if score_a < score_b {
			return 1
		} else {
			return 0
		}
	})

	if len(results) < 1 {
		return ""
	}

	return tapology_url + results[0].Url
}

func tapology_collector() *colly.Collector {
	return colly.NewCollector(
		colly.UserAgent(USER_AGENT),
		colly.AllowedDomains("tapology.com", "www.tapology.com"),
	)
}
