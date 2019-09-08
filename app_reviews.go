package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	// appID, page
	iosReviewFeed = "https://itunes.apple.com/rss/customerreviews/id=%s/page=%s/sortby=mostrecent/json"
)

var client = &http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 1 * time.Minute,
		}).DialContext,
	},
	Timeout: 10 * time.Second,
}

var (
	source = rand.NewSource(time.Now().UnixNano())
	random = rand.New(source)
)

type page struct {
	mu        *sync.RWMutex
	position  int
	semaphore semaphore
}

func (p *page) next() string {
	p.mu.Lock()
	p.position++
	p.mu.Unlock()

	p.mu.RLock()
	n := p.position
	p.mu.RUnlock()
	return strconv.Itoa(n)
}

type semaphore chan struct{}

func (s semaphore) release() { <-s }
func (s semaphore) load() {
	// add some jitter
	time.Sleep(time.Duration(random.Intn(600)) * time.Millisecond)
	s <- struct{}{}
}

// FetchIOSAppReviews ...
func FetchIOSAppReviews(appID string) []IOSReviews {
	pagger := &page{mu: new(sync.RWMutex), semaphore: make(semaphore, 10)}
	results := make(chan IOSReviews)
	// will cancel to exit the for loop
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context, cancel context.CancelFunc, pagger *page, appID string, results chan IOSReviews) {
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			default:
				go func(ctx context.Context, cancel context.CancelFunc, pagger *page, appID string, results chan IOSReviews) {
					pagger.semaphore.load()
					defer pagger.semaphore.release()

					uri := fmt.Sprintf(iosReviewFeed, appID, pagger.next())

					req, err := http.NewRequest(http.MethodGet, uri, nil)
					if err != nil {
						cancel()
						return
					}
					req.Header.Add("User-Agent", RotatedUserAgents.GetUA())

					res, err := client.Do(req)
					if err != nil {
						return
					}
					defer res.Body.Close()

					if res.StatusCode != http.StatusOK {
						cancel()
						return
					}

					var result IOSReviews
					err = json.NewDecoder(res.Body).Decode(&result)
					if err != nil {
						cancel()
						return
					}
					results <- result

				}(ctx, cancel, pagger, appID, results)
			}
		}
	}(ctx, cancel, pagger, appID, results)

	var data []IOSReviews
	for {
		select {
		case <-ctx.Done():
			return data
		case res := <-results:
			data = append(data, res)
		}
	}
}
