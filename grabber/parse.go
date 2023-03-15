package grabber

type RuleTree struct {
	Root  func() []*Request //根节点
	Trunk map[string]*Rule  //规则存储
}

//采集规则节点
type Rule struct {
	ParseFunc func(*Context) ParseResult
}

type Context struct {
	Body []byte
	Req  *Request
}
