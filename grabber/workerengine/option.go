package workerengine

import (
	"github.com/unnull0/crawler/grabber"
	"go.uber.org/zap"
)

type options struct {
	WorkCount int
	Fetcher   grabber.Fetcher
	Logger    *zap.Logger
	Seeds     []*grabber.Task
	scheduler Scheduler
}

var defaultOptions = options{
	Logger: zap.NewNop(),
}

type Option func(opts *options)

func WithWorkCount(workCount int) Option {
	return func(opts *options) {
		opts.WorkCount = workCount
	}
}

func WithFetcher(fetcher grabber.Fetcher) Option {
	return func(opts *options) {
		opts.Fetcher = fetcher
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.Logger = logger
	}
}

func WithSeeds(task []*grabber.Task) Option {
	return func(opts *options) {
		opts.Seeds = task
	}
}

func WithScheduler(scheduler Scheduler) Option {
	return func(opts *options) {
		opts.scheduler = scheduler
	}
}
