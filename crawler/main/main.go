package main

import (
	"crawler/Beiwo/parser"
	"crawler/engine"
	"crawler/persist"
	"crawler/scheduler"
	"crawler_distributed/config"
)

func main() {
	itemChan, err := persist.ItemSaver("movie")
	if err != nil {
		panic(err)
	}
	e := engine.ConcurrentEngine{
		Scheduler:        &scheduler.QueueScheduler{},
		WorkerCount:      10,
		ItemsChan:        itemChan,
		RequestProcessor: engine.Worker,
	}
	e.Run(engine.Request{
		Url:    "http://www.beiwo888.com",
		Parser: engine.NewFuncParser(parser.MovieTypeParser, config.MovieTypeParser),
	})
}
