package service

import "strings"

type StorageExporter struct {
	metricMap map[string]string
}

func (exporter *StorageExporter) buildWithLabels() {

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

	exporter.metricMap[k] = NotCollect
	return exporter
}