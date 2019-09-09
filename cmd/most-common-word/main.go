package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"sort"
	"strings"

	parser "github.com/marcsantiago/app-review-parser"
	prose "gopkg.in/jdkato/prose.v2"
)

// example run
// go run main.go -appid 639881495 > ~/Desktop/imgur_reviews.tsv
func main() {
	var appID string
	flag.StringVar(&appID, "appid", "", "the apple app store numeric id")
	flag.Parse()
	if len(appID) == 0 {
		flag.PrintDefaults()
		log.Fatal(1)
	}

	// tabs are used over commas because people may have commas in sentences
	var buf bytes.Buffer
	iosData := parser.FetchIOSAppReviews(appID)
	for _, data := range iosData {
		for _, entry := range data.Feed.Entry {

			if len(entry.QuickRow()) > 0 {
				buf.WriteString(strings.Join(entry.QuickRow(), " ") + "\n")
			}
		}
	}

	doc, err := prose.NewDocument(buf.String())
	if err != nil {
		log.Fatal(err)
	}

	// Iterate over the doc's tokens:
	count := make(map[string]int)
	for _, tok := range doc.Tokens() {
		if validTags(tok.Tag) {
			count[strings.ToLower(tok.Text)]++
		}
	}

	counts := make([]countSorted, 0, len(count))
	for k, v := range count {
		counts = append(counts, countSorted{Text: k, Count: v})
	}
	sort.Sort(countsSorted(counts))

	max := len(counts)
	if max > 100 {
		max = 100
	}

	for i := 0; i <= max; i++ {
		fmt.Println(counts[i].Text, counts[i].Count)
	}

}

type countSorted struct {
	Text  string
	Count int
}

type countsSorted []countSorted

func (c countsSorted) Len() int           { return len(c) }
func (c countsSorted) Less(i, j int) bool { return c[i].Count > c[j].Count }
func (c countsSorted) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

func validTags(tag string) bool {
	switch tag {
	case "JJ", "JJR", "JJS", "RB", "RBR", "RBS", "RP", "VB", "VBD", "VBG", "VBN", "VBZ", "WRB":
		return true
	}
	return false
}
