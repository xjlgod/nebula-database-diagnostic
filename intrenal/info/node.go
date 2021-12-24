package info

import (
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/info/physical"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/logger"
)

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
