package main

import (
	"crawler_distributed/rpcsupport"
	"crawler_distributed/worker"
	"flag"
	"fmt"
	"log"
)

var port = flag.Int("port",0,"the port from me to listen on")



func main() {
	flag.Parse()
	if *port == 0{
		log.Printf("must specify a port")
		return
	}

	log.Fatal(rpcsupport.ServeRpc(fmt.Sprintf(":%d", *port),
		worker.CrawlService{}))
}
