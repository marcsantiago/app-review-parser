package parser

import "strings"

// IOSReviews response strucure from the ios rss feed
type IOSReviews struct {
	Feed Feed `json:"feed"`
}

// Attributes ...
type Attributes struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
	Type string `json:"type"`
}

// Author ...
type Author struct {
	Label
	Name Name `json:"name"`
	URI  URI  `json:"uri"`
}

// Content ...
type Content struct {
	Attributes Attributes
	Label
}

// Entry ...
type Entry struct {
	Author        Author        `json:"author"`
	Content       Content       `json:"content"`
	ID            ID            `json:"id"`
	ImContentType ImContentType `json:"im:contentType"`
	ImRating      ImRating      `json:"im:rating"`
	ImVersion     ImVoteSum     `json:"im:version"`
	ImVoteCount   ImVoteCount   `json:"im:voteCount"`
	ImVoteSum     ImVoteSum     `json:"im:voteSum"`
	Link          Link          `json:"link"`
	Title         Title         `json:"title"`
}

// Filter allows ads to be filtered out of the overall slice
// because when QuickRow is called empty reviews are skipped
// and this method artificially makes the review blank
func (e Entry) Filter(conditions ...func(Entry) bool) Entry {
	for _, c := range conditions {
		if c(e) {
			return Entry{}
		}
	}
	return e
}

// QuickRow returns same quick data from the Entry struct that can be turned into a TSV
// returns nothing if the review message is empty
// review_id, title, author, author_url, version, rating, review, vote_count
func (e Entry) QuickRow() []string {
	if len(strings.TrimSpace(e.Content.String())) == 0 {
		return nil
	}

	// clean up the review abit
	review := strings.TrimSpace(e.Content.String())
	review = strings.ReplaceAll(review, "\n", " ")
	review = strings.ReplaceAll(review, "\t", " ")

	return []string{
		strings.TrimSpace(e.ID.String()),
		strings.TrimSpace(e.Title.String()),
		strings.TrimSpace(e.Author.Name.String()),
		strings.TrimSpace(e.Author.URI.String()),
		strings.TrimSpace(e.ImVersion.String()),
		strings.TrimSpace(e.ImRating.String()),
		review,
		strings.TrimSpace(e.ImVoteCount.String()),
	}
}

// QuickHeaders returns headers for QuickRow
// review_id, title, author, author_url, version, rating, review, vote_count
func (e Entry) QuickHeaders() []string {
	return []string{
		"review_id",
		"title",
		"author",
		"author_url",
		"version",
		"rating",
		"review",
		"vote_count",
	}
}

// Feed ...
type Feed struct {
	Author  Author  `json:"author"`
	Entry   []Entry `json:"entry"`
	Icon    Icon    `json:"icon"`
	ID      ID      `json:"id"`
	Link    []Link  `json:"link"`
	Rights  Rights  `json:"rights"`
	Title   Title   `json:"title"`
	Updated Updated `json:"updated"`
}

// Icon ...
type Icon struct {
	Label
}

// ID ...
type ID struct {
	Label
}

// ImContentType ...
type ImContentType struct {
	Attributes Attributes `json:"attributes,omitempty"`
}

// ImRating ...
type ImRating struct {
	Label
}

// ImVoteCount ...
type ImVoteCount struct {
	Label
}

// ImVoteSum ...
type ImVoteSum struct {
	Label
}

// Label ...
type Label struct {
	Label string `json:"label"`
}

func (l Label) String() string {
	return l.Label
}

// Link ...
type Link struct {
	Attributes Attributes `json:"attributes,omitempty"`
}

// Name ...
type Name struct {
	Label
}

// Rights ...
type Rights struct {
	Label
}

// Title ...
type Title struct {
	Label
}

// URI ...
type URI struct {
	Label
}

// Updated ...
type Updated struct {
	Label
}
