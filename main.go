package main

import (
	"fmt"

	"github.com/unnull0/crawler/grabber"
	"github.com/unnull0/crawler/log"
	"go.uber.org/zap/zapcore"
)

func main() {
	core, c := log.NewFileCore("./text.log", zapcore.InfoLevel)
	defer c.Close()
	logger := log.NewLogger(core)
	logger.Info("log init end")

	// task := &grabber.Request{URL: "https://www.baidu.com"}
	task := &grabber.Request{URL: "www.baidu.com"}
	fetcher := grabber.NewFetcher(grabber.BrowserFetchType)
	content, err := fetcher.Get(task)
	if err != nil {
		logger.Error(fmt.Sprintf("get failed:%v", err))
	}
	fmt.Println("body:", string(content))
}
