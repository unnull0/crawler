package grabber

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"
)

type Task struct {
	URL      string
	Cookie   string
	Reload   bool
	WaitTime time.Duration
	MaxDepth int
	RootReq  *Request
	Fetcher  Fetcher
}

type Request struct {
	Task      *Task
	URL       string
	Method    string
	Depth     int
	ParseFunc func([]byte, *Request) ParseResult
}

type ParseResult struct {
	Requests []*Request
	Items    []interface{}
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
