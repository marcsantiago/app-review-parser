package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"sort"
	"strconv"
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
		text := strings.TrimSpace(strings.ToLower(tok.Text))
		if !isStopWord(text) && validTags(tok.Tag) {
			count[text]++
		}
	}

	counts := make(countsSorted, 0, len(count))
	for k, v := range count {
		counts = append(counts, countSorted{Text: k, Count: v})
	}

	sort.Sort(counts)
	fmt.Println(counts.String())

}

type countSorted struct {
	Text  string
	Count int
}

type countsSorted []countSorted

func (c countsSorted) Len() int           { return len(c) }
func (c countsSorted) Less(i, j int) bool { return c[i].Count > c[j].Count }
func (c countsSorted) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

func (c countsSorted) String() string {
	var buf strings.Builder
	for _, e := range c {
		buf.WriteString("Word: ")
		buf.WriteString(e.Text)
		buf.WriteString(", Count: ")
		buf.WriteString(strconv.Itoa(e.Count))
		buf.WriteString("\n")
	}
	return buf.String()
}

func validTags(tag string) bool {
	switch tag {
	case "JJ", "JJR", "JJS", "RB", "RBR", "RBS", "RP", "VB", "VBD", "VBG", "VBN", "VBZ", "WRB":
		return true
	}
	return false
}

func isStopWord(tar string) bool {
	if strings.Contains(tar, "'") {
		return true
	}

	n := sort.SearchStrings(stopWords, tar)
	return n < len(stopWords) && stopWords[n] == tar
}

var stopWords = []string{
	"i",
	"me",
	"my",
	"myself",
	"we",
	"our",
	"ours",
	"ourselves",
	"you",
	"your",
	"yours",
	"yourself",
	"yourselves",
	"he",
	"him",
	"his",
	"himself",
	"she",
	"her",
	"hers",
	"herself",
	"it",
	"its",
	"itself",
	"they",
	"them",
	"their",
	"theirs",
	"themselves",
	"what",
	"which",
	"who",
	"whom",
	"this",
	"that",
	"these",
	"those",
	"am",
	"is",
	"are",
	"was",
	"were",
	"be",
	"been",
	"being",
	"have",
	"has",
	"had",
	"having",
	"do",
	"does",
	"did",
	"doing",
	"a",
	"an",
	"the",
	"and",
	"but",
	"if",
	"or",
	"because",
	"as",
	"until",
	"while",
	"of",
	"at",
	"by",
	"for",
	"with",
	"about",
	"against",
	"between",
	"into",
	"through",
	"during",
	"before",
	"after",
	"above",
	"below",
	"to",
	"from",
	"up",
	"down",
	"in",
	"out",
	"on",
	"off",
	"over",
	"under",
	"again",
	"further",
	"then",
	"once",
	"here",
	"there",
	"when",
	"where",
	"why",
	"how",
	"all",
	"any",
	"both",
	"each",
	"few",
	"more",
	"most",
	"other",
	"some",
	"such",
	"no",
	"nor",
	"not",
	"only",
	"own",
	"same",
	"so",
	"than",
	"too",
	"very",
	"s",
	"t",
	"can",
	"will",
	"just",
	"don",
	"should",
	"nowi",
	"me",
	"my",
	"myself",
	"we",
	"our",
	"ours",
	"ourselves",
	"you",
	"your",
	"yours",
	"yourself",
	"yourselves",
	"he",
	"him",
	"his",
	"himself",
	"she",
	"her",
	"hers",
	"herself",
	"it",
	"its",
	"itself",
	"they",
	"them",
	"their",
	"theirs",
	"themselves",
	"what",
	"which",
	"who",
	"whom",
	"this",
	"that",
	"these",
	"those",
	"am",
	"is",
	"are",
	"was",
	"were",
	"be",
	"been",
	"being",
	"have",
	"has",
	"had",
	"having",
	"do",
	"does",
	"did",
	"doing",
	"a",
	"an",
	"the",
	"and",
	"but",
	"if",
	"or",
	"because",
	"as",
	"until",
	"while",
	"of",
	"at",
	"by",
	"for",
	"with",
	"about",
	"against",
	"between",
	"into",
	"through",
	"during",
	"before",
	"after",
	"above",
	"below",
	"to",
	"from",
	"up",
	"down",
	"in",
	"out",
	"on",
	"off",
	"over",
	"under",
	"again",
	"further",
	"then",
	"once",
	"here",
	"there",
	"when",
	"where",
	"why",
	"how",
	"all",
	"any",
	"both",
	"each",
	"few",
	"more",
	"most",
	"other",
	"some",
	"such",
	"no",
	"nor",
	"not",
	"only",
	"own",
	"same",
	"so",
	"than",
	"too",
	"very",
	"s",
	"t",
	"can",
	"will",
	"just",
	"don",
	"should",
	"nowi",
	"me",
	"my",
	"myself",
	"we",
	"our",
	"ours",
	"ourselves",
	"you",
	"your",
	"yours",
	"yourself",
	"yourselves",
	"he",
	"him",
	"his",
	"himself",
	"she",
	"her",
	"hers",
	"herself",
	"it",
	"its",
	"itself",
	"they",
	"them",
	"their",
	"theirs",
	"themselves",
	"what",
	"which",
	"who",
	"whom",
	"this",
	"that",
	"these",
	"those",
	"am",
	"is",
	"are",
	"was",
	"were",
	"be",
	"been",
	"being",
	"have",
	"has",
	"had",
	"having",
	"do",
	"does",
	"did",
	"doing",
	"a",
	"an",
	"the",
	"and",
	"but",
	"if",
	"or",
	"because",
	"as",
	"until",
	"while",
	"of",
	"at",
	"by",
	"for",
	"with",
	"about",
	"against",
	"between",
	"into",
	"through",
	"during",
	"before",
	"after",
	"above",
	"below",
	"to",
	"from",
	"up",
	"down",
	"in",
	"out",
	"on",
	"off",
	"over",
	"under",
	"again",
	"further",
	"then",
	"once",
	"here",
	"there",
	"when",
	"where",
	"why",
	"how",
	"all",
	"any",
	"both",
	"each",
	"few",
	"more",
	"most",
	"other",
	"some",
	"such",
	"no",
	"nor",
	"not",
	"only",
	"own",
	"same",
	"so",
	"than",
	"too",
	"very",
	"s",
	"t",
	"can",
	"will",
	"just",
	"don",
	"should",
	"now",
}
