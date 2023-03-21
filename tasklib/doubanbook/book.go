package doubanbook

import (
	"regexp"
	"strconv"

	"github.com/unnull0/crawler/grabber"
)

var DoubanBookTask = &grabber.Task{
	Name:     "douban_book_list",
	WaitTime: 5,
	MaxDepth: 5,
	Cookie:   `douban-fav-remind=1; viewed="1007305"; bid=-LBsoZdQj8k; gr_user_id=9381f0c8-c9a0-4ef7-9e10-1b4812fb785e; __gads=ID=a486bca57b09d9fc-224d131d48d900aa:T=1673615555:RT=1673615555:S=ALNI_MZvaEzveWJzMynGcfoXl7iHp8T1WA; __yadk_uid=fLQE8dmo4P1yoZDK9bOwSzQusWeSs75z; push_noty_num=0; push_doumail_num=0; __utmv=30149280.26736; ap_v=0,6.0; __gpi=UID=00000ba35654ce16:T=1673615555:RT=1678430610:S=ALNI_Mb-V6rgnowh8JznMa81jye-o5mh8g; apiKey=; dbcl2="267368696:rLWKgwWt9MY"; ck=k7w3; _pk_ref.100001.8cb4=["","",1678433761,"https://accounts.douban.com/"]; _pk_id.100001.8cb4=55cd98221bb9ae38.1612278505.9.1678433761.1678430715.; _pk_ses.100001.8cb4=*; __utma=30149280.1521688018.1592847100.1678430612.1678433762.13; __utmc=30149280; __utmz=30149280.1678433762.13.4.utmcsr=accounts.douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; __utmt=1; __utmb=30149280.7.5.1678433762`,
	Rule: grabber.RuleTree{
		Root: func() []*grabber.Request {
			roots := []*grabber.Request{
				&grabber.Request{
					Method:   "GET",
					URL:      "https://book.douban.com",
					RuleName: "解析卡片",
				},
			}
			return roots
		},
		Trunk: map[string]*grabber.Rule{
			"解析卡片": &grabber.Rule{ParseFunc: ParseTag},
			"书籍列表": &grabber.Rule{ParseFunc: ParseBookList},

			"书籍简介": &grabber.Rule{
				ParseFunc:  ParseBookDetail,
				ItemFields: []string{"书名", "作者", "页数", "出版社", "得分", "价格", "简介"},
			},
		},
	},
}

const tagRe = `<a href="([^"]+)" class="tag">([^<]+)</a>`

//rulename解析卡片
func ParseTag(ctx *grabber.Context) grabber.ParseResult {
	re := regexp.MustCompile(tagRe)
	matches := re.FindAllSubmatch(ctx.Body, -1)
	// matches := re.FindSubmatch(ctx.Body)
	result := grabber.ParseResult{}

	for _, m := range matches {
		url := "https://book.douban.com" + string(m[1])
		// url := "https://book.douban.com" + string(matches[1])
		result.Requests = append(result.Requests, &grabber.Request{
			URL:      url,
			Task:     ctx.Req.Task,
			Method:   "GET",
			Depth:    ctx.Req.Depth + 1,
			RuleName: "书籍列表",
		})
	}
	return result
}

const booklistRe = `<a href="([^"]+)" title="([^"]+)"`

//rulename书籍列表
func ParseBookList(ctx *grabber.Context) grabber.ParseResult {
	re := regexp.MustCompile(booklistRe)
	matches := re.FindAllSubmatch(ctx.Body, -1)
	result := grabber.ParseResult{}

	for _, m := range matches {
		req := &grabber.Request{
			Method:   "GET",
			URL:      string(m[1]),
			Task:     ctx.Req.Task,
			Depth:    ctx.Req.Depth + 1,
			RuleName: "书籍简介",
		}
		req.TmpData = &grabber.Temp{}
		req.TmpData.Set("书名", string(m[2]))
		result.Requests = append(result.Requests, req)
	}
	return result
}

var autoRe = regexp.MustCompile(`<span class="pl"> 作者</span>:[\d\D]*?<a.*?>([^<]+)</a>`)
var public = regexp.MustCompile(`<span class="pl">出版社:</span>([^<]+)<br/>`)
var pageRe = regexp.MustCompile(`<span class="pl">页数:</span> ([^<]+)<br/>`)
var priceRe = regexp.MustCompile(`<span class="pl">定价:</span>([^<]+)<br/>`)
var scoreRe = regexp.MustCompile(`<strong class="ll rating_num " property="v:average">([^<]+)</strong>`)
var intoRe = regexp.MustCompile(`<div class="intro">[\d\D]*?<p>([^<]+)</p></div>`)

//书籍简介
func ParseBookDetail(ctx *grabber.Context) grabber.ParseResult {
	bookName := ctx.Req.TmpData.Get("书名")
	page, _ := strconv.Atoi(GetRegexpStr(ctx.Body, pageRe))

	book := map[string]interface{}{
		"书名":  bookName,
		"作者":  GetRegexpStr(ctx.Body, autoRe),
		"页数":  page,
		"出版社": GetRegexpStr(ctx.Body, public),
		"得分":  GetRegexpStr(ctx.Body, scoreRe),
		"价格":  GetRegexpStr(ctx.Body, priceRe),
		"简介":  GetRegexpStr(ctx.Body, intoRe),
	}
	data := ctx.Output(book)

	result := grabber.ParseResult{Items: []interface{}{data}}
	return result
}

func GetRegexpStr(content []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(content)

	if len(match) >= 2 {
		return string(match[1])
	}
	return ""
}
