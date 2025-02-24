package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
)

type Values = url.Values

func TextContent(doc *goquery.Selection, selector string) string {
	return CleanupWhitespace(doc.Find(selector).First().Text())
}

var whitespace *regexp.Regexp = regexp.MustCompile(`\s+`)

func CleanupWhitespace(s string) string {
	return strings.TrimSpace(whitespace.ReplaceAllString(s, " "))
}

func PrintJson(v any) {
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Println(err)
	} else {
		log.Println(string(out))
	}
}

var hamming *metrics.Hamming = metrics.NewHamming()

func HammingScore(a, b string) float64 {
	return strutil.Similarity(a, b, hamming)
}

type ScraperClient struct {
	headers map[string]string
	client  http.Client
}

type RequestOption func(*http.Request)

func (c *ScraperClient) Get(url string, options ...RequestOption) (*goquery.Selection, error) {
	resp, err := c.GetResponse(url, options...)
	if err != nil {
		return nil, err
	}
	selection, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	return selection.Find("html").First(), nil
}

func (c *ScraperClient) GetResponse(url string, options ...RequestOption) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}
	for _, option := range options {
		if option != nil {
			option(req)
		}
	}
	return c.client.Do(req)
}

func (c *ScraperClient) AddHeader(key string, value string) {
	c.headers[key] = value
}
func (c *ScraperClient) HasHeader(key string) bool {
	_, exists := c.headers[key]
	return exists
}

func NewScraperClient() *ScraperClient {
	client := ScraperClient{
		headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
		},
		client: http.Client{},
	}
	return &client
}
