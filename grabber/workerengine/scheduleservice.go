package workerengine

import (
	"github.com/unnull0/crawler/grabber"
	"go.uber.org/zap"
)

type Scheduler interface {
	Schedule()
	Push(...*grabber.Request)
	Pull() *grabber.Request
}

type Schedule struct {
	requestCh chan *grabber.Request
	workerCh  chan *grabber.Request
	reqQueue  []*grabber.Request
	Logger    *zap.Logger
}

func NewSchedule() *Schedule {
	s := &Schedule{}
	s.requestCh = make(chan *grabber.Request)
	s.workerCh = make(chan *grabber.Request)
	return s
}

func (s *Schedule) Schedule() {
	var req *grabber.Request
	var workerCh chan *grabber.Request
	for {
		if req == nil && len(s.reqQueue) > 0 {
			req = s.reqQueue[0]
			s.reqQueue = s.reqQueue[1:]
			workerCh = s.workerCh
		}
		select {
		case r := <-s.requestCh:
			s.reqQueue = append(s.reqQueue, r)
		case workerCh <- req:
			req = nil
			workerCh = nil
		}
	}

}

func (s *Schedule) Push(reqs ...*grabber.Request) {
	for _, req := range reqs {
		s.requestCh <- req
	}
}

func (s *Schedule) Pull() *grabber.Request {
	r := <-s.workerCh
	return r
}
