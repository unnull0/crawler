package grabber

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"

	"github.com/unnull0/crawler/collector"
)

type Task struct {
	Name     string
	URL      string
	Cookie   string
	Reload   bool
	WaitTime time.Duration
	MaxDepth int
	// RootReq  *Request
	Fetcher Fetcher
	Rule    RuleTree
	Storage collector.Storage
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
