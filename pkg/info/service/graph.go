package service

import (
	"github.com/xjlgod/nebula-database-diagnostic/pkg/remote"
	"strings"
)

type GraphExporter struct {
	WithMetricMap map[string]string
	WithoutMetricMap map[string]string
	ipAddress string
	port int32
	ServiceStatus string
}

func (exporter *GraphExporter) IsAlive() bool{

	status, _ := remote.GetNebulaComponentStatus(exporter.ipAddress, exporter.port);
	exporter.ServiceStatus = status[0]
	if exporter.ServiceStatus == "running" {
		return true
	}
	return false

}

func (exporter *GraphExporter) Collect() {

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

func (exporter *GraphExporter) BuildMetricMap() {
	exporter.buildWithLabels()
	exporter.buildWithoutLabels()
}

func (exporter GraphExporter) GetWithMetricMap() map[string]string{
	return exporter.WithMetricMap
}

func (exporter GraphExporter) GetWithoutMetricMap() map[string]string{
	return exporter.WithoutMetricMap
}

func (exporter *GraphExporter) Config(ipAddress string, port int32)() {
	exporter.ipAddress = ipAddress
	exporter.port = port
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