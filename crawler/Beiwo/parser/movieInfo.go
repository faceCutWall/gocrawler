package parser

import (
	"crawler/engine"
	"crawler/model"
	"regexp"
)

var (
	movieActorsRegex      = regexp.MustCompile(`<a href="/index\.php\?s=vod-search-actor-.+\.html" target="_blank">([^<]+)</a>`)
	movieTagsRegex        = regexp.MustCompile(`<a href="/index\.php\?s=vod-search-tid-[0-9]+\.html">([^<]+)</a>`)
	movieReleaseTimeRegex = regexp.MustCompile(`<li>上映年代：([0-9]+)&nbsp;`)
	movieIdRegex          = regexp.MustCompile(`vod/([0-9]+)`)
)

func MovieInfoParser(body []byte, name string, url string) engine.ParseResult {
	var movie model.MovieInfo
	movie.Name = name
	movie.ReleaseTime = assignString(body, movieReleaseTimeRegex)

	matchActors := movieActorsRegex.FindAllSubmatch(body, -1)
	matchTags := movieTagsRegex.FindAllSubmatch(body, -1)

	for _, subMatchActors := range matchActors {
		for _, m := range subMatchActors[1:] {
			movie.Actors = append(movie.Actors, string(m))
		}
	}
	for _, subMatchTags := range matchTags {
		for _, m := range subMatchTags[1:] {
			if string(m) != "未知ID98" {
				movie.Tags = append(movie.Tags, string(m))
			}
		}
	}

	result := engine.ParseResult{
		Items: []engine.Item{
			{
				Url:     url,
				ID:      assignString([]byte(url), movieIdRegex),
				Type:    "beiwo",
				Payload: movie,
			},
		},
	}
	return result
}

func assignString(body []byte, regex *regexp.Regexp) string {
	subMatch := regex.FindSubmatch(body)
	return string(subMatch[1])
}
