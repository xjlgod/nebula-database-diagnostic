package service

import (
	"errors"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/remote"
	"path/filepath"
	"strconv"
	"strings"
)

type GraphExporter struct {
	WithMetricMap    map[string]string
	WithoutMetricMap map[string]string
	ConfigMap        map[string]string
	ipAddress        string
	port             int
	ServiceStatus    string
	InfoConfig       *config.InfoConfig
	ServiceConfig    *config.ServiceConfig
}

func (exporter *GraphExporter) IsAlive() bool {

	status, _ := remote.GetNebulaComponentStatus(exporter.ipAddress, exporter.port)
	exporter.ServiceStatus = status[0]
	if exporter.ServiceStatus == "running" {
		return true
	}
	return false

}

func (exporter *GraphExporter) Collect(isMetrics bool, isConfig bool) error {

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

func (exporter *GraphExporter) BuildAllMap() {
	exporter.buildWithLabels()
	exporter.buildWithoutLabels()
	exporter.buildConfigLabels()
}

func (exporter GraphExporter) GetWithMetricMap() map[string]string {
	return exporter.WithMetricMap
}

func (exporter GraphExporter) GetWithoutMetricMap() map[string]string {
	return exporter.WithoutMetricMap
}

func (exporter GraphExporter) GetConfigMap() map[string]string {
	return exporter.ConfigMap
}

func (exporter *GraphExporter) Config(conf *config.InfoConfig, serviceConfig *config.ServiceConfig) {
	exporter.InfoConfig = conf
	exporter.ServiceConfig = serviceConfig
	exporter.ipAddress = conf.Node.SSH.Address
	exporter.port = serviceConfig.HTTPPort
}

func (exporter *GraphExporter) GetLogsInLogDir() error {

	logDir, ok := exporter.ConfigMap["log_dir"]
	newlogDir := strings.Replace(logDir, "\"", "", 10)

	if logDir == NotCollect || exporter.ServiceConfig.RuntimeDir == "" {
		return errors.New("logdir is nut exist")
	}
	if !ok {
		return errors.New("logdir is nut exist")
	}
	if !strings.HasPrefix(logDir, "/") {
		newlogDir = exporter.ServiceConfig.RuntimeDir + "/" + newlogDir
	}
	newDir := exporter.InfoConfig.Node.SSH.Address + "-" + strconv.Itoa(exporter.ServiceConfig.Port)
	localDir := filepath.Join(exporter.InfoConfig.Output.DirPath, newDir)
	err := remote.GetFilesInRemoteDir(exporter.InfoConfig.Node.SSH.Username, exporter.InfoConfig.Node.SSH, newlogDir, localDir)
	return err

}

func (exporter *GraphExporter) buildConfigLabels() {
	exporter.ConfigMap = make(map[string]string)
	for _, metaLabel := range GraphConfigLabels {
		exporter.ConfigMap[metaLabel] = NotCollect
	}
}

func (exporter *GraphExporter) buildWithLabels() {

	exporter.WithMetricMap = make(map[string]string)
	for _, graphdLabel := range GraphWithLabels {
		if strings.HasPrefix(graphdLabel, "num_") {
			for _, secondLabel := range ThirdLabels[:2] {
				for _, lastLabel := range LastLabels {
					exporter.buildMetricMap(graphdLabel, secondLabel, lastLabel)
				}
			}
		} else {
			for _, secondLabel := range SecondaryLabels {
				for _, thirdLabel := range ThirdLabels[2:] {
					for _, lastLabel := range LastLabels {
						exporter.buildMetricMap(graphdLabel, secondLabel, thirdLabel, lastLabel)
					}
				}
			}
		}
	}

}

// TODO
func (exporter *GraphExporter) buildWithoutLabels() {

	exporter.WithoutMetricMap = make(map[string]string)
	for _, metaLabel := range GraphWithouLabels {
		exporter.WithoutMetricMap[metaLabel] = NotCollect
	}

}

func (exporter *GraphExporter) buildMetricMap(
	labels ...string) *GraphExporter {
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
