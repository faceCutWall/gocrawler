package main

import (
	"crawler/Beiwo/parser"
	"crawler/engine"
	"crawler/scheduler"
	"crawler_distributed/config"
	itemSaver "crawler_distributed/persist/client"
	"crawler_distributed/rpcsupport"
	worker "crawler_distributed/worker/client"
	"flag"
	"log"
	"net/rpc"
	"strings"
)

var(
	itemSaverHost = flag.String("itemsaver_host","","itemsaver host")
	workerHost = flag.String("worker_hosts","","worker host (comma separated)")
)

func main() {
	flag.Parse()

	itemChan, err := itemSaver.ItemSaver(*itemSaverHost)
	if err != nil {
		panic(err)
	}
	pool := createClientPool( strings.Split(*workerHost,",") )
	processor, err := worker.CreateProcessor(pool)
	if err != nil {
		panic(err)
	}
	e := engine.ConcurrentEngine{
		Scheduler:        &scheduler.QueueScheduler{},
		WorkerCount:      10,
		ItemsChan:        itemChan,
		RequestProcessor: processor,
	}
	e.Run(engine.Request{
		Url:    "http://www.beiwo888.com",
		Parser: engine.NewFuncParser(parser.MovieTypeParser, config.MovieTypeParser),
	})
}

func createClientPool(hosts []string) chan *rpc.Client {
	var clients []*rpc.Client
	for _, host := range hosts {
		client, err := rpcsupport.NewClient(host)
		if err == nil {
			clients = append(clients, client)
			log.Printf("connected to %s", host)
		} else {
			log.Printf("error connecting to %s:%v", host, err)
		}
	}
	out := make(chan *rpc.Client)
	go func() {
		for {
			for _, client := range clients {
				out <- client
			}
		}
	}()
	return out
}
