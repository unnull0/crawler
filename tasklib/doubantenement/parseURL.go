package doubantenement

import (
	"fmt"
	"regexp"
	"time"

	"github.com/unnull0/crawler/grabber"
)

var DoubantenementTask = &grabber.Task{
	Name:     "find_douban_sun_rom",
	WaitTime: 1 * time.Second,
	Cookie:   `douban-fav-remind=1; viewed="1007305"; bid=-LBsoZdQj8k; gr_user_id=9381f0c8-c9a0-4ef7-9e10-1b4812fb785e; __gads=ID=a486bca57b09d9fc-224d131d48d900aa:T=1673615555:RT=1673615555:S=ALNI_MZvaEzveWJzMynGcfoXl7iHp8T1WA; __yadk_uid=fLQE8dmo4P1yoZDK9bOwSzQusWeSs75z; push_noty_num=0; push_doumail_num=0; __utmv=30149280.26736; ap_v=0,6.0; __gpi=UID=00000ba35654ce16:T=1673615555:RT=1678430610:S=ALNI_Mb-V6rgnowh8JznMa81jye-o5mh8g; apiKey=; dbcl2="267368696:rLWKgwWt9MY"; ck=k7w3; _pk_ref.100001.8cb4=["","",1678433761,"https://accounts.douban.com/"]; _pk_id.100001.8cb4=55cd98221bb9ae38.1612278505.9.1678433761.1678430715.; _pk_ses.100001.8cb4=*; __utma=30149280.1521688018.1592847100.1678430612.1678433762.13; __utmc=30149280; __utmz=30149280.1678433762.13.4.utmcsr=accounts.douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; __utmt=1; __utmb=30149280.7.5.1678433762`,
	MaxDepth: 5,
	Rule: grabber.RuleTree{
		Root: func() []*grabber.Request {
			var root []*grabber.Request
			for i := 0; i < 25; i += 25 {
				str := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", i)
				root = append(root, &grabber.Request{
					URL:      str,
					Method:   "GET",
					RuleName: "解析网站URL",
				})
			}
			return root
		},
		Trunk: map[string]*grabber.Rule{
			"解析网站URL": &grabber.Rule{ParseURL},
			"扫描阳台":    &grabber.Rule{GetContent},
		},
	},
}

//获取网页所有租房文章链接，内容包含阳台就返回对应文章的链接

const cityListRe = `(https://www.douban.com/group/topic/[0-9a-z]+/)"[^>]*>([^<]+)</a>`

func ParseURL(ctx *grabber.Context) grabber.ParseResult {
	re := regexp.MustCompile(cityListRe)
	matches := re.FindAllSubmatch(ctx.Body, -1) //找到所有匹配的正则内容
	result := grabber.ParseResult{}

	for _, m := range matches {
		URL := string(m[1])
		result.Requests = append(result.Requests, &grabber.Request{
			URL:      URL,
			Method:   "GET",
			Task:     ctx.Req.Task,
			Depth:    ctx.Req.Depth + 1,
			RuleName: "扫描阳台",
		})
	}
	return result
}

const ContentRe = `<div class="topic-content">[\s\S]*?阳台[\s\S]*?<div class="aside">`

func GetContent(ctx *grabber.Context) grabber.ParseResult {
	re := regexp.MustCompile(ContentRe)
	ok := re.Match(ctx.Body) //返回是否有匹配的正则内容
	if !ok {
		return grabber.ParseResult{
			Items: []interface{}{},
		}
	}

	result := grabber.ParseResult{
		Items: []interface{}{ctx.Req.URL},
	}
	return result
}
