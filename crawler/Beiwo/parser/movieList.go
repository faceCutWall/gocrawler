package parser

import (
	"crawler/engine"
	"crawler_distributed/config"
	"regexp"
)

var (
	movieListRegex = regexp.MustCompile(`<a href="(/vod/[0-9]+/)" title="([^"]+)" target="_blank">`)
	nextPageRegex  = regexp.MustCompile(`<a href="([^"]+)" class="pagegbk">下一页</a>`)
)

func MovieListParser(body []byte, _ string) engine.ParseResult {
	movieListSubMatch := movieListRegex.FindAllSubmatch(body, -1)
	nextPageSubMatch := nextPageRegex.FindAllSubmatch(body, -1)
	var result engine.ParseResult
	for _, m := range movieListSubMatch {
		name := string(m[2])
		url := "http://www.beiwo888.com" + string(m[1])
		result.Requests = append(result.Requests, engine.Request{
			Url:    url,
			Parser: NewMovieInfoFuncParser(name),
		})
	}
	for _, m := range nextPageSubMatch {
		result.Requests = append(result.Requests, engine.Request{
			Url:    "http://www.beiwo888.com" + string(m[1]),
			Parser: engine.NewFuncParser(MovieListParser, config.MovieListParser),
		})
	}
	return result
}

type MovieInfoFuncParser struct {
	movieName string
}

func (p *MovieInfoFuncParser) Parse(contents []byte, url string) engine.ParseResult {
	return MovieInfoParser(contents, p.movieName, url)
}

func (p *MovieInfoFuncParser) Serialize() (name string, args interface{}) {
	return config.MovieInfoParser, p.movieName
}

func NewMovieInfoFuncParser(name string) *MovieInfoFuncParser {
	return &MovieInfoFuncParser{
		movieName: name,
	}

}
