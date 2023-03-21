package grabber

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"math/rand"
	"time"

	"github.com/unnull0/crawler/collector"
	"github.com/unnull0/crawler/limiter"
)

type Task struct {
	Name     string
	URL      string
	Cookie   string
	Reload   bool
	WaitTime int
	MaxDepth int
	// RootReq  *Request
	Fetcher Fetcher
	Rule    RuleTree
	Storage collector.Storage
	Limit   limiter.RateLimiter
}

type Request struct {
	Task     *Task
	URL      string
	Method   string
	Depth    int
	RuleName string
	TmpData  *Temp
}

type ParseResult struct {
	Requests []*Request
	Items    []interface{}
}

func (r *Request) Fetch() ([]byte, error) {
	if err := r.Task.Limit.Wait(context.Background()); err != nil {
		return nil, err
	}

	sleepTime := rand.Intn(r.Task.WaitTime * 1000)
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	return r.Task.Fetcher.Get(r)
}

func (r *Request) Check() error {
	if r.Depth > r.Task.MaxDepth {
		return errors.New("overdepth")
	}
	return nil
}

func (r *Request) Unique() string {
	block := md5.Sum([]byte(r.URL + r.Method))
	return hex.EncodeToString(block[:])
}
