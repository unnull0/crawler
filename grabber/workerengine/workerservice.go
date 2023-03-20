package workerengine

import (
	"sync"

	"github.com/unnull0/crawler/collector"
	"github.com/unnull0/crawler/grabber"
	"go.uber.org/zap"
)

type Engine struct {
	out         chan grabber.ParseResult
	Visited     map[string]bool
	VisitedLock sync.Mutex

	failures    map[string]*grabber.Request
	failureLock sync.Mutex
	options
}

func NewEngine(opts ...Option) *Engine {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}

	e := &Engine{}
	e.Visited = make(map[string]bool)
	e.failures = make(map[string]*grabber.Request)
	e.options = options
	e.out = make(chan grabber.ParseResult)
	return e
}

func (e *Engine) Run() {
	go e.Schedule()
	for i := 0; i < e.WorkCount; i++ {
		go e.CreateWork()
	}
	e.HandleResult()
}

func (e *Engine) Schedule() {
	var reqs []*grabber.Request
	for _, seed := range e.Seeds {
		task := Tkstore.Hash[seed.Name]
		task.Fetcher = seed.Fetcher
		task.Storage = seed.Storage
		rootReqs := task.Rule.Root()
		for _, req := range rootReqs {
			req.Task = task
		}
		reqs = append(reqs, rootReqs...)
	}
	go e.scheduler.Schedule()
	go e.scheduler.Push(reqs...)
}

func (e *Engine) CreateWork() {
	for {
		r := e.scheduler.Pull()
		if err := r.Check(); err != nil {
			e.Logger.Error("check failed", zap.Error(err))
			continue
		}
		if !r.Task.Reload && e.HasVisited(r) {
			continue
		}
		e.StoreVisited(r)

		body, err := r.Task.Fetcher.Get(r)
		if err != nil {
			e.Logger.Error("get content error", zap.Error(err))
			e.SetFailure(r)
			continue
		}

		if len(body) < 6000 {
			e.Logger.Error("fetch failed", zap.Int("body length", len(body)), zap.String("URL", r.URL))
			e.SetFailure(r)
			continue
		}

		rule := r.Task.Rule.Trunk[r.RuleName]
		result := rule.ParseFunc(&grabber.Context{Body: body, Req: r})
		if len(result.Requests) > 0 {
			go e.scheduler.Push(result.Requests...)
		}
		e.out <- result
	}
}

func (e *Engine) HandleResult() {
	for {
		select {
		case result := <-e.out:
			for _, item := range result.Items {
				switch d := item.(type) {
				case *collector.DataCell:
					name := d.GetTaskName()
					task := Tkstore.Hash[name]
					task.Storage.Save(d)
				}
				e.Logger.Info("result", zap.String("url", item.(string)))
			}
		}
	}
}

func (e *Engine) HasVisited(r *grabber.Request) bool {
	e.VisitedLock.Lock()
	defer e.VisitedLock.Unlock()
	unique := r.Unique()
	return e.Visited[unique]
}

func (e *Engine) StoreVisited(reqs ...*grabber.Request) {
	e.VisitedLock.Lock()
	defer e.VisitedLock.Unlock()
	for _, r := range reqs {
		unique := r.Unique()
		e.Visited[unique] = true
	}
}

func (e *Engine) SetFailure(req *grabber.Request) {
	if !req.Task.Reload {
		e.VisitedLock.Lock()
		unique := req.Unique()
		delete(e.Visited, unique)
		e.VisitedLock.Unlock()
	}
	e.failureLock.Lock()
	defer e.failureLock.Unlock()
	if _, ok := e.failures[req.Unique()]; !ok {
		e.failures[req.Unique()] = req
		e.scheduler.Push(req)
	}
}
