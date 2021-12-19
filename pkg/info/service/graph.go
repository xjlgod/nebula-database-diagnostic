package service

import "strings"

type GraphExporter struct {
	metricMap map[string]string
}

func (exporter *GraphExporter) buildWithLabels() {

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

	exporter.metricMap[k] = NotCollect
	return exporter
}