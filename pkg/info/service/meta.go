package service

import (
	"errors"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/remote"
	"path/filepath"
	"strconv"
	"strings"
)

type MetaExporter struct {
	WithMetricMap    map[string]string
	WithoutMetricMap map[string]string
	ConfigMap        map[string]string
	ipAddress        string
	port             int
	ServiceStatus    string
	NodeConfig       *config.NodeConfig
	ServiceConfig    *config.ServiceConfig
}

func (exporter *MetaExporter) IsAlive() bool {

	status, _ := remote.GetNebulaComponentStatus(exporter.ipAddress, exporter.port)
	exporter.ServiceStatus = status[0]
	if exporter.ServiceStatus == "running" {
		return true
	}
	return false

}

func (exporter *MetaExporter) Collect(isMetrics bool, isConfig bool) error {

	if !exporter.IsAlive() {
		return errors.New("service is not alive")
	}

	// 服务存活，开始通过服务接口收集信息
	if isMetrics {
		metrics, err := remote.GetNebulaMetrics(exporter.ipAddress, exporter.port)
		if err == nil {
			matches := convertToMap(metrics)
			for metric, value := range matches {
				if err != nil {
					continue
				}
				if _, ok := exporter.WithMetricMap[metric]; ok {
					exporter.WithMetricMap[metric] = value
				}
			}
		}
	}

	// TODO 自动生成服务目前不能提供的信息

	// 获取服务配置信息
	if isConfig {
		configs, err := remote.GetNebulaConfigs(exporter.ipAddress, exporter.port)
		if err == nil {
			matches := convertToMap(configs)
			for config, value := range matches {
				if _, ok := exporter.ConfigMap[config]; ok {
					exporter.ConfigMap[config] = value
				}
			}
		} else {
			return err
		}
	}

	return nil

}

func (exporter *MetaExporter) BuildAllMap() {
	exporter.buildWithLabels()
	exporter.buildWithoutLabels()
	exporter.buildConfigLabels()
}

func (exporter MetaExporter) GetWithMetricMap() map[string]string {
	return exporter.WithMetricMap
}

func (exporter MetaExporter) GetWithoutMetricMap() map[string]string {
	return exporter.WithoutMetricMap
}

func (exporter MetaExporter) GetConfigMap() map[string]string {
	return exporter.ConfigMap
}

func (exporter *MetaExporter) Config(nodeConfig *config.NodeConfig, serviceConfig *config.ServiceConfig) {
	exporter.NodeConfig = nodeConfig
	exporter.ServiceConfig = serviceConfig
	exporter.ipAddress = nodeConfig.SSH.Address
	exporter.port = serviceConfig.HTTPPort
}

func (exporter *MetaExporter) GetLogsInLogDir() error {

	logDir, ok := exporter.ConfigMap["log_dir"]

	if (logDir == NotCollect || exporter.ServiceConfig.RuntimeDir == "") {
		return errors.New("logdir is nut exist")
	}
	if !ok {
		return errors.New("logdir is nut exist")
	}
	if (!strings.HasPrefix(logDir, "/")) {
		logDir = exporter.ServiceConfig.RuntimeDir + "/" + logDir
	}
	newDir := exporter.NodeConfig.SSH.Address + "-" + strconv.Itoa(exporter.ServiceConfig.Port)
	localDir := filepath.Join(exporter.NodeConfig.Output.DirPath , newDir)
	err := remote.GetFilesInRemoteDir(exporter.NodeConfig.SSH.Username, exporter.NodeConfig.SSH, logDir, localDir)
	return err

}

func (exporter *MetaExporter) buildConfigLabels() {
	exporter.ConfigMap = make(map[string]string)
	for _, metaLabel := range MetaConfigLabels {
		exporter.ConfigMap[metaLabel] = NotCollect
	}
}

func (exporter *MetaExporter) buildWithLabels() {

	exporter.WithMetricMap = make(map[string]string)
	for _, metadLabel := range MetadWithLabels {
		if strings.HasPrefix(metadLabel, "num_") {
			for _, secondLabel := range ThirdLabels[:2] {
				for _, lastLabel := range LastLabels {
					exporter.buildMetricMap(metadLabel, secondLabel, lastLabel)
				}
			}
		} else {
			for _, secondLabel := range SecondaryLabels {
				for _, thirdLabel := range ThirdLabels[2 : len(ThirdLabels)-1] {
					for _, lastLabel := range LastLabels {
						exporter.buildMetricMap(metadLabel, secondLabel, thirdLabel, lastLabel)
					}
				}
			}
		}
	}

}

func (exporter *MetaExporter) buildWithoutLabels() {

	exporter.WithoutMetricMap = make(map[string]string)
	for _, metaLabel := range MetaWithouLabels {
		exporter.WithoutMetricMap[metaLabel] = NotCollect
	}

}

func (exporter *MetaExporter) buildMetricMap(
	labels ...string) *MetaExporter {
	if len(labels) == 0 {
		return exporter
	}

	var k string
	var metricName string

	last := len(labels) - 2
	if last <= 0 {
		return exporter
	}

	for i, label := range labels {
		if i == 0 {
			metricName = label
			k = label
		} else {
			metricName = metricName + "_" + label
			if i < last {
				k = k + "_" + label
			}
		}
	}

	for _, label := range labels[last:] {
		k = k + "." + label
	}

	exporter.WithMetricMap[k] = NotCollect
	return exporter
}
