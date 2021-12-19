package service

import (
	"github.com/xjlgod/nebula-database-diagnostic/pkg/remote"
	"strings"
)

type StorageExporter struct {
	WithMetricMap map[string]string
	WithoutMetricMap map[string]string
	ipAddress string
	port int32
	ServiceStatus string
}

func (exporter *StorageExporter) IsAlive() bool{

	status, _ := remote.GetNebulaComponentStatus(exporter.ipAddress, exporter.port);
	exporter.ServiceStatus = status[0]
	if exporter.ServiceStatus == "running" {
		return true
	}
	return false

}

func (exporter *StorageExporter) Collect() {

	if (!exporter.IsAlive()) {
		return
	}

	// 服务存活，开始通过服务接口收集信息
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

	// TODO 自动生成服务目前不能提供的信息

}

func (exporter *StorageExporter) BuildMetricMap() {
	exporter.buildWithLabels()
	exporter.buildWithoutLabels()
}

func (exporter StorageExporter) GetWithMetricMap() map[string]string{
	return exporter.WithMetricMap
}

func (exporter StorageExporter) GetWithoutMetricMap() map[string]string{
	return exporter.WithoutMetricMap
}

func (exporter *StorageExporter) Config(ipAddress string, port int32)() {
	exporter.ipAddress = ipAddress
	exporter.port = port
}


func (exporter *StorageExporter) buildWithLabels() {

	exporter.WithMetricMap = make(map[string]string)
	for _, storagedLabel := range StorageWithLabels {
		if strings.HasPrefix(storagedLabel, "num_") {
			for _, secondLabel := range ThirdLabels[:2] {
				for _, lastLabel := range LastLabels {
					exporter.buildMetricMap(storagedLabel, secondLabel, lastLabel)
				}
			}
		} else {
			for _, secondLabel := range SecondaryLabels {
				for _, thirdLabel := range ThirdLabels[2 : len(ThirdLabels)-1] {
					for _, lastLabel := range LastLabels {
						exporter.buildMetricMap(storagedLabel, secondLabel, thirdLabel, lastLabel)
					}
				}
			}
		}
	}

}

// TODO
func (exporter *StorageExporter) buildWithoutLabels() {

	exporter.WithoutMetricMap = make(map[string]string)
	for _, metaLabel := range StorageWithouLabels {
		exporter.WithoutMetricMap[metaLabel] = NotCollect
	}

}

func (exporter *StorageExporter) buildMetricMap(
	labels ...string) *StorageExporter {
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