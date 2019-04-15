package scheduler

import "crawler/engine"

type QueueScheduler struct {
	RequestChan chan engine.Request
	WorkerChan  chan chan engine.Request
}

func (p *QueueScheduler) GetRequestChan() chan engine.Request {
	return make(chan engine.Request)
}

func (p *QueueScheduler) Submit(request engine.Request) {
		p.RequestChan <- request
}

func (p *QueueScheduler) WorkerReady(ch chan engine.Request) {
	p.WorkerChan <- ch
}

func (p *QueueScheduler) Run() {
	p.RequestChan = make(chan engine.Request)
	p.WorkerChan = make(chan chan engine.Request)
	var requestQ []engine.Request
	var workerQ []chan engine.Request
	go func() {
		for {
			var activeRequest engine.Request
			var activeWorker chan engine.Request
			if len(requestQ) > 0 && len(workerQ) > 0 {
				activeRequest = requestQ[0]
				activeWorker = workerQ[0]
			}
			select {
			case r := <-p.RequestChan:
				requestQ = append(requestQ, r)
			case w := <-p.WorkerChan:
				workerQ = append(workerQ, w)
			case activeWorker <- activeRequest:
				requestQ = requestQ[1:]
				workerQ = workerQ[1:]
			}
		}
	}()

}
