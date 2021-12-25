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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for name, node := range conf.Nodes {
		node := node
		name := name
		go func() {
			var _logger logger.Logger
			if node.Output.LogToFile {
				_logger = logger.GetFileLogger(name, node.Output)
			} else {
				_logger = logger.GetCmdLogger(name)
			}
			// the conf has been verified, so don't need to handle error
			d, _ := time.ParseDuration(node.Duration)
			if node.Duration == "-1" {
				runWithInfinity(node, _logger)
			} else if d == 0 {
				run(node, _logger)
			} else {
				runWithDuration(node, _logger)
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

func run(conf *config.NodeConfig, defaultLogger logger.Logger) {
	for _, option := range conf.Infos {
		//fetchInfo(conf, option, defaultLogger)
		fetchAndSaveInfo(conf, option, defaultLogger)
	}
}

func runWithInfinity(conf *config.NodeConfig, defaultLogger logger.Logger) {
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

func runWithDuration(conf *config.NodeConfig, defaultLogger logger.Logger) {
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
