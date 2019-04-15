package worker

import (
	"crawler/Beiwo/parser"
	"crawler/engine"
	"crawler_distributed/config"
	"fmt"
	"github.com/pkg/errors"
	"log"
)

type SerializedParser struct {
	Name string
	Args interface{}
}

type Request struct {
	Url    string
	Parser SerializedParser
}

type ParseResult struct {
	Items   []engine.Item
	Request []Request
}

func SerializeRequest(r engine.Request) Request {
	name, args := r.Parser.Serialize()
	return Request{
		Url: r.Url,
		Parser: SerializedParser{
			Name: name,
			Args: args,
		},
	}
}

func SerializeResult(r engine.ParseResult) ParseResult {

	result := ParseResult{
		Items: r.Items,
	}
	for _, r := range r.Requests {
		result.Request = append(result.Request, SerializeRequest(r))
	}
	return result
}

func deserializeRequest(r Request) (engine.Request, error) {
	parser, err := DeserializeParser(r.Parser)
	if err != nil {
		return engine.Request{}, err
	}
	return engine.Request{
		Url:    r.Url,
		Parser: parser,
	}, nil

}
func DeserializeResult(r ParseResult) engine.ParseResult {
	result := engine.ParseResult{
		Items: r.Items,
	}
	for _, r := range r.Request {
		engineRequest, err := deserializeRequest(r)
		if err != nil {
			log.Printf("error deserializing request:%v", err)
			continue
		}
		result.Requests = append(result.Requests, engineRequest)
	}
	return result
}

func DeserializeParser(p SerializedParser) (engine.Parser, error) {
	switch p.Name {
	case config.MovieInfoParser:
		if movieName, ok := p.Args.(string); ok {
			return parser.NewMovieInfoFuncParser(movieName), nil
		} else {
			return nil, fmt.Errorf("invaild arg: %v", p.Args)
		}
	case config.MovieListParser:
		return engine.NewFuncParser(parser.MovieListParser, config.MovieListParser), nil
	case config.MovieTypeParser:
		return engine.NewFuncParser(parser.MovieTypeParser, config.MovieTypeParser), nil
	default:
		return nil, errors.New("unknown parser name")
	}

}

