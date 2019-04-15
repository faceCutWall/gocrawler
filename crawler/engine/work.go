package engine

import (
	"crawler/fetcher"
	"log"
)

func Worker(r Request) (ParseResult, error) {
	log.Printf("Fetching Url:%s", r.Url)
	bytes, err := fetcher.Fetcher(r.Url)
	if err != nil {
		return ParseResult{}, err
	}
	return r.Parser.Parse(bytes,r.Url),nil
}

var visited = make(map[string]bool)

func isDuplicate(url string) bool {
	if visited[url] {
		return true
	}
	visited[url] = true
	return false
}
