package doubantenement

import (
	"regexp"

	"github.com/unnull0/crawler/grabber"
)

//获取网页所有租房文章链接，内容包含阳台就返回对应文章的链接

const cityListRe = `(https://www.douban.com/group/topic/[0-9a-z]+/)"[^>]*>([^<]+)</a>`

func ParseURL(contents []byte, req *grabber.Request) grabber.ParseResult {
	re := regexp.MustCompile(cityListRe)
	matches := re.FindAllSubmatch(contents, -1) //找到所有匹配的正则内容
	result := grabber.ParseResult{}

	for _, m := range matches {
		URL := string(m[1])
		result.Requests = append(result.Requests, &grabber.Request{
			URL:    URL,
			Method: "GET",
			Task:   req.Task,
			Depth:  req.Depth + 1,
			ParseFunc: func(c []byte, request *grabber.Request) grabber.ParseResult {
				return GetContent(c, URL)
			},
		})
	}
	return result
}

const ContentRe = `<div class="topic-content">[\s\S]*?阳台\s\S]*?<div`

func GetContent(contents []byte, URL string) grabber.ParseResult {
	re := regexp.MustCompile(ContentRe)
	ok := re.Match(contents) //返回是否有匹配的正则内容
	if !ok {
		return grabber.ParseResult{
			Items: []interface{}{},
		}
	}

	result := grabber.ParseResult{
		Items: []interface{}{URL},
	}
	return result
}
