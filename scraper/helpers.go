package scraper

import (
	"bytes"
	"encoding/json"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
)

func text_content(doc *goquery.Selection, selector string) string {
	return cleanup_whitespace(doc.Find(selector).First().Text())
}

var whitespace *regexp.Regexp = regexp.MustCompile(`\s+`)

func cleanup_whitespace(s string) string {
	return strings.TrimSpace(whitespace.ReplaceAllString(s, " "))
}

func document_from_bytes(b []byte) *goquery.Selection {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(b))
	return doc.Find("html").First()
}

func print_json(v any) {
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Println(err)
	} else {
		log.Println(string(out))
	}
}

var hamming *metrics.Hamming = metrics.NewHamming()

func hamming_score(a, b string) float64 {
    return strutil.Similarity(a, b, hamming)
}
