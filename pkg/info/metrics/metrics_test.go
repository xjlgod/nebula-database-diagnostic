package metrics

import (
	"fmt"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	"testing"
)

func TestGetMetrics(t *testing.T) {
	nodeConf := &config.NodeConfig{
		SSH: config.SSHConfig{
			Address:  "192.168.8.169",
			Port:     22,
			Timeout:  "1s",
			Username: "katz.zhang",
			Password: "nebula",
		},
	}
	serviceConf := &config.ServiceConfig{
		Type: config.GraphService,
		HTTPPort: 19669,
	}
	info, _ := GetMetricsInfo(nodeConf, serviceConf)
	for key, value := range info {
		fmt.Printf("key: %s, value :%s \n", key, value)
	}
}