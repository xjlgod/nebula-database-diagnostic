package info

import (
	"encoding/json"
	"fmt"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	configinfo "github.com/xjlgod/nebula-database-diagnostic/pkg/info/config"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/info/logs"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/info/metrics"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/info/physical"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/info/service"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/logger"
	"os"
	"path/filepath"
	"strings"
)

var NowAllInfo AllInfo

type AllInfo struct {
	Time        string                       `json:"time"`
	PhyInfo     *physical.PhyInfo            `json:"phy_info"`
	MetricsInfo []*service.ServiceMetricInfo `json:"metrics_info,omitempty"`
	ConfigInfo  []*service.ServiceConfigInfo `json:"config_info,omitempty"`
}

func fetchInfo(conf *config.InfoConfig, option config.InfoOption, defaultLogger logger.Logger) (*physical.PhyInfo,
	[]*service.ServiceMetricInfo, []*service.ServiceConfigInfo) {
	phyInfo, err := fetchPhyInfo(option, conf.Node.SSH)
	if err != nil {
		defaultLogger.Errorf("fetch phy info failed: %s\n", err.Error())
	} else {
		// defaultLogger.Info(phyInfo.String())
		if phyInfo != nil {
			defaultLogger.Infof("%s physical info: %+v\n", conf.Node.Host, phyInfo)
		}
	}

	// fetch all services metrics info
	servicesMetricsInfo, err := fetchMetricsInfo(conf, option, defaultLogger)
	if err != nil {
		defaultLogger.Errorf("fetch services metrics failed: %s\n", err.Error())
	} else {
		// defaultLogger.Info(phyInfo.String())
		if servicesMetricsInfo != nil {
			for i := range servicesMetricsInfo {
				defaultLogger.Infof("%s metrics info: %+v\n", conf.Node.Host, servicesMetricsInfo[i])
			}
		}
	}

	// fetch all services config info
	serviceConfigsInfo, err := fetchConfigsInfo(conf, option, defaultLogger)
	if err != nil {
		defaultLogger.Errorf("fetch services metrics failed: %s", err.Error())
	} else {
		// defaultLogger.Info(phyInfo.String())
		if serviceConfigsInfo != nil {
			for i := range serviceConfigsInfo {
				defaultLogger.Infof("%s config info: %+v\n", conf.Node.Host, serviceConfigsInfo[i])
			}
		}
	}

	// pack all services log
	err = packageLogs(conf, option, defaultLogger)
	defaultLogger.Info("packaging service logs...\n")
	if err != nil {
		defaultLogger.Errorf("service package: failed, %s\n", err.Error())
	} else {
		defaultLogger.Info(conf.Node.Host, " service package: success!\n")
	}

	return phyInfo, servicesMetricsInfo, serviceConfigsInfo
}

func packageLogs(conf *config.InfoConfig, option config.InfoOption, defaultLogger logger.Logger) error {
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

func fetchMetricsInfo(conf *config.InfoConfig, option config.InfoOption, defaultLogger logger.Logger) ([]*service.ServiceMetricInfo, error) {
	if option == config.AllInfo || option == config.Metrics {
		servicesConfig := conf.Services
		servicesMetricsInfo := make([]*service.ServiceMetricInfo, 0, len(servicesConfig))
		for key, serviceConfig := range servicesConfig {
			metrics, err := metrics.GetMetricsInfo(conf, &serviceConfig)
			serviceMetricsInfo := &service.ServiceMetricInfo{
				Name:    key,
				Metrics: metrics,
				Type:    serviceConfig.Type,
			}
			if err != nil {
				defaultLogger.Errorf("fetch metrics info failed: %s, stop fetch services metrics!\n", err.Error())
				return servicesMetricsInfo, err
			}
			servicesMetricsInfo = append(servicesMetricsInfo, serviceMetricsInfo)
		}
		return servicesMetricsInfo, nil
	}
	return nil, nil

}

func fetchConfigsInfo(conf *config.InfoConfig, option config.InfoOption, defaultLogger logger.Logger) ([]*service.ServiceConfigInfo, error) {
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
				defaultLogger.Errorf("fetch configs info failed: %s, stop fetch services configs!\n", err.Error())
				return serviceConfigInfo, err
			}
			serviceConfigInfo = append(serviceConfigInfo, serviceMetricsInfo)
		}
		return serviceConfigInfo, nil
	}
	return nil, nil
}

func fetchAndSaveInfo(conf *config.InfoConfig, option config.InfoOption, defaultLogger logger.Logger) {

	phyInfo, metricsInfo, configInfo := fetchInfo(conf, option, defaultLogger)
	allInfo := &AllInfo{
		PhyInfo:     phyInfo,
		MetricsInfo: metricsInfo,
		ConfigInfo:  configInfo,
	}
	marshal, err := json.Marshal(allInfo)
	if err != nil {
		defaultLogger.Errorf("save json data failed: %s\n", err.Error())
	}

	dir := filepath.Join(conf.Output.DirPath, conf.Node.Host.Address)
	p, _ := filepath.Abs(dir)
	_, err = os.Stat(p)
	if os.IsNotExist(err) {
		os.Mkdir(p, os.ModePerm)
	}

	// node0_160.data
	// node0_161.data
	// node0/ 160.data
	//

	loggerFileName := defaultLogger.Filename()
	slice := strings.Split(loggerFileName, "_")
	time := slice[1]

	filename := fmt.Sprintf("%s%s", time, ".data")
	filePath := filepath.Join(p, filename)
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			defaultLogger.Fatal(err)
		}
		_, err = file.Write(marshal)
		_, err = file.Write([]byte("\n"))
		if err != nil {
			defaultLogger.Errorf("save json data failed: %s\n", err.Error())
		}
	} else {
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			defaultLogger.Fatal(err)
		}
		_, err = file.Write(marshal)
		_, err = file.Write([]byte("\n"))
		if err != nil {
			defaultLogger.Errorf("save json data failed: %s\n", err.Error())
		}
	}

}
