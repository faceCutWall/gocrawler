package parser

import (
	"crawler/engine"
	"crawler_distributed/config"
	"regexp"
)

var movieTypeRegex = regexp.MustCompile(
	`<a href="(/index\.php\?s=vod-search-id-1-tid-[0-9]+\.html)">([^<]+)</a>`)

func MovieTypeParser(body []byte, _ string) engine.ParseResult {
	subMatch := movieTypeRegex.FindAllSubmatch(body, -1)
	var result = engine.ParseResult{}
	for _, m := range subMatch {
		result.Requests = append(result.Requests, engine.Request{
			Url:       "http://www.beiwo888.com" + string(m[1]),
			Parser:    engine.NewFuncParser(MovieListParser, config.MovieListParser),
		})
	}
	return result
}
