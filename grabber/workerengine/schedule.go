package workerengine

import (
	"github.com/unnull0/crawler/grabber"
	"go.uber.org/zap"
)

type Schedule struct {
	requestCh chan *grabber.Request
	workerCh  chan *grabber.Request
	out       chan grabber.ParseResult
	options
}

func NewSchedule(opts ...Option) *Schedule {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	s := &Schedule{}
	s.options = options
	return s
}

func (s *Schedule) Run() {
	requestch := make(chan *grabber.Request)
	workerCh := make(chan *grabber.Request)
	out := make(chan grabber.ParseResult)
	s.requestCh = requestch
	s.workerCh = workerCh
	s.out = out
	go s.Schedule()
	for i := 0; i < s.WorkCount; i++ {
		go s.CreateWork()
	}
	s.HandleResult()
}

func (s *Schedule) Schedule() {

	var req *grabber.Request
	var workerCh chan *grabber.Request
	for {
		if req == nil && len(s.Seeds) > 0 {
			req = s.Seeds[0]
			s.Seeds = s.Seeds[1:]
			workerCh = s.workerCh
		}
		select {
		case r := <-s.requestCh:
			s.Seeds = append(s.Seeds, r)
		case workerCh <- req:
			req = nil
			workerCh = nil
		}
	}

}

func (s *Schedule) CreateWork() {
	for {
		r := <-s.workerCh
		body, err := s.Fetcher.Get(r)
		if err != nil {
			s.Logger.Error("get content error", zap.Error(err))
			continue
		}
		result := r.ParseFunc(body, r)
		s.out <- result
	}
}

func (s *Schedule) HandleResult() {
	for {
		select {
		case result := <-s.out:
			for _, req := range result.Requests {
				s.requestCh <- req
			}
			for _, item := range result.Items {
				s.Logger.Info("result", zap.String("url", item.(string)))
			}
		}
	}
}

func (s *Schedule) run() {
	requestCh := make(chan *grabber.Request)
	workerCh := make(chan *grabber.Request)
	out := make(chan grabber.ParseResult)
	s.requestCh = requestCh
	s.workerCh = workerCh
	s.out = out
	go s.Schedule()
	for i := 0; i < s.WorkCount; i++ {
		go s.CreateWork()
	}
	s.HandleResult()
}
