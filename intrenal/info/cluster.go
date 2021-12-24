package info

import (
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/info/physical"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/logger"
	"time"
)

type (
	Cluster struct {
		NodeNum int
	}
)

func Run(conf *config.Config) {
	for _, node := range conf.Nodes {
		var _logger logger.Logger
		logger.InitCmdLogger()
		logger.InitFileLogger(node.Output)
		if node.Output.LogToFile {
			_logger = logger.FileLogger
		} else {
			_logger = logger.CMDLogger
		}
		// the conf has been verified, so don't need to handle error
		d, _ := time.ParseDuration(node.Duration)

		if d == 0 {
			go run(node, _logger)
		} else if d < 0 {
			go runWithInfinity(node, _logger)
		} else {
			go runWithDuration(node, _logger)
		}
	}
}

func run(conf *config.NodeConfig, defaultLogger logger.Logger) {
	for _, option := range conf.Infos {
		fetchInfo(conf, option, defaultLogger)
	}
}

func runWithInfinity(conf *config.NodeConfig, defaultLogger logger.Logger) {
	go func() {
		p, _ := time.ParseDuration(conf.Period)
		for {
			time.AfterFunc(p, func() {
				run(conf, defaultLogger)
			})
		}
	}()
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

func fetchInfo(conf *config.NodeConfig, option config.InfoOption, defaultLogger logger.Logger) {
	phyInfo, err := fetchPhyInfo(option, conf.SSH)
	if err != nil {
		defaultLogger.Errorf("fetch phy info failed: %s", err.Error())
	} else {
		//defaultLogger.Info(phyInfo.String())
		defaultLogger.Info(phyInfo)
	}
	// ... fetch more info here
}

func fetchPhyInfo(option config.InfoOption, sshConfig config.SSHConfig) (*physical.PhyInfo, error) {
	if option == config.AllInfo || option == config.Physical {
		return physical.GetPhyInfo(sshConfig)
	}
	return nil, nil
}
