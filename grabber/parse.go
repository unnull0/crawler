package grabber

import (
	"time"

	"github.com/unnull0/crawler/collector"
)

type RuleTree struct {
	Root  func() []*Request //根节点
	Trunk map[string]*Rule  //规则存储
}

//采集规则节点
type Rule struct {
	ItemFields []string
	ParseFunc  func(*Context) ParseResult
}

type Context struct {
	Body []byte
	Req  *Request
}

func (c *Context) Output(data interface{}) *collector.DataCell {
	res := &collector.DataCell{}
	res.Data = make(map[string]interface{})
	res.Data["Task"] = c.Req.Task.Name
	res.Data["Rule"] = c.Req.RuleName
	res.Data["Data"] = data
	res.Data["Url"] = c.Req.URL
	res.Data["Time"] = time.Now().Format("2006-01-02 15:04:05")
	return res
}
