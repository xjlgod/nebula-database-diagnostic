package info

import (
	"encoding/json"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	configinfo "github.com/xjlgod/nebula-database-diagnostic/pkg/info/config"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/info/logs"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/info/metrics"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/info/physical"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/info/service"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/logger"
	"os"
	"path/filepath"
)

var NowAllInfo AllInfo

type AllInfo struct {
	Time string `json:"time"`
	PhyInfo     *physical.PhyInfo            `json:"phy_info"`
	MetricsInfo []*service.ServiceMetricInfo `json:"metrics_info,omitempty"`
	ConfigInfo  []*service.ServiceConfigInfo `json:"config_info,omitempty"`
}

func fetchInfo(conf *config.NodeConfig, option config.InfoOption, defaultLogger logger.Logger) (*physical.PhyInfo,
	[]*service.ServiceMetricInfo, []*service.ServiceConfigInfo) {
	phyInfo, err := fetchPhyInfo(option, conf.SSH)
	if err != nil {
		defaultLogger.Errorf("fetch phy info failed: %s", err.Error())
	} else {
		// defaultLogger.Info(phyInfo.String())
		if phyInfo != nil {
			defaultLogger.Info(conf.Host, ": ", phyInfo)
		}
	}

	// fetch all services metrics info
	servicesMetricsInfo, err := fetchMetricsInfo(conf, option, defaultLogger)
	if err != nil {
		defaultLogger.Errorf("fetch services metrics failed: %s", err.Error())
	} else {
		// defaultLogger.Info(phyInfo.String())
		if servicesMetricsInfo != nil {
			defaultLogger.Info(conf.Host, ": ", servicesMetricsInfo)
		}
	}

	// fetch all services config info
	serviceConfigsInfo, err := fetchConfigsInfo(conf, option, defaultLogger)
	if err != nil {
		defaultLogger.Errorf("fetch services metrics failed: %s", err.Error())
	} else {
		// defaultLogger.Info(phyInfo.String())
		if serviceConfigsInfo != nil {
			defaultLogger.Info(conf.Host, ": ", serviceConfigsInfo)
		}
	}

	// pack all services log
	err = packageLogs(conf, option, defaultLogger)
	if err != nil {
		defaultLogger.Errorf("package logs failed: %s", err.Error())
	} else {
		defaultLogger.Info(conf.Host, ": ", "package logs success!")
	}

	return phyInfo, servicesMetricsInfo, serviceConfigsInfo
}

func packageLogs(conf *config.NodeConfig, option config.InfoOption, defaultLogger logger.Logger) error {
	if option == config.AllInfo || option == config.Metrics {
		servicesConfig := conf.Services
		for _, serviceConfig := range servicesConfig {
			err := logs.GetAllLog(conf, &serviceConfig)
			if err != nil {
				defaultLogger.Errorf("package logs failed: %s, stop package logs!", err.Error())
				return err
			}
		}
	}
	return nil
}


func fetchPhyInfo(option config.InfoOption, sshConfig config.SSHConfig) (*physical.PhyInfo, error) {
	if option == config.AllInfo || option == config.Physical {
		return physical.GetPhyInfo(sshConfig)
	}
	return nil, nil
}

func fetchMetricsInfo(conf *config.NodeConfig, option config.InfoOption, defaultLogger logger.Logger) ([]*service.ServiceMetricInfo, error) {
	if option == config.AllInfo || option == config.Metrics {
		servicesConfig := conf.Services
		servicesMetricsInfo := make([]*service.ServiceMetricInfo, 0, len(servicesConfig))
		for key, serviceConfig := range servicesConfig {
			metrics, err := metrics.GetMetricsInfo(conf, &serviceConfig)
			serviceMetricsInfo := &service.ServiceMetricInfo{
				Name:    key,
				Metrics: metrics,
			}
			if err != nil {
				defaultLogger.Errorf("fetch metrics info failed: %s, stop fetch services metrics!", err.Error())
				return servicesMetricsInfo, err
			}
			servicesMetricsInfo = append(servicesMetricsInfo, serviceMetricsInfo)
		}
		return servicesMetricsInfo, nil
	}
	return nil, nil

}

func fetchConfigsInfo(conf *config.NodeConfig, option config.InfoOption, defaultLogger logger.Logger) ([]*service.ServiceConfigInfo, error) {
	// config belongs to Metrics
	if option == config.AllInfo || option == config.Metrics {
		servicesConfig := conf.Services
		serviceConfigInfo := make([]*service.ServiceConfigInfo, 0, len(servicesConfig))
		for key, serviceConfig := range servicesConfig {
			configs, err := configinfo.GetConfigInfo(conf, &serviceConfig)
			serviceMetricsInfo := &service.ServiceConfigInfo{
				Name:    key,
				Configs: configs,
			}
			if err != nil {
				defaultLogger.Errorf("fetch configs info failed: %s, stop fetch services metrics!", err.Error())
				return serviceConfigInfo, err
			}
			serviceConfigInfo = append(serviceConfigInfo, serviceMetricsInfo)
		}
		return serviceConfigInfo, nil
	}
	return nil, nil
}

func fetchAndSaveInfo(conf *config.NodeConfig, option config.InfoOption, defaultLogger logger.Logger) {

	phyInfo, metricsInfo, configInfo := fetchInfo(conf, option, defaultLogger)
	allInfo := &AllInfo{
		PhyInfo:     phyInfo,
		MetricsInfo: metricsInfo,
		ConfigInfo:  configInfo,
	}
	marshal, err := json.Marshal(allInfo)
	if err != nil {
		defaultLogger.Errorf("save json data fail: %s", err.Error())
	}
	p, _ := filepath.Abs(conf.Output.DirPath)
	_, err = os.Stat(p)
	if os.IsNotExist(err) {
		os.Mkdir(p, os.ModePerm)
	}


	filename := defaultLogger.Filename()
	filePath := filepath.Join(p, filename+".data")
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			defaultLogger.Fatal(err)
		}
		_, err = file.Write(marshal)
		if err != nil {
			defaultLogger.Errorf("save json data fail: %s", err.Error())
		}
	} else {
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			defaultLogger.Fatal(err)
		}
		_, err = file.Write([]byte("\n"))
		_, err = file.Write(marshal)
		if err != nil {
			defaultLogger.Errorf("save json data fail: %s", err.Error())
		}
	}



}
