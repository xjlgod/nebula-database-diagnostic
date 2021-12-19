package service

import "strings"

type (
	ComponentType string

	ServiceExporter interface {
		// IsAlive 判断需要收集的服务是否存活，并更新服务状态
		IsAlive() bool
		// Collect 收集所有服务指标信息
		Collect()
		// BuildMetricMap 进行所需服务指标信息的map构建
		BuildMetricMap()
		// GetWithMetricMap 返回能够收集到的服务指标信息 key为指标名，value为值
		GetWithMetricMap() map[string]string
		// GetWithoutMetricMap 返回目前服务不能提供，自动生成的服务指标信息 key为指标名，value为值
		GetWithoutMetricMap() map[string]string
		// Config TODO 配置Exporter，目前只配置需收集的服务的ip与端口
		Config(ipAddress string, port int32)
	}
)

const (
	GraphdComponent   ComponentType = "graphd"
	MetadComponent    ComponentType = "metad"
	StoragedComponent ComponentType = "storaged"

	NotCollect string = "wait for collect"
)

var (

	// WithLabels为服务目前能够提供的信息
	GraphWithLabels = []string{
		"slow_query",
		"query",
		"num_query_errors",
		"num_queries",
		"num_slow_queries",
	}

	// WithouLabels为服务目前不能够提供的信息
	GraphWithouLabels = []string{
		"version",
		"pid",
		"port",
		"deploy_dir",
		"log",
		"config",
		"runtime_dir",
	}

	MetadWithLabels = []string{
		"heartbeat",
		"num_heartbeats",
	}

	MetaWithouLabels = []string{
		"version",
		"pid",
		"port",
		"deploy_dir",
		"log",
		"config",
		"runtime_dir",
	}

	StorageWithLabels = []string{
		"num_add_vertices",
		"num_add_edges",
		"add_vertices_latency_us",
		"add_edges_latency_us",
	}

	StorageWithouLabels = []string{
		"version",
		"pid",
		"port",
		"deploy_dir",
		"log",
		"config",
		"runtime_dir",
	}

	SecondaryLabels = []string{
		"latency_us",
	}

	ThirdLabels = []string{
		"rate", "sum", "avg", "p75", "p95", "p99", "p999",
	}

	LastLabels = []string{
		"5", "60", "600", "3600",
	}

)

func convertToMap(metrics []string) map[string]string {
	matches := make(map[string]string)
	// match metric format
	// value=0
	// name=num_queries.rate.60
	if strings.Contains(metrics[0], "value=") {
		for t := 0; ; t += 2 {
			if t+1 >= len(metrics) {
				break
			}
			ok, value := getRValue(metrics[t])
			if !ok {
				continue
			}
			ok, metric := getRValue(metrics[t+1])
			if !ok {
				continue
			}
			if value != "" && metric != "" {
				matches[metric] = value
			}
		}
		return matches
	}

	// match metric format
	// slow_query_latency_us.p95.5=0
	for _, metric := range metrics {
		s := strings.Split(metric, "=")
		if len(s) != 2 {
			continue
		}
		matches[s[0]] = s[1]
	}
	return matches
}

func getRValue(metric string) (bool, string) {
	seps := strings.Split(metric, "=")
	if len(seps) != 2 {
		return false, ""
	}
	return true, seps[1]
}
