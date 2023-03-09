package main

import (
	"fmt"
	"time"

	"github.com/unnull0/crawler/grabber"
	"github.com/unnull0/crawler/log"
	"github.com/unnull0/crawler/tasklib/doubantenement"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	core, c := log.NewFileCore("./text.log", zapcore.InfoLevel)
	defer c.Close()
	logger := log.NewLogger(core)
	logger.Info("log init end")

	cookie := ""

	fetcher := grabber.NewFetcher(grabber.BrowserFetchType)

	var workerlist []*grabber.Request
	for i := 0; i <= 25; i += 25 {
		str := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", i)
		workerlist = append(workerlist, &grabber.Request{
			URL:       str,
			Cookie:    cookie,
			Timeout:   3000 * time.Millisecond,
			ParseFunc: doubantenement.ParseURL,
		})
	}

	for len(workerlist) > 0 {
		req := workerlist[0]
		workerlist = workerlist[1:]
		body, err := fetcher.Get(req)
		time.Sleep(1 * time.Second)
		if err != nil {
			logger.Error("get content error", zap.Error(err))
		}

		result := req.ParseFunc(body)
		if len(result.Requests) == 0 {
			logger.Error("get content error:regexp")
			return
		}
		for _, item := range result.Items {
			logger.Info("result", zap.String("get url", item.(string)))
		}
		workerlist = append(workerlist, result.Requests...)
	}
}
