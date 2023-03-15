package main

import (
	"time"

	"github.com/unnull0/crawler/grabber"
	"github.com/unnull0/crawler/grabber/workerengine"
	"github.com/unnull0/crawler/log"
	"go.uber.org/zap/zapcore"
)

func main() {
	core, c := log.NewFileCore("./text.log", zapcore.InfoLevel)
	defer c.Close()
	logger := log.NewLogger(core)
	logger.Info("log init end")

	var fetcher grabber.Fetcher = &grabber.BrowserFetch{
		Timeout: 3000 * time.Millisecond,
	}
	seeds := make([]*grabber.Task, 0, 1000)
	seeds = append(seeds, &grabber.Task{
		Name:    "find_douban_sun_rom",
		Fetcher: fetcher,
	})

	e := workerengine.NewEngine(
		workerengine.WithFetcher(fetcher),
		workerengine.WithLogger(logger),
		workerengine.WithScheduler(workerengine.NewSchedule()),
		workerengine.WithWorkCount(5),
		workerengine.WithSeeds(seeds),
	)
	e.Run()
}
