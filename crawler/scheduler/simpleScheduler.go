package scheduler

import "crawler/engine"

type simpleScheduler struct {
	RequestChan chan engine.Request
}

func (p *simpleScheduler) Run() {
	p.RequestChan = make(chan engine.Request)
}

func (p *simpleScheduler) GetRequestChan() chan engine.Request {
	return p.RequestChan
}
func (p *simpleScheduler) Submit(r engine.Request) {
	go func() {
		p.RequestChan <- r
	}()
}

func (p *simpleScheduler) WorkerReady(ch chan engine.Request) {
	p.RequestChan = ch
}
