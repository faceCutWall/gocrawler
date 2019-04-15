package engine

import "log"

type ConcurrentEngine struct {
	Scheduler   Scheduler
	WorkerCount int
	ItemsChan   chan Item
	RequestProcessor Processor
}

type Processor func (Request)(ParseResult, error)

type Scheduler interface {
	Run()
	GetRequestChan() chan Request
	Submit(Request)
	ReadyNotifier
}

type ReadyNotifier interface {
	WorkerReady(chan Request)
}

func (p *ConcurrentEngine) Run(seeds ...Request) {
	p.Scheduler.Run()
	out := make(chan ParseResult)
	for i := 0; i < p.WorkerCount; i++ {
		p.createWorker(p.Scheduler.GetRequestChan(), out, p.Scheduler)
	}

	for _, seed := range seeds {
		if !isDuplicate(seed.Url) {
			p.Scheduler.Submit(seed)
		}
	}

	for {
		result := <-out
		for _, item := range result.Items {
			p.ItemsChan <- item

		}
		for _, request := range result.Requests {
			if !isDuplicate(request.Url) {
				p.Scheduler.Submit(request)
			}
		}
	}
}

func (p *ConcurrentEngine)createWorker(in chan Request, out chan ParseResult, ready ReadyNotifier) {
	go func() {
		for {
			ready.WorkerReady(in)
			request := <-in
			result, err := p.RequestProcessor(request)
			if err != nil {
				log.Printf("work error: %v", err)
				continue
			}
			out <- result
		}
	}()
}
