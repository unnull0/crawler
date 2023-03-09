package grabber

import "time"

//任务请求
type Request struct {
	URL       string
	Cookie    string
	Timeout   time.Duration
	ParseFunc func([]byte) ParseResult
}

type ParseResult struct {
	Requests []*Request
	Items    []interface{}
}
