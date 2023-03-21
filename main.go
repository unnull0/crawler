package main

import (
	"time"

	"github.com/unnull0/crawler/collector/sqlstorage"
	"github.com/unnull0/crawler/grabber"
	"github.com/unnull0/crawler/grabber/workerengine"
	"github.com/unnull0/crawler/limiter"
	"github.com/unnull0/crawler/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/time/rate"
)

func main() {
	core, c := log.NewFileCore("./text.log", zapcore.InfoLevel)
	defer c.Close()
	logger := log.NewLogger(core)
	logger.Info("log init end")

	var fetcher grabber.Fetcher = &grabber.BrowserFetch{
		Timeout: 10 * time.Second,
	}

	storage, err := sqlstorage.New(
		sqlstorage.WithSqlUrl("root:9516@tcp(127.0.0.1:3306)/crawler?charset=utf8"),
		sqlstorage.WithLogger(logger),
		sqlstorage.WithBatchCount(5),
	)
	if err != nil {
		logger.Error("creat storage failed", zap.Error(err))
		return
	}

	secondLimit := rate.NewLimiter(limiter.Per(1, 5*time.Second), 1)
	minuteLimit := rate.NewLimiter(limiter.Per(10, 1*time.Minute), 10)
	multiLimiter := limiter.MultiLimiter(secondLimit, minuteLimit)

	seeds := make([]*grabber.Task, 0, 1000)
	seeds = append(seeds, &grabber.Task{
		Name:    "douban_book_list",
		Fetcher: fetcher,
		Storage: storage,
		Limit:   multiLimiter,
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
