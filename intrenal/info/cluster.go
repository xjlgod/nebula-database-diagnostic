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
		// the conf has been verified, so don't need to handle error
		d, _ := time.ParseDuration(node.Duration)
		if node.Output.DirPath == "" {
			if d == 0 {
				run(node, logger.CMDLogger)
			} else {
				runWithDuration(node, d, logger.CMDLogger)
			}
		} else {
			if d == 0 {
				run(node, logger.FileLogger)
			} else {
				runWithDuration(node, d, logger.FileLogger)
			}
		}
	}
}

func run(conf *config.NodeConfig, defaultLogger logger.Logger) {
	for _, option := range conf.Infos {
		fetchInfo(conf, option, defaultLogger)
	}
}

func runWithDuration(conf *config.NodeConfig, duration time.Duration, defaultLogger logger.Logger) {

}

func fetchInfo(conf *config.NodeConfig, option config.InfoOption, defaultLogger logger.Logger) {
	// ...
	phyInfo, err := fetchPhyInfo(option, conf.SSH)
	if err != nil {
		defaultLogger.Errorf("fetch phy info failed: %s", err.Error())
	} else {
		//defaultLogger.Info(phyInfo.String())
		defaultLogger.Info(phyInfo)
	}
	// ...
}

func fetchPhyInfo(option config.InfoOption, sshConfig config.SSHConfig) (*physical.PhyInfo, error) {
	if option == config.AllInfo || option == config.Physical {
		return physical.GetPhyInfo(sshConfig)
	}
	return nil, nil
}
