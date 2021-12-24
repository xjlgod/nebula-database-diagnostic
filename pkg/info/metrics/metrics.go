package metrics

import (
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/info/service"
	"strconv"
)


func GetMetricsInfo(nodeConfig *config.NodeConfig, serviceConfig *config.ServiceConfig) (map[string]string, error) {

	seid := nodeConfig.SSH.Address + ":" + strconv.Itoa(serviceConfig.Port)
	exporter, err := service.GetServiceExporter(seid, nodeConfig, serviceConfig)
	if err != nil {
		return nil, err
	}
	err = exporter.Collect(true, false)
	if err != nil {
		return nil, err
	}
	res := make(map[string]string)
	for key, value := range exporter.GetWithMetricMap() {
		res[key] = value
	}
	for key, value := range exporter.GetWithoutMetricMap() {
		res[key] = value
	}
	return res, nil

}
