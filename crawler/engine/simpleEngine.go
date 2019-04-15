package engine

import (
	"log"
)

type SimpleEngine struct {
}

func (p *SimpleEngine) Run(seeds ...Request) {
	var queue []Request
	for _, seed := range seeds {
		if !isDuplicate(seed.Url) {
			queue = append(queue, seed)
		}
	}
	for {
		request := queue[0]
		queue = queue[1:]
		result, err := Worker(request)
		if err != nil {
			log.Printf("work error:%v", err)
			continue
		}
		for _, item := range result.Items {
			log.Printf("got item:%v", item)
		}
		for _, request := range result.Requests {
			if !isDuplicate(request.Url) {
				queue = append(queue, request)
			}
		}
	}
}
