package info

import (
	"context"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/logger"
	"time"
)

type (
	Cluster struct {
		NodeNum int
	}
)

func Run(conf *config.Config) {
	// TODO fix the cancel bugs
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for name, info := range conf.Infos {
		info := info
		name := name
		go func() {
			var _logger logger.Logger
			if info.Output.LogToFile {
				_logger = logger.GetFileLogger(name, info.Output)
			} else {
				_logger = logger.GetCmdLogger(name)
			}
			// the conf has been verified, so don't need to handle error
			d, _ := time.ParseDuration(info.Duration)
			if info.Duration == "-1" {
				runWithInfinity(info, _logger)
			} else if d == 0 {
				run(info, _logger)
				cancel() // TODO temp code
			} else {
				runWithDuration(info, _logger)
			}
		}()
	}

	for {
		select {
		case <-ctx.Done():
			return
		}
	}

}

func run(conf *config.InfoConfig, defaultLogger logger.Logger) {
	for _, option := range conf.Options {
		//fetchInfo(conf, option, defaultLogger)
		fetchAndSaveInfo(conf, option, defaultLogger)
	}
}

func runWithInfinity(conf *config.InfoConfig, defaultLogger logger.Logger) {
	p, _ := time.ParseDuration(conf.Period)
	ticker := time.NewTicker(p)
	for {
		select {
		case <-ticker.C:
			run(conf, defaultLogger)
		default:

		}
	}
}

func runWithDuration(conf *config.InfoConfig, defaultLogger logger.Logger) {
	p, _ := time.ParseDuration(conf.Period)
	ticker := time.NewTicker(p)
	ch := make(chan bool)
	go func(ticker *time.Ticker) {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				run(conf, defaultLogger)
			case stop := <-ch:
				if stop {
					return
				}
			default:

			}
		}
	}(ticker)

	d, _ := time.ParseDuration(conf.Duration)
	time.Sleep(d)
	ch <- true
	close(ch)
}
