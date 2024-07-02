package scraper

import (
	"bytes"
	"encoding/json"
	"hash/fnv"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func TextContent(doc *goquery.Selection, selector string) string {
	return CleanupString(doc.Find(selector).First().Text())
}

var whitespace *regexp.Regexp = regexp.MustCompile(`\s+`)

func CleanupString(s string) string {
	return strings.TrimSpace(whitespace.ReplaceAllString(s, " "))
}

func DocumentFromBytes(b []byte) *goquery.Selection {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(b))
	return doc.Find("html").First()
}

func HashJSON(v any) uint32 {
	out, err := json.Marshal(v)
	if err != nil {
		return 0
	}
	h := fnv.New32a()
	h.Write(out)
	return h.Sum32()
}

func PrintJSON(v any) {
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Println(err)
	} else {
		log.Println(string(out))
	}
}
