package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	parser "github.com/marcsantiago/app-review-parser"
)

var headersOnce sync.Once

func printHeadersOnce(e parser.Entry) {
	headersOnce.Do(func() {
		fmt.Println(strings.Join(e.QuickHeaders(), "\t"))
	})
}

// example run
// go run main.go -appid 639881495 > ~/Desktop/imgur_reviews.tsv
func main() {
	var appID string
	var reviewFilter int
	flag.StringVar(&appID, "appid", "", "the apple app store numeric id")
	flag.IntVar(&reviewFilter, "filter-review", 0, "the apple app store numeric id")
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
			printHeadersOnce(entry)
			entry = entry.Filter(reviewGreaterThan(entry, reviewFilter))
			if len(entry.QuickRow()) > 0 {
				buf.WriteString(strings.Join(entry.QuickRow(), "\t") + "\n")
			}
		}
	}
	fmt.Println(buf.String())
}

func reviewGreaterThan(e parser.Entry, n int) func(parser.Entry) bool {
	i, err := strconv.Atoi(e.ImRating.String())
	if err != nil {
		return func(e parser.Entry) bool { return false }
	}
	switch n {
	case 0:
		// return everything
		return func(e parser.Entry) bool { return false }
	case 1:
		return func(e parser.Entry) bool {
			if i > 1 {
				return true
			}
			return false
		}
	case 2:
		return func(e parser.Entry) bool {
			if i > 2 {
				return true
			}
			return false
		}
	case 3:
		return func(e parser.Entry) bool {
			if i > 3 {
				return true
			}
			return false
		}
	case 4:
		return func(e parser.Entry) bool {
			if i > 4 {
				return true
			}
			return false
		}
	default:
		// return everything
		return func(e parser.Entry) bool { return false }
	}
}
