package config

import (
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/info/service"
	"strconv"
)

func GetConfigInfo(nodeConfig *config.NodeConfig, serviceConfig *config.ServiceConfig) (map[string]string, error) {

	seid := nodeConfig.SSH.Address + ":" + strconv.Itoa(serviceConfig.Port)
	exporter, err := service.GetServiceExporter(seid, nodeConfig, serviceConfig)
	if err != nil {
		return nil, err
	}
	err = exporter.Collect(false, true)
	if err != nil {
		return nil, err
	}
	res := make(map[string]string)
	for key, value := range exporter.GetConfigMap() {
		res[key] = value
	}
	return res, nil

}
