package logs

import (
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/info/service"
	"strconv"
)

func GetAllLog(nodeConfig *config.NodeConfig, serviceConfig *config.ServiceConfig) error {

	seid := nodeConfig.SSH.Address + ":" + strconv.Itoa(serviceConfig.Port)
	exporter, err := service.GetServiceExporter(seid, nodeConfig, serviceConfig)
	if err != nil {
		return err
	}
	err = exporter.GetLogsInLogDir()
	return err

}
