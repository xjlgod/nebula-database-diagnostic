package service

import (
	"errors"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	"strconv"
	"strings"
	"sync"
)

type (
	ServiceExporter interface {
		// IsAlive 判断需要收集的服务是否存活，并更新服务状态
		IsAlive() bool
		// Collect 收集所有服务指标信息
		Collect(isMetrics bool, isConfig bool) error
		// BuildAllMap 进行所需服务所有信息的map构建
		BuildAllMap()
		// GetWithMetricMap 返回能够收集到的服务指标信息 key为指标名，value为值
		GetWithMetricMap() map[string]string
		// GetWithoutMetricMap 返回目前服务不能提供，自动生成的服务指标信息 key为指标名，value为值
		GetWithoutMetricMap() map[string]string
		// GetConfigMap 返回能够收集到的服务配置信息
		GetConfigMap() map[string]string
		// Config 配置config
		Config(nodeConfig *config.NodeConfig, serviceConfig *config.ServiceConfig)
		// GetLogsInLogDir 通过指定路径获取日志
		GetLogsInLogDir() error
	}
)

type (
	ServiceMetricInfo struct {
		Name    string            `json:"name,omitempty"`
		Metrics map[string]string `json:"metrics,omitempty"`
		Type config.ComponentType `json:"type"`
	}
	ServiceConfigInfo struct {
		Name    string            `json:"name,omitempty"`
		Configs map[string]string `json:"configs,omitempty"`
	}
)

const (
	NotCollect    string = "wait for collect"
	LOCAL_LOG_DIR        = "data/logs"
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
		"logs",
		"config",
		"runtime_dir",
	}

	GraphConfigLabels = []string{
		"log_dir",
		"pid_file",
		"port",
		"ws_h2_port",
		"ws_http_port",
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
		"logs",
		"config",
		"runtime_dir",
	}

	MetaConfigLabels = []string{
		"log_dir",
		"pid_file",
		"port",
		"ws_h2_port",
		"ws_http_port",
	}

	StorageWithLabels = []string{
		"num_add_vertices",
		"num_add_edges",
		"add_vertices",
		"add_edges",
	}

	StorageWithouLabels = []string{
		"version",
		"pid",
		"port",
		"deploy_dir",
		"logs",
		"config",
		"runtime_dir",
	}

	StorageConfigLabels = []string{
		"log_dir",
		"pid_file",
		"port",
		"ws_h2_port",
		"ws_http_port",
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

var exporters = make(map[string]ServiceExporter)
var mux sync.RWMutex

func GetServiceExporter(seid string, nodeConfig *config.NodeConfig, serviceConfig *config.ServiceConfig) (ServiceExporter, error) {
	mux.Lock()
	if _, ok := exporters[seid]; !ok {
		e, err := newServiceExporter(nodeConfig, serviceConfig)
		if err != nil {
			return nil, err
		}

		exporters[seid] = e
	}
	mux.Unlock()

	mux.RLock()
	defer mux.RUnlock()

	return exporters[seid], nil
}

func newServiceExporter(nodeConfig *config.NodeConfig, serviceConfig *config.ServiceConfig) (ServiceExporter, error) {

	serviceType := serviceConfig.Type
	switch serviceType {
	case config.GraphService:
		seid := nodeConfig.SSH.Address + ":" + strconv.Itoa(serviceConfig.Port)
		graphExporter := new(GraphExporter)
		graphExporter.Config(nodeConfig, serviceConfig)
		graphExporter.BuildAllMap()
		exporters[seid] = graphExporter
		return graphExporter, nil
	case config.MetaService:
		seid := nodeConfig.SSH.Address + ":" + strconv.Itoa(serviceConfig.Port)
		metaExporter := new(MetaExporter)
		metaExporter.Config(nodeConfig, serviceConfig)
		metaExporter.BuildAllMap()
		exporters[seid] = metaExporter
		return metaExporter, nil
	case config.StorageService:
		seid := nodeConfig.SSH.Address + ":" + strconv.Itoa(serviceConfig.Port)
		storageExporter := new(StorageExporter)
		storageExporter.Config(nodeConfig, serviceConfig)
		storageExporter.BuildAllMap()
		exporters[seid] = storageExporter
		return storageExporter, nil
	default:
		return nil, errors.New("init service exporter fail")
	}

}

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
