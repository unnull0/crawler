package main

import (
	"fmt"

	"github.com/unnull0/crawler/grabber"
)

func main() {
	task := &grabber.Request{URL: "https://www.baidu.com"}
	fetcher := grabber.NewFetcher(grabber.BrowserFetchType)
	content, err := fetcher.Get(task)
	if err != nil {
		fmt.Printf("get failed:%v\n", err)
	}
	fmt.Println("body:", string(content))
}
