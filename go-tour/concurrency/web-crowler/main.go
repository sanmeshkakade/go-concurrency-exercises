package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type FetchRecorder interface {
	WasFetched(url string) bool
	AddFetched(url string)
}

func NewFetchRecorder() FetchRecorder {
	return &fetchRecorder{
		fMap: map[string]bool{},
		mu:   &sync.Mutex{},
	}
}

type fetchRecorder struct {
	fMap map[string]bool
	mu   *sync.Mutex
}

func (fr *fetchRecorder) WasFetched(url string) bool {
	fr.mu.Lock()
	defer fr.mu.Unlock()
	_, exists := fr.fMap[url]
	return exists
}

func (fr *fetchRecorder) AddFetched(url string) {
	fr.mu.Lock()
	defer fr.mu.Unlock()
	fr.fMap[url] = true
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, wg *sync.WaitGroup, fr FetchRecorder) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:

	defer wg.Done()

	if depth <= 0 {
		return
	}

	if fr.WasFetched(url) {
		return
	}

	defer fr.AddFetched(url)

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)

	wg.Add(len(urls))
	for _, u := range urls {
		go Crawl(u, depth-1, fetcher, wg, fr)
	}
	return
}

func main() {
	var fr FetchRecorder = NewFetchRecorder()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go Crawl("https://golang.org/", 4, fetcher, wg, fr)
	wg.Wait()
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
