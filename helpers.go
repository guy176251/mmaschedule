package main

import (
	"bytes"
	"encoding/json"
	"hash/fnv"
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
