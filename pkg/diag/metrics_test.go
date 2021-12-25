package diag

import (
	"fmt"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/info/metrics"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/info/service"
	"testing"
)

func TestGetMetricsDiag(t *testing.T) {
	serviceConf := config.ServiceConfig{
		Type: config.GraphService,
		HTTPPort: 19669,
		RuntimeDir: "/home/katz.zhang/.nebula/clusters/graphd",
	}
	nodeConf := &config.NodeConfig{
		SSH: config.SSHConfig{
			Address:  "192.168.8.49",
			Port:     22,
			Timeout:  "1s",
			Username: "katz.zhang",
			Password: "nebula",
		},
		Output: config.OutputConfig{
			DirPath: "",
		},
		Services: map[string]config.ServiceConfig{
			"graph1":serviceConf,
		},
	}
	metrics, _ := metrics.GetMetricsInfo(nodeConf, &serviceConf)
	serviceMetricsInfo := &service.ServiceMetricInfo{
		Name:    "graph",
		Metrics: metrics,
		Type: serviceConf.Type,
	}
	var metricsSlice []*service.ServiceMetricInfo
	metricsSlice = append(metricsSlice, serviceMetricsInfo)
	diags := GetMetricsDiag(metricsSlice)
	fmt.Println(diags)

}
